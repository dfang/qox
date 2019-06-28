package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
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
	"github.com/dfang/qor-demo/app/api"
	"github.com/dfang/qor-demo/app/home"
	"github.com/dfang/qor-demo/app/orders"
	"github.com/dfang/qor-demo/app/pages"
	"github.com/dfang/qor-demo/app/products"
	"github.com/dfang/qor-demo/app/static"
	"github.com/dfang/qor-demo/app/stores"
	"github.com/dfang/qor-demo/config"
	"github.com/dfang/qor-demo/config/auth"
	"github.com/dfang/qor-demo/config/db"
	"github.com/dfang/qor-demo/utils/funcmapmaker"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/application"
	"github.com/qor/publish2"
	"github.com/qor/qor"
	"github.com/qor/qor/utils"

	// https://github.com/qor/qor-example/issues/129
	"github.com/dfang/qor-demo/config/db/migrations"
)

var (
	Router      *chi.Mux
	Admin       *admin.Admin
	Application *application.Application
)

func main() {
	start := time.Now()

	cmdLine := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	compileTemplate := cmdLine.Bool("compile-templates", false, "Compile Templates")
	isDebug, _ := strconv.ParseBool(os.Getenv("DEBUG"))
	debug := cmdLine.Bool("debug", isDebug, "Set log level to debug")
	runMigration := cmdLine.Bool("migration", false, "Run migration")
	// runSeed := cmdLine.Bool("seed", false, "Run seed")

	cmdLine.Parse(os.Args[1:])

	if *runMigration {
		migrations.Migrate()
		os.Exit(0)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Logger = log.With().Caller().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

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

	Router.Use(middleware.RealIP)
	Router.Use(middleware.Logger)
	Router.Use(middleware.Recoverer)

	// subdomain support
	Admin.GetRouter().Use(&admin.Middleware{
		Name: "switch_db",
		Handler: func(context *admin.Context, middleware *admin.Middleware) {
			fmt.Println("Host: ", context.Request.Host)
			host := context.Request.Host
			//Figure out if a subdomain exists in the host given.
			hostParts := strings.Split(host, ".")
			var subDomain string
			if len(hostParts) > 2 {
				subDomain = hostParts[0]
			}
			fmt.Println("SubDomain: ", subDomain)
			if subDomain != "" && subDomain != "www" {
				// http://codepodu.com/subdomains-with-golang/
				// https://stackoverflow.com/questions/26517636/getting-the-subdomain
				// https://www.learngoogle.com/2013/07/14/extract-subdomain-from-request-in-go/
				subDomain := strings.Split(context.Request.Host, ".")[0]
				fmt.Println(fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", "postgres", "postgres", "localhost", "5432", subDomain))
				DB, err := gorm.Open("postgres", fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", "postgres", "postgres", "localhost", "5432", subDomain))
				if err != nil {
					log.Fatal().Str("err", err.Error())
				}
				err = DB.DB().Ping()
				if err != nil {
					log.Fatal().Str("err", err.Error())
				}
				context.SetDB(DB)
			}
			middleware.Next(context)
		},
	})

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

	Application.Use(api.New(&api.Config{}))
	Application.Use(adminapp.New(&adminapp.Config{}))
	// Application.Use(home.New(&home.Config{}))
	Application.Use(account.NewWithDefault())
	Application.Use(home.NewWithDefault())
	Application.Use(products.NewWithDefault())
	Application.Use(orders.NewWithDefault())
	Application.Use(pages.NewWithDefault())
	Application.Use(stores.NewWithDefault())
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

	Application.Use(static.New(&static.Config{
		Prefixs: []string{"/system"},
		Handler: utils.FileServer(http.Dir(filepath.Join(config.Root, "public"))),
	}))
	Application.Use(static.New(&static.Config{
		Prefixs: []string{"javascripts", "stylesheets", "images", "dist", "fonts", "vendors", "downloads", "favicon.ico"},
		Handler: bindatafs.AssetFS.FileServer(http.Dir("public"), "javascripts", "stylesheets", "images", "dist", "downloads", "fonts", "vendors", "favicon.ico"),
	}))

	if *compileTemplate {
		bindatafs.AssetFS.Compile()
		os.Exit(0)
	} else {

		elapsed := time.Since(start)
		fmt.Printf("Startup took %s\n", elapsed)
		fmt.Printf("Listening on: %v\n", config.Config.Port)
		if os.Getenv("GO_ENV") != "production" {
			if err := http.ListenAndServe(fmt.Sprintf("app.localhost:%d", config.Config.Port), Application.NewServeMux()); err != nil {
				panic(err)
			}
		} else {
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
	}
}
