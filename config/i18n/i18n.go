package i18n

import (
	"path/filepath"

	"github.com/qor/i18n"
	"github.com/qor/i18n/backends/database"
	"github.com/qor/i18n/backends/yaml"

	"github.com/dfang/qor-demo/config"
	"github.com/dfang/qor-demo/config/db"
	"github.com/rs/zerolog/log"
)

// I18n i18n.I18n
var I18n *i18n.I18n

// Initialize changed init to Initialize
func Initialize() {
	// I18n will look up the translation in order
	// I18n = i18n.New(database.New(db.DB), yaml.New(filepath.Join(config.Root, "config/locales")))
	log.Debug().Str("locales", "locales").Msg(filepath.Join(config.Root, "config/locales"))

	I18n = i18n.New(
		yaml.New(filepath.Join(config.Root, "config/locales")),
		database.New(db.DB),
	)
	// I18n.AddTranslation(&i18n.Translation{Key: "qor_admin.menus.Dashboard", Locale: "en-US", Value: "dashbord"})
	// I18n.AddTranslation(&i18n.Translation{Key: "qor_admin.menus.Dashboard", Locale: "zh-CN", Value: "控制面板"})
}
