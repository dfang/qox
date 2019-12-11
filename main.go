package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/acme/autocert"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/dfang/qor-demo/config/bindatafs"

	"github.com/dfang/qor-demo/app/account"
	adminapp "github.com/dfang/qor-demo/app/admin"
	"github.com/dfang/qor-demo/app/aftersale"
	"github.com/dfang/qor-demo/app/api"
	"github.com/dfang/qor-demo/app/home"
	"github.com/dfang/qor-demo/app/orders"
	"github.com/dfang/qor-demo/app/pages"
	"github.com/dfang/qor-demo/app/products"

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

const version = "0.0.1" // must follow semver spec, https://github.com/motemen/gobump

var (
	Router      *chi.Mux
	Admin       *admin.Admin
	Application *application.Application
)

func main() {
	cmdLine := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	compileTemplate := cmdLine.Bool("compile-templates", false, "Compile Templates")
	isDebug, _ := strconv.ParseBool(os.Getenv("DEBUG"))
	debug := cmdLine.Bool("debug", isDebug, "Set log level to debug")
	runMigration := cmdLine.Bool("migrate", false, "Run migration")
	runSeeds := cmdLine.Bool("seeds", false, "Run seeds, never run this on production")
	ui := cmdLine.Bool("ui", false, "Serves gocraft/work ui")
	// runSeed := cmdLine.Bool("seed", false, "Run seed")
	eval := cmdLine.Bool("eval", false, "Evaluate rules")

	cmdLine.Parse(os.Args[1:])

	Compile()

	if *eval {
		Evaluate()
		os.Exit(0)
	}

	fmt.Println("check availables enviroments variables ......")
	checkAvailableEnvs()

	fmt.Println("initialze configurations ......")
	initialzeConfigs()

	fmt.Println("set log level ......")
	setupLogLevel(*debug)

	if *runMigration {
		fmt.Println("just run migrations and exit ......")
		migrations.Migrate()
		os.Exit(0)
	}

	// fmt.Println("start auto migrations ......")
	// migrations.Migrate()

	if *runSeeds && os.Getenv("QOR_ENV") != "production" {
		fmt.Println("just run seeds and exit ......")
		fmt.Println("start truncate tables ......")
		seeds.TruncateTables()
		fmt.Println("start seeding samples data for testing ......")
		seeds.Run()
		fmt.Println("seeding done, exit ......")
		os.Exit(0)
	}

	fmt.Println("setup middlewares and routes ......")
	setupMiddlewaresAndRoutes()

	if *compileTemplate {
		fmt.Println("just compile templates and exit ......")
		bindatafs.AssetFS.Compile()
		os.Exit(0)
	}

	// fmt.Println("start workerPool ......")
	// if os.Getenv("ENV") == "development" {
	if false {
		go startWorkerPool()
	}

	fmt.Println("start health check ......")
	go startHealthCheck()

	if *ui || os.Getenv("UI") == "true" {
		fmt.Println("serves gocraft/work web ui ......")
		go startWorkWebUI()
	}

	if os.Getenv("HTTPS") == "true" && os.Getenv("DOMAIN") == "" {
		log.Info().Msg("If set HTTPS=true, this app will get ssl certificates automatically and serve on 443,  so you must also set DOMAIN")
		log.Info().Msg("By default (HTTPS not set), you need to config caddy as reverse proxy to serve https requests")
		log.Info().Msg("If you plan to use caddy as frontend (reverse proxy), you don't need to set HTTPS, or just set to false")
		os.Exit(1)
	}

	fmt.Println("NOW is ", time.Now().Format("2006-01-02 15:04:05"))
	elapsed := time.Since(*config.StartUpStartTime)
	fmt.Printf("Startup took %s\n", elapsed)

	fmt.Printf("Listening on: %v\n\n", config.Config.Port)
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
			fmt.Printf("\tfound env var %s=%s\n", e, os.Getenv(e))
		}
	}
}

func setupLogLevel(debug bool) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Logger = log.With().Caller().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func initialzeConfigs() {
	// config.Initialize()
	// db.Initialize()
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
	Application.Use(pages.NewWithDefault())
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

	Router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("PONG"))
	})

	// fs := http.FileServer(http.Dir("public/"))
	// http.Handle("/", http.StripPrefix("/", fs))
	// Router.Handle("/static/", http.StripPrefix("/static/", fs))

	// Router.Handle("/static/", http.RedirectHandler("http://baidu.com", 303))
	// Router.Handle("/public/", http.StripPrefix("/public/", fs))

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "public")
	FileServer(Router, "/", http.Dir(filesDir))
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

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
