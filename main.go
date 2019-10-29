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
	"github.com/dfang/qor-demo/app/products"
	"github.com/dfang/qor-demo/app/static"
	"github.com/dfang/qor-demo/config"
	"github.com/dfang/qor-demo/config/auth"
	"github.com/dfang/qor-demo/config/db"
	"github.com/dfang/qor-demo/models/aftersales"
	"github.com/dfang/qor-demo/models/users"
	"github.com/dfang/qor-demo/utils/funcmapmaker"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/heptiolabs/healthcheck"
	"github.com/qor/admin"
	"github.com/qor/application"
	"github.com/qor/publish2"
	"github.com/qor/qor"
	"github.com/qor/qor/utils"

	"github.com/gocraft/work"

	// https://github.com/qor/qor-example/issues/129
	"github.com/dfang/qor-demo/config/db/migrations"
)

const version = "0.0.1" // must follow semver spec, https://github.com/motemen/gobump

var (
	Router      *chi.Mux
	Admin       *admin.Admin
	Application *application.Application
)

func main() {
	start := time.Now()
	fmt.Println("Now is ", time.Now().Format("2006-01-02 15:04:05"))

	cmdLine := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	compileTemplate := cmdLine.Bool("compile-templates", false, "Compile Templates")
	isDebug, _ := strconv.ParseBool(os.Getenv("DEBUG"))
	debug := cmdLine.Bool("debug", isDebug, "Set log level to debug")
	runMigration := cmdLine.Bool("migrate", false, "Run migration")
	// runSeed := cmdLine.Bool("seed", false, "Run seed")

	cmdLine.Parse(os.Args[1:])

	migrations.Migrate()

	if *runMigration {
		migrations.Migrate()
		os.Exit(0)
	}

	health := healthcheck.NewHandler()
	// Our app is not happy if we've got more than 100 goroutines running.
	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(100))
	// Our app is not ready if we can't resolve our upstream dependency in DNS.
	// health.AddReadinessCheck(
	// 	"upstream-dep-dns",
	// 	healthcheck.DNSResolveCheck("upstream.example.com", 50*time.Millisecond))

	// Our app is not ready if we can't connect to our database (`var db *sql.DB`) in <1s.
	health.AddReadinessCheck("database", healthcheck.DatabasePingCheck(db.DB.DB(), 5*time.Second))
	go http.ListenAndServe("0.0.0.0:8086", health)

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

	Router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			DumpHTTPRequest(req)
			next.ServeHTTP(w, req)
		})
	})

	Application.Use(aftersale.NewWithDefault())
	Application.Use(api.New(&api.Config{}))
	Application.Use(adminapp.New(&adminapp.Config{}))
	// Application.Use(home.New(&home.Config{}))
	Application.Use(account.NewWithDefault())
	Application.Use(home.NewWithDefault())
	Application.Use(products.NewWithDefault())
	Application.Use(orders.NewWithDefault())
	// Application.Use(pages.NewWithDefault())
	// Application.Use(stores.NewWithDefault())

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

	fmt.Println("start cron job ......")
	go startWorkerPool()

	if *compileTemplate {
		bindatafs.AssetFS.Compile()
		os.Exit(0)
	}

	envs := []string{
		"GO_ENV", "DEBUG", "HTTPS", "DOMAIN", "PORT",
	}

	for _, e := range envs {
		if os.Getenv(e) != "" {
			fmt.Printf("Found env var %s=%s\n", e, os.Getenv(e))
		}
	}

	if os.Getenv("HTTPS") == "true" && os.Getenv("DOMAIN") == "" {
		log.Info().Msg("If set HTTPS=true, this app will get ssl certificates automatically and serve on 443,  so you must also set DOMAIN")
		log.Info().Msg("By default (HTTPS not set), you need to config caddy as reverse proxy to serve https requests")
		log.Info().Msg("If you plan to use caddy as frontend (reverse proxy), you don't need to set HTTPS, or just set to false")
		os.Exit(1)
	}

	elapsed := time.Since(start)
	fmt.Printf("Startup took %s\n", elapsed)
	fmt.Printf("Listening on: %v\n", config.Config.Port)

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

// DumpHTTPRequest for debugging
func DumpHTTPRequest(r *http.Request) {
	output, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println("Error when dumping request:", err)
		return
	}
	fmt.Println(string(output))
}

// Context For gocraft/work
type Context struct {
	userID int64
}

