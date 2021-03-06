package home

import (
	"github.com/qor/application"
	// "github.com/dfang/qor-demo/config/bindatafas"
	"github.com/dfang/qor-demo/utils/funcmapmaker"
	"github.com/qor/render"
)

// New new home app
func New(config *Config) *App {
	return &App{Config: config}
}

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
func (App) ConfigureApplication(application *application.Application) {
	controller := &Controller{View: render.New(&render.Config{AssetFileSystem: application.AssetFS.NameSpace("home")}, "app/home/views")}
	// controller := &Controller{ View: render.New( &render.Config{ AssetFileSystem: bindatafs.AssetFS.NameSpace("home"), }),

	funcmapmaker.AddFuncMapMaker(controller.View)

	application.Router.Get("/", controller.RedirectToAdmin)
	// application.Router.Get("/", controller.Index)

	application.Router.Get("/switch_locale", controller.SwitchLocale)
}
