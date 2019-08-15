package aftersale

// "net/http"
import (
	"github.com/dfang/qor-demo/models/aftersales"
	"github.com/qor/admin"
	"github.com/qor/application"
)

// New new home app
func New(config *Config) *App {
	return &App{Config: config}
}

// NewWithDefault new home app
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
	// 售后后台
	app.ConfigureAdmin(application.Admin)
}

// ConfigureAdmin configure admin interface
func (App) ConfigureAdmin(Admin *admin.Admin) {
	Admin.AddMenu(&admin.Menu{Name: "Aftersale Management", Priority: 6})

	// Add Aftersale
	aftersale := Admin.AddResource(&aftersales.AfterSale{}, &admin.Config{Menu: []string{"Aftersale Management"}})

	aftersale.Meta(&admin.Meta{
		Name:       "ServiceType",
		Type:       "select_one",
		Collection: []string{"安装", "维修", "清洗"},
	})
}