// ExpireAftersales 任务指派后 after_sale的状态为scheduled， 如果师傅20分钟之内没有响应，自动变为overdue状态
func ExpireAftersales(job *work.Job) error {
	// time.Sleep(10 * time.Second)
	fmt.Println("now is", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("expires all scheduled aftersales that idle for 20 minutes ......")
	var items []aftersales.Aftersale
	db.DB.Model(aftersales.Aftersale{}).Where("state = ?", "scheduled").Where("updated_at <= NOW() - INTERVAL '20 minutes'").Find(&items)
	// .Update("state", "overdue")
	fmt.Println(len(items))
	for _, item := range items {
		fmt.Println("before: ", item.State)
		aftersales.OrderStateMachine.Trigger("expire", &item, db.DB, "expires aftersale with id: "+fmt.Sprintf("%d", item.ID))
		fmt.Println("after:", item.State)
		// db.DB.Model(&item).Update("state", "overdue")
		db.DB.Save(&item)
	}
	fmt.Println("expires aftersales done ")

	return nil
}

// FreezeAftersales 已审核的服务单冻结7天才能结算
func FreezeAftersales(job *work.Job) error {
	fmt.Println("now is", time.Now().Format("2006-01-02 15:04:05"))
	// time.Sleep(70 * time.Second)
	fmt.Println("freeze aftersales ......")
	var items []aftersales.Aftersale
	// db.DB.Model(aftersales.Aftersale{}).Where("state = ?", "audited").Update("state", "frozen")
	db.DB.Model(aftersales.Aftersale{}).Where("state = ?", "audited").Find(&items)
	for _, item := range items {
		aftersales.OrderStateMachine.Trigger("freeze", &item, db.DB, "freeze aftersale with id: "+fmt.Sprintf("%d", item.ID))
		db.DB.Save(&item)
	}
	fmt.Println("freeze aftersales done ......")

	return nil
}

// UnfreezeAftersales 解冻超过7天的，自动结算，金额算到师傅名下
func UnfreezeAftersales(job *work.Job) error {
	fmt.Println("now is", time.Now().Format("2006-01-02 15:04:05"))
	// time.Sleep(55 * time.Second)
	fmt.Println("unfreeze aftersales ......")
	var items []aftersales.Aftersale
	db.DB.Model(aftersales.Aftersale{}).Where("state = ?", "frozen").Find(&items)
	for _, item := range items {
		aftersales.OrderStateMachine.Trigger("unfreeze", &item, db.DB, "unfreeze aftersale with id: "+fmt.Sprintf("%d", item.ID))
		db.DB.Save(&item)
	}
	fmt.Println("unfreeze aftersales done ......")

	return nil
}

// UpdateBalances 统计每个师傅的冻结金额和可结算金额并更新到Balances表
func UpdateBalances(job *work.Job) error {
	var workmen []users.User
	db.DB.Select("name, id").Where("role = ?", "workman").Find(&workmen)

	for _, item := range workmen {
		// 计算frozen_amount
		// 计算free_amount
		// update balance by user_id
		var balance aftersales.Balance
		db.DB.Model(aftersales.Balance{}).Where("user_id = ?", item.ID).Assign(aftersales.Balance{UserID: item.ID}).FirstOrInit(&balance)

		// select sum(amount) from settlements where user_id = 73 and state='frozen';

		// var frozenResult float32
		// var freeResult float32
		// db.DB.Table("settlements").Select("sum(amount)").Where("state = 'frozen'").Where("user_id = ?", item.ID).Take(&frozenResult)
		// db.DB.Table("settlements").Select("sum(amount)").Where("state = 'free'").Where("user_id = ?", item.ID).Take(&freeResult)
		type Result struct {
			State string
			Total float32
		}
		// rows, err :=
		var results []Result
		var f1, f2, f3 float32

		db.DB.Table("settlements").Select("state, sum(amount) as total").Group("state").Where("user_id = ?", item.ID).Scan(&results)
		for _, i := range results {
			fmt.Println(i.State)
			fmt.Println(i.Total)
			if i.State == "frozen" {
				f1 = i.Total
			}

			if i.State == "free" {
				f2 = i.Total
			}

			if i.State == "withdrawed" {
				f3 = i.Total
			}
		}

		balance.FrozenAmount = f1
		balance.FreeAmount = f2 + f3
		balance.WithdrawAmount = f3
		balance.TotalAmount = f2 + f1

		// balance.FrozenAmount = balance.FrozenAmount + balance.FreeAmount + balance.WithdrawAmount

		// balance.UserID = item.ID
		// balance.FrozenAmount = frozenResult
		// balance.FreeAmount = freeResult

		db.DB.Save(&balance)
	}

	return nil
}

func startWorkerPool() {
	// Periodic Enqueueing (Cron)
	pool := work.NewWorkerPool(Context{}, 10, "qor", db.RedisPool)
	pool.PeriodicallyEnqueue("30 * * * * *", "expire_aftersales") // This will enqueue a "expire_aftersales" job every minutes
	pool.PeriodicallyEnqueue("30 * * * * *", "freeze_audited_aftersales")
	pool.PeriodicallyEnqueue("5 * * * *", "unfreeze_aftersales")
	pool.PeriodicallyEnqueue("30 * * * * *", "update_balances")

	pool.Job("expire_aftersales", ExpireAftersales) // Still need to register a handler for this job separately
	pool.Job("freeze_audited_aftersales", FreezeAftersales)
	pool.Job("unfreeze_aftersales", UnfreezeAftersales)
	pool.Job("update_balances", UpdateBalances)

	// Start processing jobs
	pool.Start()
	// // Wait for a signal to quit:
	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt, os.Kill)
	// <-signalChan

	// // Stop the pool
	// pool.Stop()
}
