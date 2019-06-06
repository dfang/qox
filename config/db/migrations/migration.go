package migrations

import (
	"fmt"

	"github.com/dfang/qor-demo/app/admin"
	"github.com/dfang/qor-demo/config/db"
	"github.com/dfang/qor-demo/models/blogs"
	"github.com/dfang/qor-demo/models/orders"
	"github.com/dfang/qor-demo/models/products"
	"github.com/dfang/qor-demo/models/seo"
	"github.com/dfang/qor-demo/models/settings"
	"github.com/dfang/qor-demo/models/stores"
	"github.com/dfang/qor-demo/models/users"
	"github.com/qor/activity"
	"github.com/qor/auth/auth_identity"
	"github.com/qor/banner_editor"
	"github.com/qor/help"
	i18n_database "github.com/qor/i18n/backends/database"
	"github.com/qor/media/asset_manager"
	"github.com/qor/transition"
)

func Migrate() {
	fmt.Println("running miration .......")

	AutoMigrate(&asset_manager.AssetManager{})

	AutoMigrate(&products.Product{}, &products.ProductVariation{}, &products.ProductImage{}, &products.ColorVariation{}, &products.ColorVariationImage{}, &products.SizeVariation{})
	AutoMigrate(&products.Color{}, &products.Size{}, &products.Material{}, &products.Category{}, &products.Collection{})

	AutoMigrate(&users.User{}, &users.Address{})

	AutoMigrate(&orders.Order{}, &orders.OrderItem{})

	AutoMigrate(&orders.DeliveryMethod{})

	AutoMigrate(&stores.Store{})

	AutoMigrate(&settings.Setting{}, &settings.MediaLibrary{})

	AutoMigrate(&transition.StateChangeLog{})

	AutoMigrate(&activity.QorActivity{})

	AutoMigrate(&admin.QorWidgetSetting{})

	AutoMigrate(&blogs.Page{}, &blogs.Article{})

	AutoMigrate(&seo.MySEOSetting{})

	AutoMigrate(&help.QorHelpEntry{})

	AutoMigrate(&auth_identity.AuthIdentity{})

	AutoMigrate(&banner_editor.QorBannerEditorSetting{})

	AutoMigrate(&i18n_database.Translation{})
}

// AutoMigrate run auto migration
func AutoMigrate(values ...interface{}) {
	for _, value := range values {
		db.DB.AutoMigrate(value)
	}
}
