package auth

import (
	"os"
	"strconv"
	"time"

	"github.com/dfang/qor-demo/config"
	"github.com/dfang/qor-demo/config/bindatafs"
	"github.com/dfang/qor-demo/config/db"
	"github.com/dfang/qor-demo/models/users"
	"github.com/qor/auth"
	"github.com/qor/auth/authority"
	"github.com/qor/auth/providers/wechat_work"
	"github.com/qor/auth_themes/clean"
	"github.com/qor/render"
)

// 修改密码
// 如何直接改默认管理员或者其他用户的密码
// https://play.golang.org/p/fEefNQ48-L9
var (
	Auth      *auth.Auth
	Authority *authority.Authority
)

// Initialize changed init to Initialize
func Initialize() {

	// Auth initialize Auth for Authentication
	Auth = clean.New(&auth.Config{
		DB:         db.DB,
		Mailer:     config.Mailer,
		Render:     render.New(&render.Config{AssetFileSystem: bindatafs.AssetFS.NameSpace("auth")}),
		UserModel:  users.User{},
		Redirector: auth.Redirector{RedirectBack: config.RedirectBack},
		// Redirector: AdminRedirector{},
	})

	// Authority initialize Authority for Authorization
	Authority = authority.New(&authority.Config{
		Auth: Auth,
	})

	// Auth.RegisterProvider(github.New(&config.Config.Github))
	// Auth.RegisterProvider(google.New(&config.Config.Google))
	// Auth.RegisterProvider(facebook.New(&config.Config.Facebook))
	// Auth.RegisterProvider(twitter.New(&config.Config.Twitter))
	CorpID := os.Getenv("AUTH_CORP_ID")
	CorpSecret := os.Getenv("AUTH_CORP_SECRET")
	if CorpID == "" || CorpSecret == "" {
		panic("AUTH_CORP_ID 和 AUTH_CORP_SECRET 都不能为空, 请配置")
	}

	AUTH_DOMAIN := os.Getenv("AUTH_DOMAIN")
	if AUTH_DOMAIN == "" {
		panic("请设置用于企业微信登录后台的 AUTH_DOMAIN")
	}
	redirectURI := AUTH_DOMAIN + "/wechat_work/callback"

	AgentID, err := strconv.ParseInt(os.Getenv("AUTH_AGENT_ID"), 10, 64)
	if err != nil {
		panic(err)
	}

	Auth.RegisterProvider(wechat_work.New(&wechat_work.Config{
		CorpID:      CorpID,
		CorpSecret:  CorpSecret,
		AgentID:     AgentID,
		RedirectURI: redirectURI,
	}))

	Authority.Register("logged_in_half_hour", authority.Rule{TimeoutSinceLastLogin: time.Minute * 30})
}
