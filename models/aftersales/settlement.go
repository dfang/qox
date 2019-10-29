package aftersales

import (
	"github.com/dfang/qor-demo/models/users"
	"github.com/jinzhu/gorm"
	"github.com/qor/audited"
	"github.com/qor/transition"
)

// Settlement Aftersale 售后厂家管理
type Settlement struct {
	gorm.Model

	UserID uint
	User   users.User

	Amount    float32
	Direction string

	AftersaleID uint
	Aftersale   Aftersale

	transition.Transition
	audited.AuditedModel

	// deposit 存入（完成服务）
	// withdraw 提现 ()
	// Direction string
}

// BeforeCreate 初始状态
func (item *Settlement) BeforeCreate(scope *gorm.Scope) error {
	if item.Direction != "提现" {
		scope.SetColumn("State", "frozen")
	}
	return nil
}

// BeforeSave BeforeSave Callback
func (item *Settlement) BeforeSave(scope *gorm.Scope) error {
	if item.Direction == "提现" {
		if item.Amount > 0 {
			item.Amount = -item.Amount
			item.State = "withdrawed"
		}
	}

	// 检查是否可提现
	if item.Direction == "提现" {
		// Amount <= Balance.FreeAmount
	}

	return nil
}
