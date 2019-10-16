package db

import (
	"errors"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"

	// _ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	// _ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/dfang/qor-demo/config"
	// "github.com/dfang/qor-demo/config/logger"
	"github.com/qor/audited"
	"github.com/qor/l10n"
	"github.com/qor/media"
	"github.com/qor/publish2"
	"github.com/qor/sorting"
	"github.com/qor/validations"
	"github.com/rs/zerolog/log"

	"github.com/dfang/qor-demo/models/aftersales"
)

// DB Global DB connection
var DB *gorm.DB

func init() {
	var err error

	// logger, _ := zap.NewDevelopment()
	// logger.Info("Hello zap", zap.String("key", "value"), zap.Time("now", time.Now()))
	// logger.Logger.Info("Hello zap", zap.String("key", "value"), zap.Time("now", time.Now()))

	dbConfig := config.Config.DB
	// Logger, _ := zap.NewDevelopment()
	// zap.ReplaceGlobals(Logger)
	// // zap.L().Debug("Connection info from config/database.yml or env",
	// // 	zap.String("adapter", dbConfig.Adapter),
	// // 	zap.String("host", dbConfig.Host),
	// // 	zap.String("db", dbConfig.Name),
	// // 	zap.String("port", dbConfig.Port),
	// // 	zap.String("user", dbConfig.User),
	// // 	zap.String("password", dbConfig.Password),
	// // )
	// zap.S().Debugw("database connection info",
	// 	"adapter", dbConfig.Adapter,
	// 	"host", dbConfig.Host,
	// 	"db", dbConfig.Name,
	// 	"port", dbConfig.Port,
	// 	"user", dbConfig.User,
	// 	"password", dbConfig.Password,
	// )

	if config.Config.DB.Adapter == "mysql" {
		log.Debug().Msg(fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name))
		DB, err = gorm.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name))
		// DB = DB.Set("gorm:table_options", "CHARSET=utf8")
	} else if config.Config.DB.Adapter == "postgres" {
		log.Debug().Msg(fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name))
		DB, err = gorm.Open("postgres", fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name))
	} else if config.Config.DB.Adapter == "sqlite" {
		log.Debug().Msg(fmt.Sprintf("%v/%v", os.TempDir(), dbConfig.Name))
		DB, err = gorm.Open("sqlite3", fmt.Sprintf("%v/%v", os.TempDir(), dbConfig.Name))
	} else {
		panic(errors.New("not supported database adapter"))
	}

	if err == nil {
		// if os.Getenv("DEBUG") != "" {
		// 	DB.LogMode(true)
		// }

		// l10n.Global = "zh-CN"

		l10n.RegisterCallbacks(DB)
		sorting.RegisterCallbacks(DB)
		validations.RegisterCallbacks(DB)
		audited.RegisterCallbacks(DB)
		media.RegisterCallbacks(DB)
		publish2.RegisterCallbacks(DB)

		aftersales.RegisterCallbacks(DB)
	} else {
		panic(err)
	}
}

// func GetDB() *gorm.DB {

// }
