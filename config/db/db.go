package db

import (
	"errors"
	"fmt"
	"os"

	"github.com/gomodule/redigo/redis"
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
)

// DB Global DB connection
var DB *gorm.DB

// RedisPool Workpool connection
var RedisPool *redis.Pool

// Initialize changed init to Initialize
func Initialize() error {
	var err error

	// Make a redis pool
	RedisPool = &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			log.Debug().Msg(fmt.Sprintf("%s:%s", config.Config.Redis.Host, config.Config.Redis.Port))
			return redis.Dial("tcp", fmt.Sprintf("%s:%s", config.Config.Redis.Host, config.Config.Redis.Port))
		},
	}

	dbConfig := config.Config.DB

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

	if err := DB.DB().Ping(); err != nil {
		log.Panic().Msgf("database/sql db.ping() failed with err: %s", err.Error())
	}

	if err == nil {
		if os.Getenv("DEBUG") != "false" && os.Getenv("DEBUG_GORM_LOG_SQL") == "true" {
			DB.LogMode(true)
		}

		// l10n.Global = "zh-CN"

		l10n.RegisterCallbacks(DB)
		sorting.RegisterCallbacks(DB)
		validations.RegisterCallbacks(DB)
		audited.RegisterCallbacks(DB)
		media.RegisterCallbacks(DB)
		publish2.RegisterCallbacks(DB)

		// aftersales.RegisterCallbacks(DB)
	} else {
		panic(err)
	}

	return nil
}

// func init() {
// 	config.Initialize()
// 	Initialize()
// }
