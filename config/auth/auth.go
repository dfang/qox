package auth

import "github.com/dfang/qor-demo/config/bindatafs"

import (
	"time"

	"github.com/dfang/qor-demo/config"
	"github.com/dfang/qor-demo/config/db"
	"github.com/dfang/qor-demo/models/users"
	"github.com/dfang/auth"
	"github.com/dfang/auth/authority"
	"github.com/dfang/auth_themes/clean"
	"github.com/qor/render"
)

var (
	// Auth initialize Auth for Authentication
	Auth = clean.New(&auth.Config{
		DB:         db.DB,
		Mailer:     config.Mailer,
		Render:     render.New(&render.Config{AssetFileSystem: bindatafs.AssetFS.NameSpace("auth")}),
		UserModel:  users.User{},
		Redirector: auth.Redirector{RedirectBack: config.RedirectBack},
	})

	// Authority initialize Authority for Authorization
	Authority = authority.New(&authority.Config{
		Auth: Auth,
	})
)

func init() {
	// Auth.RegisterProvider(github.New(&config.Config.Github))
	// Auth.RegisterProvider(google.New(&config.Config.Google))
	// Auth.RegisterProvider(facebook.New(&config.Config.Facebook))
	// Auth.RegisterProvider(twitter.New(&config.Config.Twitter))

	Authority.Register("logged_in_half_hour", authority.Rule{TimeoutSinceLastLogin: time.Minute * 30})
}
