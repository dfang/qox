package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/acme/autocert"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/dfang/qor-demo/app/account"
	adminapp "github.com/dfang/qor-demo/app/admin"
	"github.com/dfang/qor-demo/app/aftersale"
	"github.com/dfang/qor-demo/app/api"
	"github.com/dfang/qor-demo/app/home"
	"github.com/dfang/qor-demo/app/orders"
	"github.com/dfang/qor-demo/app/products"
	"github.com/dfang/qor-demo/cmd"

	"github.com/dfang/qor-demo/app/reports"
	"github.com/dfang/qor-demo/app/stores"
	"github.com/dfang/qor-demo/config"
	"github.com/dfang/qor-demo/config/auth"
	"github.com/dfang/qor-demo/config/db"
	"github.com/dfang/qor-demo/config/i18n"
	"github.com/dfang/qor-demo/utils/funcmapmaker"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/qor/admin"
	"github.com/qor/application"
	"github.com/qor/publish2"
	"github.com/qor/qor"
	"github.com/qor/qor/utils"
	"github.com/rs/cors"

	// https://github.com/qor/qor-example/issues/129

	"github.com/dfang/qor-demo/config/db/migrations"
	"github.com/dfang/qor-demo/config/db/seeds"
)

var version = "0.0.1" // must follow semver spec, https://github.com/motemen/gobump
var buildVersion = "development"

var (
	// Router Chi Router
	Router *chi.Mux
	// Admin Qor Admin (后台管理页面)
	Admin *admin.Admin
	// Application Qor Application (前端)
	Application *application.Application
)

func checkAvailableEnvs() {
	envs := []string{
		"QOR_ENV",
		"DEBUG",
		"HTTPS",
		"DOMAIN",
		"PORT",
		"DBAdapter",
		"DBName",
		"DBPort",
		"DBHost",
		"DBUser",
		"DBPassword",
		"REDIS_HOST",
		"REDIS_PORT",
		"DEMO_MODE",
	}

	for _, e := range envs {
		if os.Getenv(e) != "" {
			log.Debug().Msgf("\tfound env var %s=%s", e, os.Getenv(e))
		}
	}
}

func setLogLevel(level int) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Logger = log.With().Caller().Logger()

	var l zerolog.Level
	switch level {
	case -1:
		l = zerolog.InfoLevel
	case 0:
		l = zerolog.DebugLevel
	case 1:
		l = zerolog.InfoLevel
	case 2:
		l = zerolog.WarnLevel
	case 3:
		l = zerolog.ErrorLevel
	case 4:
		l = zerolog.FatalLevel
	case 5:
		l = zerolog.PanicLevel
	default:
		l = zerolog.NoLevel
	}

	log.Info().Msgf("set log level to %s ......", l.String())
	zerolog.SetGlobalLevel(l)
}

func initialzeConfigs() {
	config.Initialize()
	db.Initialize()
	i18n.Initialize()
	auth.Initialize()
}

