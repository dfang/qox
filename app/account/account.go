package account

import (
	"strconv"

	"github.com/dfang/qor-demo/config/auth"
	"github.com/dfang/qor-demo/models/users"
	"github.com/dfang/qor-demo/utils/funcmapmaker"
	"github.com/go-chi/chi"
	"github.com/qor/application"
	"github.com/qor/qor"
	"github.com/qor/render"
)

// New new home app
func New(config *Config) *App {
	return &App{Config: config}
}

// NewWithDefault New App With No Config
func NewWithDefault() *App {
	return &App{Config: &Config{}}
}

// App home app
type App struct {
	Config *Config
}

// Config home config struct
type Config struct {
}

// ConfigureApplication configure application
func (app App) ConfigureApplication(application *application.Application) {
	controller := &Controller{View: render.New(&render.Config{AssetFileSystem: application.AssetFS.NameSpace("account")}, "app/account/views")}

	funcmapmaker.AddFuncMapMaker(controller.View)
	app.ConfigureAdmin(application.Admin)

	application.Router.Mount("/auth/", auth.Auth.NewServeMux())

	application.Router.With(auth.Authority.Authorize()).Route("/account", func(r chi.Router) {
		r.Get("/", controller.Orders)
		r.With(auth.Authority.Authorize("logged_in_half_hour")).Post("/add_user_credit", controller.AddCredit)
		r.Get("/profile", controller.Profile)
		r.Post("/profile", controller.Update)
	})
}

func userAddressesCollection(resource interface{}, context *qor.Context) (results [][]string) {
	var (
		user users.User
		DB   = context.DB
	)

	DB.Preload("Addresses").Where(context.ResourceID).First(&user)

	for _, address := range user.Addresses {
		results = append(results, []string{strconv.Itoa(int(address.ID)), address.Stringify()})
	}
	return
}
