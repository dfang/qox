package users

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/qor/media"
	"github.com/qor/media/oss"
)

var ROLES_VALUES = []string{"admin", "operator", "workman"}
var ROLES_TEXTS = []string{"管理员", "调度员", "师傅"}
var GENDERS_VALUES = []string{"male", "female"}
var GENDERS_TEXTS = []string{"男", "女"}

type User struct {
	gorm.Model
	Email    string `form:"email"`
	Password string
	Name     string `form:"name"`
	Gender   string
	Role     string
	Birthday *time.Time

	// 身份证号码
	IdentityCardNum string
	// 手机号码
	MobilePhone string
	// 车牌号码
	CarLicencePlateNum string
	// 车型 东风小货
	CarType string
	// 驾照类型 C1
	CarLicenseType string
	// 是否临时工
	IsCasual bool

	JDAppUser string
	HireDate  *time.Time

	// UserType DeliveryMan、SetupMan
	// Type string

	Balance                float32
	DefaultBillingAddress  uint `form:"default-billing-address"`
	DefaultShippingAddress uint `form:"default-shipping-address"`
	Addresses              []Address
	Avatar                 AvatarImageStorage

	// Confirm
	ConfirmToken string
	Confirmed    bool

	// Recover
	RecoverToken       string
	RecoverTokenExpiry *time.Time

	// Accepts
	AcceptPrivate bool `form:"accept-private"`
	AcceptLicense bool `form:"accept-license"`
	AcceptNews    bool `form:"accept-news"`
}

func (user User) DisplayName() string {
	return user.Email
}

func (user User) AvailableLocales() []string {
	return []string{"en-US", "zh-CN"}
}

type AvatarImageStorage struct{ oss.OSS }

func (AvatarImageStorage) GetSizes() map[string]*media.Size {
	return map[string]*media.Size{
		"small":  {Width: 50, Height: 50},
		"middle": {Width: 120, Height: 120},
		"big":    {Width: 320, Height: 320},
	}
}

// T2V Text To Value
func T2V(tArr []string, vArr []string, t string) string {
	var index int
	for i, item := range tArr {
		if t == item {
			index = i
			break
		}
	}
	return vArr[index]
}

// V2T Value To Text
func V2T(vArr []string, tArr []string, t string) string {
	var index int
	for i, item := range vArr {
		if t == item {
			index = i
			break
		}
	}
	return tArr[index]
}
