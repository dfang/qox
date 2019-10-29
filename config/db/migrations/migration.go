package migrations

import (
	"fmt"

	"github.com/dfang/qor-demo/app/admin"
	"github.com/dfang/qor-demo/config/db"
	"github.com/dfang/qor-demo/models/aftersales"
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
	"github.com/qor/notification"
	"github.com/qor/transition"
)

// Migrate Run Migration
func Migrate() {
	fmt.Println("running migration .......")

	AutoMigrate(&aftersales.Aftersale{})
	AutoMigrate(&aftersales.Manufacturer{})
	AutoMigrate(&aftersales.Settlement{})
	AutoMigrate(&aftersales.Balance{})

	AutoMigrate(&products.Product{}, &products.ProductVariation{}, &products.ProductImage{}, &products.ColorVariation{}, &products.ColorVariationImage{}, &products.SizeVariation{})
	AutoMigrate(&products.Color{}, &products.Size{}, &products.Material{}, &products.Category{}, &products.Collection{})

	AutoMigrate(&users.User{}, &users.Address{})

	AutoMigrate(&users.WechatProfile{})

	AutoMigrate(&auth_identity.AuthIdentity{})

	AutoMigrate(&orders.Order{}, &orders.OrderItem{}, &orders.OrderFollowUp{})

	AutoMigrate(&orders.DeliveryMethod{})

	AutoMigrate(&stores.Store{})

	AutoMigrate(&notification.QorNotification{})
	AutoMigrate(&i18n_database.Translation{})
	AutoMigrate(&transition.StateChangeLog{})
	AutoMigrate(&activity.QorActivity{})

	AutoMigrate(&settings.Setting{}, &settings.MediaLibrary{}, &settings.Brand{}, &settings.ServiceType{})
	AutoMigrate(&asset_manager.AssetManager{})
	AutoMigrate(&admin.QorWidgetSetting{})
	AutoMigrate(&banner_editor.QorBannerEditorSetting{})
	AutoMigrate(&seo.MySEOSetting{})

	AutoMigrate(&blogs.Page{}, &blogs.Article{})
	AutoMigrate(&help.QorHelpEntry{})
}

// AutoMigrate run auto migration
func AutoMigrate(values ...interface{}) {
	for _, value := range values {
		db.DB.AutoMigrate(value)
	}
}
