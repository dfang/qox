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

	Direction string
	Amount    float32

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

	if item.Direction == "罚款" {
		if item.Amount > 0 {
			item.Amount = -item.Amount
		}
	}
	// 检查是否可提现
	if item.Direction == "提现" {
		// Amount <= Balance.FreeAmount
	}

	return nil
}
