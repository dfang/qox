package config

import (
	"fmt"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/jinzhu/configor"
	"github.com/qor/auth/providers/facebook"
	"github.com/qor/auth/providers/github"
	"github.com/qor/auth/providers/google"
	"github.com/qor/auth/providers/twitter"

	// amazonpay "github.com/qor/amazon-pay-sdk-go"
	"github.com/qor/gomerchant"
	"github.com/qor/location"
	"github.com/qor/mailer"
	"github.com/qor/mailer/logger"
	// "github.com/qor/media/oss"

	// "github.com/qor/oss/qiniu"
	// "github.com/qor/oss/s3"
	"github.com/qor/redirect_back"
	"github.com/qor/session/manager"
	"github.com/unrolled/render"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

var Config = struct {
	HTTPS bool `default:"false" env:"HTTPS"`
	Port  uint `default:"7000" env:"PORT"`
	Redis struct {
		Host string `env:"REDIS_HOST" default:"localhost"`
		Port string `env:"REDIS_PORT" default:"6379"`
	}
	DB struct {
		Name     string `env:"DBName" default:"qor_example"`
		Adapter  string `env:"DBAdapter" default:"postgres"`
		Host     string `env:"DBHost" default:"localhost"`
		Port     string `env:"DBPort" default:"5432"`
		User     string `env:"DBUser"`
		Password string `env:"DBPassword"`
	}
	S3 struct {
		AccessKeyID     string `env:"QOR_AWS_ACCESS_KEY_ID"`
		SecretAccessKey string `env:"QOR_AWS_SECRET_ACCESS_KEY"`
		Region          string `env:"QOR_AWS_REGION"`
		S3Bucket        string `env:"QOR_AWS_BUCKET"`
	}

	Qiniu struct {
		AccessID  string `env:"QOR_QINIU_ACCESS_ID"`
		AccessKey string `env:"QOR_QINIU_ACCESS_KEY"`
		Bucket    string `env:"QOR_QINIU_BUCKET"`
		Region    string `env:"QOR_QINIU_REGION"`
		Endpoint  string `env:"QOR_QINIU_ENDPOINT"`
	}

	AmazonPay struct {
		MerchantID   string `env:"AmazonPayMerchantID"`
		AccessKey    string `env:"AmazonPayAccessKey"`
		SecretKey    string `env:"AmazonPaySecretKey"`
		ClientID     string `env:"AmazonPayClientID"`
		ClientSecret string `env:"AmazonPayClientSecret"`
		Sandbox      bool   `env:"AmazonPaySandbox"`
		CurrencyCode string `env:"AmazonPayCurrencyCode" default:"JPY"`
	}

	Cron struct {
		ExpireAftersales         string `env:"EXPIRE_AFTERSALES" default:"*/60 * * * * *"`
		FreezeAuditedAftersales  string `env:"FREEZE_AUDITED_AFTERSALES" default:"0 */2 * * * *"`
		UnfreezeAftersales       string `env:"UNFREEZE_AFTERSALES" default:"0 */5 * * * *`
		UpdateBalances           string `env:"UPDATE_BALANCES" default:"0 */5 * * * *`
		AutoExportMobilePhones   string `env:"AUTO_EXPORT_MOBILE_PHONES" default:"0 0 1 * * *"`
		AutoExportOrderDetails   string `env:"AUTO_EXPORT_ORDER_DETAILS" default:"0 0 1 * * *"`
		AutoExportOrderFollowUps string `env:"AUTO_EXPORT_ORDER_FOLLOWUPS" default:"0 0 1 * * *"`
		AutoExportOrderFees      string `env:"AUTO_EXPORT_ORDER_FEES" default:"0 */1 * * * *"`
		AutoUpdateOrderItems     string `env:"AUTO_UPDATE_ORDER_ITEMS" default:"0 */5 * * * *"`
		AutoDeliverOrders        string `env:"AUTO_DELIVERY_ORDERS" default:"0 0 1 * *"`

		// DEMO_MODE = true 才生效
		AutoInquire            string `env:"AutoInquire", default:"*/30 * * * * *"`
		AutoSchedule           string `env:"AutoSchedule", default:"0 */5 * * * *"`
		AutoProcess            string `env:"AutoProcess", default:"0 */2 * * * *"`
		AutoFinish             string `env:"AutoFinish", default:"0 */2 * * * *"`
		AutoAudit              string `env:"AutoAudit", default:"0 */1 * * * *"`
		AutoWithdraw           string `env:"AutoWithdraw", default:"0 */5 * * * *"`
		AutoAward              string `env:"AutoAward", default:"0 */6 * * * *"`
		AutoFine               string `env:"AutoFine", default:"0 */7 * * * *"`
		AutoGenerateAftersales string `env:"AutoGenerateAftersales", default:"0 */30 * * * *"`
	}

	SMTP         SMTPConfig
	Github       github.Config
	Google       google.Config
	Facebook     facebook.Config
	Twitter      twitter.Config
	GoogleAPIKey string `env:"GoogleAPIKey"`
	BaiduAPIKey  string `env:"BaiduAPIKey"`
}{}

var (
	// Root           = os.Getenv("GOPATH") + "/src/github.com/dfang/qor-demo"
	Root, _ = os.Getwd()
	Mailer  *mailer.Mailer
	Render  = render.New()
	// AmazonPay      amazonpay.AmazonPayService
	PaymentGateway gomerchant.PaymentGateway
	RedirectBack   = redirect_back.New(&redirect_back.Config{
		SessionManager:  manager.SessionManager,
		IgnoredPrefixes: []string{"/auth"},
	})

	StartUpStartTime *time.Time
)

// Initialize changed init to Initialize
func Initialize() {
	start := time.Now()
	StartUpStartTime = &start
	fmt.Println("STARTUP begins at", start.Format("2006-01-02 15:04:05"))

	if os.Getenv("DEBUG") != "true" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// // Auto Reload
	if err := configor.New(&configor.Config{AutoReload: true, AutoReloadInterval: time.Minute, AutoReloadCallback: func(config interface{}) {
		fmt.Printf("%v changed", config)
	}}).Load(&Config, "config/database.yml", "config/smtp.yml", "config/application.yml"); err != nil {
		// if err := configor.Load(&Config, "config/database.yml", "config/smtp.yml", "config/application.yml"); err != nil {
		panic(err)
	}

	// if Config.Cron.ExpireAftersales == "" {
	// 	// every 60 seconds
	// 	Config.Cron.ExpireAftersales = "*/60 * * * * *"
	// }

	// if Config.Cron.FreezeAuditedAftersales == "" {
	// 	// every 2 minutes
	// 	Config.Cron.FreezeAuditedAftersales = "0 */2 * * * *"
	// }

	// if Config.Cron.UnfreezeAftersales == "" {
	// 	// every 5 minutes
	// 	Config.Cron.UnfreezeAftersales = "0 */5 * * * *"
	// }

	// if Config.Cron.UpdateBalances == "" {
	// 	// every 5 minutes
	// 	Config.Cron.UpdateBalances = "0 */5 * * * *"
	// }

	// if Config.Cron.AutoExportMobilePhones == "" {
	// 	// 1:00 AM every day
	// 	Config.Cron.AutoExportMobilePhones = "0 0 1 * * *"
	// }

	log.Debug().Msg("Cron settings: ")
	log.Debug().Msgf("Cron settings: %+v", Config.Cron)
	// log.Debug().Msgf("ExpireAftersales: %s", Config.Cron.ExpireAftersales)
	// log.Debug().Msgf("FreezeAuditedAftersales: %s", Config.Cron.FreezeAuditedAftersales)
	// log.Debug().Msgf("UnfreezeAftersales: %s", Config.Cron.UnfreezeAftersales)
	// log.Debug().Msgf("UpdateBalances: %s", Config.Cron.UpdateBalances)
	// log.Debug().Msgf("AutoExportMobilePhones: %s", Config.Cron.AutoExportMobilePhones)

	location.GoogleAPIKey = Config.GoogleAPIKey
	location.BaiduAPIKey = Config.BaiduAPIKey

	// log.Println(Config.Qiniu)

	// if Config.S3.AccessKeyID == "" {
	// 	log.Println("Please set env QOR_AWS_ACCESS_KEY_ID")
	// 	os.Exit(1)
	// }

	// if Config.S3.SecretAccessKey == "" {
	// 	log.Println("Please set env QOR_AWS_SECRET_ACCESS_KEY")
	// 	os.Exit(1)
	// }

	// if Config.S3.Region == "" {
	// 	log.Println("Please set env QOR_AWS_REGION")
	// 	os.Exit(1)
	// }

	// if Config.S3.S3Bucket == "" {
	// 	log.Println("Please set env QOR_AWS_BUCKET")
	// 	os.Exit(1)
	// }

	// if Config.S3.AccessKeyID != "" {
	// 	oss.Storage = s3.New(&s3.Config{
	// 		AccessID:  Config.S3.AccessKeyID,
	// 		AccessKey: Config.S3.SecretAccessKey,
	// 		Region:    Config.S3.Region,
	// 		Bucket:    Config.S3.S3Bucket,
	// 	})
	// }

	// if Config.Qiniu.AccessID != "" {
	// 	oss.Storage = qiniu.New(&qiniu.Config{
	// 		AccessID:  Config.Qiniu.AccessID,
	// 		AccessKey: Config.Qiniu.AccessKey,
	// 		Bucket:    Config.Qiniu.Bucket,
	// 		Region:    Config.Qiniu.Region,
	// 		Endpoint:  Config.Qiniu.Endpoint,
	// 	})
	// }

	// AmazonPay = amazonpay.New(&amazonpay.Config{
	// 	MerchantID: Config.AmazonPay.MerchantID,
	// 	AccessKey:  Config.AmazonPay.AccessKey,
	// 	SecretKey:  Config.AmazonPay.SecretKey,
	// 	Sandbox:    true,
	// 	Region:     "jp",
	// })

	// dialer := gomail.NewDialer(Config.SMTP.Host, Config.SMTP.Port, Config.SMTP.User, Config.SMTP.Password)
	// sender, err := dialer.Dial()

	// Mailer = mailer.New(&mailer.Config{
	// 	Sender: gomailer.New(&gomailer.Config{Sender: sender}),
	// })
	Mailer = mailer.New(&mailer.Config{
		Sender: logger.New(&logger.Config{}),
	})
}