func setupMiddlewaresAndRoutes() {
	Router = chi.NewRouter()
	Admin = admin.New(&admin.AdminConfig{
		SiteName: "QOR DEMO",
		Auth:     auth.AdminAuth{},
		DB:       db.DB.Set(publish2.VisibleMode, publish2.ModeOff).Set(publish2.ScheduleMode, publish2.ModeOff),
	})
	Application = application.New(&application.Config{
		Router: Router,
		Admin:  Admin,
		DB:     Admin.DB,
	})

	funcmapmaker.AddFuncMapMaker(auth.Auth.Config.Render)

	// Health check for k8s pod
	Router.Use(middleware.Heartbeat("/health"))
	// TODO: implement Graceful shutdown
	// https://medium.com/over-engineering/graceful-shutdown-with-go-http-servers-and-kubernetes-rolling-updates-6697e7db17cf

	// https://github.com/qor/qor-example/commit/06835622f5feeeb90aee4f02574abbae29e40e10
	Router.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			req.Header.Del("Authorization")
			handler.ServeHTTP(w, req)
		})
	})

	Router.Use(middleware.RealIP)
	Router.Use(middleware.RequestID)
	Router.Use(middleware.Logger)
	Router.Use(middleware.Recoverer)
	Router.Use(middleware.Timeout(3 * time.Second))
	// Router.Use(middleware.WithValue("k", "v"))
	// Router.Use(middleware.NoCache)
	// Router.Use(middleware.StripSlashes)

	// Use default options
	Router.Use(cors.Default().Handler)

	Router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			var (
				tx         = db.DB
				qorContext = &qor.Context{Request: req, Writer: w}
			)

			if locale := utils.GetLocale(qorContext); locale != "" {
				tx = tx.Set("l10n:locale", locale)
			} else {
				// set default locale to zh-CN
				// tx = tx.Set("l10n:locale", "zh-CN")
				expire := time.Now().AddDate(0, 1, 0) // one month later
				cookie := http.Cookie{Name: "locale", Value: "zh-CN", Path: "/", Expires: expire, MaxAge: 90000}
				http.SetCookie(w, &cookie)
			}

			ctx := context.WithValue(req.Context(), utils.ContextDBName, publish2.PreviewByDB(tx, qorContext))
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	})

	if os.Getenv("DEBUG") == "true" {
		Router.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				DumpHTTPRequest(req)
				next.ServeHTTP(w, req)
			})
		})
	}

	Application.Use(aftersale.NewWithDefault())
	Application.Use(api.New(&api.Config{}))
	Application.Use(adminapp.New(&adminapp.Config{}))
	// Application.Use(home.New(&home.Config{}))
	Application.Use(account.NewWithDefault())
	Application.Use(home.NewWithDefault())
	Application.Use(products.NewWithDefault())
	Application.Use(orders.NewWithDefault())
	// Application.Use(pages.NewWithDefault())
	Application.Use(stores.NewWithDefault())
	Application.Use(reports.NewWithDefault())
	// Application.Use(enterprise.New(&enterprise.Config{}))

	// views := []string{
	// 	"app/home/views",
	// 	"app/account/views",
	// 	"app/order/views",
	// 	"app/page/views",
	// 	"app/products/views",
	// 	"app/views",
	// }

	// for _, v := range views {
	// 	// bindatafs.AssetFS.NameSpace("views").RegisterPath(filepath.Join(config.Root, v))
	// 	bindatafs.AssetFS.RegisterPath(filepath.Join(config.Root, v))
	// }

	Router.Mount("/debug", middleware.Profiler())
	Router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("PONG"))
	})

	// fs := http.FileServer(http.Dir("public"))
	// Router.Handle("/", fs)

	// http.Handle("/", http.StripPrefix("/", fs))
	// Router.Handle("/static/", http.StripPrefix("/static/", fs))

	// Router.Handle("/static/", http.RedirectHandler("http://baidu.com", 303))
	// Router.Handle("/public/", http.StripPrefix("/public/", fs))

	// https://www.alexedwards.net/blog/serving-static-sites-with-go
	// https://github.com/go-chi/chi/blob/master/middleware/nocache.go
	// https://github.com/go-chi/chi/blob/master/middleware/profiler.go

	// workDir, _ := os.Getwd()
	// filesDir := filepath.Join(workDir, "public")
	// FileServer(Router, "/", http.Dir(filesDir))

	// 用NotFound 来替代前面的三句
	// Serve public 文件夹里的 favicon.ico 等文件
	// 确保 /favicon.ico 等能直接访问
	// 注意 NotFound 和 pages 模块冲突，所以注释了pages
	// 因为貌似开启了pages模块 永远不会执行到NotFound里来
	Router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		// w.WriteHeader(404)
		// w.Write([]byte(r.URL.Path))
		// w.Write([]byte("nothing here"))
		log.Debug().Msg("request path: " + r.URL.Path)
		// workDir, _ := os.Getwd()
		// filesDir := filepath.Join(workDir, "public")
		fs := http.FileServer(http.Dir("./public"))
		fs.ServeHTTP(w, r)
		// http.ServeFile(w, r, "./public/favicon.ico")
	})
}

// DumpHTTPRequest for debugging
func DumpHTTPRequest(r *http.Request) {
	output, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println("Error when dumping request:", err)
		return
	}
	fmt.Println(string(output))
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
// https://github.com/go-chi/chi/blob/master/_examples/fileserver/main.go
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	// https://stackoverflow.com/questions/27945310/why-do-i-need-to-use-http-stripprefix-to-access-my-static-files
	fs := http.StripPrefix(path, http.FileServer(root))

	// fmt.Println("fuck")
	// fmt.Println(path)
	// fmt.Println(len(path) - 1)
	// fmt.Println(path[len(path)-1])

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Debug().Msg("FALLBACK")

		fmt.Println("fallback to file server")
		fmt.Println(path)

		fs.ServeHTTP(w, r)
	}))
}

