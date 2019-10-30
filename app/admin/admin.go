package admin

import (
	"github.com/dfang/qor-demo/config/i18n"
	"github.com/dfang/qor-demo/models/settings"
	"github.com/qor/action_bar"
	"github.com/qor/admin"
	"github.com/qor/application"
	"github.com/qor/help"
	"github.com/qor/media/asset_manager"
	"github.com/qor/media/media_library"
	"github.com/qor/roles"
)

// ActionBar admin action bar
var ActionBar *action_bar.ActionBar

// AssetManager asset manager
var AssetManager *admin.Resource

// New new home app
func New(config *Config) *App {
	if config.Prefix == "" {
		config.Prefix = "/admin"
	}
	return &App{Config: config}
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
	Admin := application.Admin

	AssetManager = Admin.AddResource(&asset_manager.AssetManager{}, &admin.Config{Invisible: true})

	// Add Media Library
	Admin.AddResource(&media_library.MediaLibrary{}, &admin.Config{Menu: []string{"Site Management"}})

	// Add Help
	Help := Admin.NewResource(&help.QorHelpEntry{})
	Help.Meta(&admin.Meta{Name: "Body", Config: &admin.RichEditorConfig{AssetManager: AssetManager}})

	// Add action bar
	ActionBar = action_bar.New(Admin)
	ActionBar.RegisterAction(&action_bar.Action{Name: "Admin Dashboard", Link: "/admin"})

	// i18n.Load()
	// Add Translations
	Admin.AddResource(i18n.I18n, &admin.Config{Menu: []string{"Site Management"}, Priority: -1})

	// Add Setting
	Admin.AddResource(&settings.Setting{}, &admin.Config{Name: "Shop Setting", Menu: []string{"Site Management"}, Singleton: true, Priority: 1})

	// Add Help
	Admin.AddResource(&help.QorHelpEntry{}, &admin.Config{Name: "Help", Menu: []string{"Site Management"}, Singleton: true, Priority: 1})

	// 一级菜单
	// https://doc.getqor.com/admin/authentication.html#authorization
	// Hide some menus for operator role
	Admin.AddMenu(&admin.Menu{Name: "Aftersale Management", Priority: 2})
	Admin.AddMenu(&admin.Menu{Name: "Settlement Management", Priority: 3})
	Admin.AddMenu(&admin.Menu{Name: "User Management", Priority: 3})
	Admin.AddMenu(&admin.Menu{Name: "Order Management", Priority: 4, Permission: roles.Deny(roles.Read, "operator")})
	Admin.AddMenu(&admin.Menu{Name: "Product Management", Priority: 5, Permission: roles.Deny(roles.Read, "operator")})
	Admin.AddMenu(&admin.Menu{Name: "Pages Management", Priority: 6, Permission: roles.Deny(roles.Read, "operator")})
	Admin.AddMenu(&admin.Menu{Name: "Site Management", Priority: 7, Permission: roles.Deny(roles.Read, "operator")})
	// publish2
	Admin.AddMenu(&admin.Menu{Name: "Publishing", Priority: 8, Permission: roles.Deny(roles.Read, "operator")})

	SetupNotification(Admin)
	SetupWorker(Admin)
	SetupSEO(Admin)
	SetupWidget(Admin)
	SetupDashboard(Admin)
	application.Router.Mount(app.Config.Prefix, Admin.NewServeMux(app.Config.Prefix))

	// TODO: Investigate why Admin.MountTo not works
	// mux := http.NewServeMux()
	// Admin.MountTo("/admin", mux)
	// Admin.MountTo("/admin", application.Router)
	// application.Router.Mount("/admin", Admin.NewServeMux("/admin"))
}
