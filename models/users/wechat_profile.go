package users

import (
	"github.com/jinzhu/gorm"
)

// WechatProfile 微信资料
type WechatProfile struct {
	gorm.Model

	Openid     string `json:"openid"`
	Unionid    string `json:"unionid"`
	Nickname   string `json:"nickname"`
	Sex        int    `json:"sex"`
	City       string `json:"city"`
	Province   string `json:"province"`
	Country    string `json:"country"`
	Headimgurl string `json:"headimgurl"`

	MobilePhone string `json:"mobile_phone"`
	Role        string
}