func listenAndServe() {
	if os.Getenv("HTTPS") == "true" && os.Getenv("DOMAIN") == "" {
		log.Info().Msg("If set HTTPS=true, this app will get ssl certificates automatically and serve on 443,  so you must also set DOMAIN")
		log.Info().Msg("By default (HTTPS not set), you need to config caddy as reverse proxy to serve https requests")
		log.Info().Msg("If you plan to use caddy as frontend (reverse proxy), you don't need to set HTTPS, or just set to false")
		os.Exit(1)
	}

	// fmt.Println("NOW is ", time.Now().Format("2006-01-02 15:04:05"))
	elapsed := time.Since(*config.StartUpStartTime)
	log.Debug().Msgf("Startup took %s\n", elapsed)
	log.Info().Msgf("Listening on: %v\n", config.Config.Port)

	if config.Config.HTTPS {
		certManager := autocert.Manager{
			Prompt: autocert.AcceptTOS,
			Cache:  autocert.DirCache("cert-cache"),
			// Put your domain here:
			HostPolicy: autocert.HostWhitelist(os.Getenv("DOMAIN")),
		}
		server := &http.Server{
			Addr:    ":443",
			Handler: Router,
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
			},
		}

		go http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), certManager.HTTPHandler(nil))

		server.ListenAndServeTLS("", "")

	} else {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), Application.NewServeMux()); err != nil {
			panic(err)
		}
	}
}

func runMainAction(c *cli.Context) error {
	log.Info().Msg("check availables enviroments variables ......")
	checkAvailableEnvs()

	log.Info().Msg("initialze configurations ......")
	initialzeConfigs()

	log.Info().Msg("start faktory worker")
	go cmd.StartFaktoryWorker()

	log.Info().Msg("start gocraft/work workerPool ......")
	go cmd.StartWorkerPool()
	if c.Bool("ui") || os.Getenv("UI") == "true" {
		log.Info().Msg("start gocraft/work web ui ......")
		go cmd.StartWorkWebUI()
	}

	log.Info().Msg("start health check ......")
	go cmd.StartHealthCheck()

	log.Info().Msg("start webhookd")
	go cmd.StartWebhookd()

	log.Info().Msg("setup middlewares and routes ......")
	setupMiddlewaresAndRoutes()

	listenAndServe()

	return nil
}

func printVersion() {
	fmt.Printf("Version %s, BuildVersion %s\n", version, buildVersion)
}

func runMigrations() {
	log.Info().Msg("just run migrations and exit ......")
	migrations.Migrate()
}

func runSeeds() {
	log.Info().Msg("start truncate tables ......")
	seeds.TruncateTables()
	log.Info().Msg("start seeding samples data for testing ......")
	seeds.Run()
	log.Info().Msg("seeding done, exit ......")
}

func fail(err error) {
	log.Fatal().Msg(err.Error())
	os.Exit(-1)
}

func main() {
	app := &cli.App{
		Name:  "qor",
		Usage: "make an explosive entrance",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "migrate",
				Aliases: []string{"m"},
				Value:   false,
				Usage:   "Run migrations",
			},
			&cli.BoolFlag{
				Name:    "seeds",
				Aliases: []string{"s"},
				Value:   false,
				Usage:   "Run seeds",
			},
			&cli.BoolFlag{
				Name:    "compile",
				Aliases: []string{"c"},
				Value:   false,
				Usage:   "Compile grool templates",
			},
			&cli.BoolFlag{
				Name:  "eval",
				Value: false,
				Usage: "Eval grool templates",
			},
			&cli.BoolFlag{
				Name:  "ui",
				Value: false,
				Usage: "Start go worker ui",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Value: false,
				Usage: "DEBUG MODE (v=0)",
			},
			&cli.IntFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Value:   1,
				Usage:   "Set log level",
			},
			&cli.BoolFlag{
				Name:    "version",
				Aliases: []string{"V"},
				Usage:   "Show version",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "version",
				Usage: "Show version",
				Action: func(c *cli.Context) error {
					printVersion()
					return nil
				},
			},
			{
				Name:    "migrate",
				Aliases: []string{"m"},
				Usage:   "Run migrations",
				Action: func(c *cli.Context) error {
					initialzeConfigs()
					runMigrations()
					return nil
				},
			},
			{
				Name:    "seed",
				Aliases: []string{"s"},
				Usage:   "Run seeding",
				Action: func(c *cli.Context) error {
					initialzeConfigs()
					runSeeds()
					return nil
				},
			},
		},
		Action: func(c *cli.Context) error {
			if c.Bool("V") {
				fmt.Printf("Version %s, BuildVersion %s\n", version, buildVersion)
				os.Exit(0)
			}

			level := c.Int("v")
			setLogLevel(level)

			if c.Bool("debug") {
				setLogLevel(0)
			}

			runMainAction(c)

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}
