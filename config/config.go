package config

import (
	"log"
	"os"

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
	"github.com/qor/media/oss"

	// "github.com/qor/oss/qiniu"
	"github.com/qor/oss/s3"
	"github.com/qor/redirect_back"
	"github.com/qor/session/manager"
	"github.com/unrolled/render"
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
	DB    struct {
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
)

func init() {

	if err := configor.Load(&Config, "config/database.yml", "config/smtp.yml", "config/application.yml"); err != nil {
		panic(err)
	}

	location.GoogleAPIKey = Config.GoogleAPIKey
	location.BaiduAPIKey = Config.BaiduAPIKey

	log.Println(Config.Qiniu)

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

	if Config.S3.AccessKeyID != "" {
		oss.Storage = s3.New(&s3.Config{
			AccessID:  Config.S3.AccessKeyID,
			AccessKey: Config.S3.SecretAccessKey,
			Region:    Config.S3.Region,
			Bucket:    Config.S3.S3Bucket,
		})
	}

	// oss.Storage = qiniu.New(&qiniu.Config{
	// 	AccessID:  Config.Qiniu.AccessID,
	// 	AccessKey: Config.Qiniu.AccessKey,
	// 	Bucket:    Config.Qiniu.Bucket,
	// 	Region:    Config.Qiniu.Region,
	// 	Endpoint:  Config.Qiniu.Endpoint,
	// })

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
