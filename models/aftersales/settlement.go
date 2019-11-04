package aftersales

import (
	"github.com/dfang/qor-demo/models/users"
	"github.com/jinzhu/gorm"
	"github.com/qor/audited"
	"github.com/qor/transition"

	"github.com/dfang/qor-demo/config/db"
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
			item.State = "free"
		}
	}

	if item.Direction == "奖励" {
		if item.Amount > 0 {
			item.State = "free"
		}
	}

	// 检查是否可提现
	if item.Direction == "提现" {
		// Amount <= Balance.FreeAmount
	}

	return nil
}

// UpdateBalanceFor 更新账户余额
// 罚款和奖励都是立即生效的（立即变为free状态的）
func UpdateBalanceFor(userID uint) Balance {
	var results []Result
	var f1, f2, f3 float32

	// update balance by user_id
	var balance Balance
	db.DB.Model(Balance{}).Where("user_id = ?", userID).Assign(Balance{UserID: userID}).FirstOrInit(&balance)

	// select state, sum(amount) as total from settlements where user_id = 1 group by state;
	db.DB.Table("settlements").Select("state, sum(amount) as total").Group("state").Where("user_id = ?", userID).Scan(&results)
	for _, i := range results {
		// fmt.Println(i.State)
		// fmt.Println(i.Total)
		if i.State == "frozen" {
			f1 = i.Total
		}

		if i.State == "free" {
			f2 = i.Total
		}

		if i.State == "withdrawed" {
			f3 = i.Total
		}
	}

	balance.FrozenAmount = f1
	balance.FreeAmount = f2 + f3
	balance.WithdrawAmount = f3
	balance.TotalAmount = f2 + f1

	return balance
}

type Result struct {
	State string
	Total float32
}
