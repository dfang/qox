package api

import (
	"github.com/dfang/qor-demo/config/db"
	"github.com/dfang/qor-demo/models/orders"
	// "github.com/dfang/qor-demo/models/products"
	"github.com/dfang/qor-demo/models/users"
	"github.com/dfang/qor-demo/models/aftersales"
	"github.com/qor/admin"
	"github.com/qor/application"
	"github.com/qor/qor"
)

// New new home app
func New(config *Config) *App {
	if config.Prefix == "" {
		config.Prefix = "/api"
	}
	return &App{Config: config}
}

func NewWithDefault() *App {
	return &App{Config: &Config{
		Prefix: "/api",
	}}
}

// App home app
type App struct {
	Config *Config
}

// Config home config struct
type Config struct {
	Prefix string
}

// ConfigureApplication configure application
func (app App) ConfigureApplication(application *application.Application) {
	API := admin.New(&qor.Config{DB: db.DB})

	// Product := API.AddResource(&products.Product{})

	// ColorVariationMeta := Product.Meta(&admin.Meta{Name: "ColorVariations"})
	// ColorVariation := ColorVariationMeta.Resource
	// ColorVariation.IndexAttrs("ID", "Color", "Images", "SizeVariations")
	// ColorVariation.ShowAttrs("Color", "Images", "SizeVariations")

	// // SizeVariationMeta := ColorVariation.Meta(&admin.Meta{Name: "SizeVariations"})
	// SizeVariationMeta := Product.Meta(&admin.Meta{Name: "SizeVariations"})
	// SizeVariation := SizeVariationMeta.Resource
	// SizeVariation.IndexAttrs("ID", "Size", "AvailableQuantity")
	// SizeVariation.ShowAttrs("ID", "Size", "AvailableQuantity")

	API.AddResource(&orders.Order{})

	// API.AddResource(&users.User{})
	user := API.AddResource(&users.User{})
	user.AddSubResource("Orders")
	// afterSales, _ := User.AddSubResource("Aftersales")
	user.AddSubResource("Aftersales")
	// // userOrders.AddSubResource("OrderItems", &admin.Config{Name: "Items"})

	// API.AddResource(&products.Category{})

	API.AddResource(&aftersales.Aftersale{})

	application.Router.Mount(app.Config.Prefix, API.NewServeMux(app.Config.Prefix))
}
