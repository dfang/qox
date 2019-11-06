package aftersales

import (
	"fmt"
	"math"
	"strconv"

	"github.com/dfang/qor-demo/models/users"
	"github.com/gocraft/work"
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
	fmt.Println("before save ...")
	fmt.Println(item)
	fmt.Println(item.User)
	fmt.Println(item.UserID) // 0, 取不到值得, 不对
	fmt.Println(item.User.ID)

	// 表单提交来的 item.User.ID 有值
	// AutoWithdraw 来的只有item.UserID 有值
	var userID uint
	if item.UserID != 0 {
		userID = item.UserID
	} else {
		userID = item.User.ID
	}

	if item.Direction == "罚款" {
		if item.Amount > 0 {
			item.Amount = -item.Amount
		}
		item.State = "free"
	}

	if item.Direction == "奖励" {
		if item.Amount > 0 {
			item.State = "free"
		}
	}

	// 检查是否可提现
	if item.Direction == "提现" && userID > 0 {
		// s, err := strconv.ParseUint(fmt.Sprint(item.Amount), 10, 64)
		// if err != nil {
		// 	panic(err)
		// }
		// amount := float32(s)

		balance := UpdateBalanceFor(fmt.Sprint(userID))

		fmt.Println(*balance)

		fmt.Println("剩余可提现额度是", balance.FreeAmount)
		fmt.Println("尝试提现", item.Amount)

		if math.Abs(float64(item.Amount)) >= float64(balance.FreeAmount) {
			return fmt.Errorf("提现不能超过可提现额度, 提现失败")
		}

		if item.Amount > 0 {
			item.Amount = -item.Amount
		}

		item.State = "withdrawed"
	}

	return nil
}

// AfterSave AfterSave Callback
func (item *Settlement) AfterSave(scope *gorm.Scope) error {
	var enqueuer = work.NewEnqueuer("qor", db.RedisPool)
	fmt.Printf("enqueueing update_balance for user id %d .....\n", item.UserID)
	id := strconv.FormatUint(uint64(item.User.ID), 10)
	j, err := enqueuer.Enqueue("update_balance", work.Q{"user_id": id})
	if err != nil {
		return err
	}
	fmt.Printf("Job ID: %s, Name: %s, EnqueueAt: %d\n", j.ID, j.Name, j.EnqueuedAt)

	return nil
}

// UpdateBalanceFor 更新账户余额
// 罚款和奖励都是立即生效的（立即变为free状态的）
// update balance by user_id
func UpdateBalanceFor(userID string) *Balance {
	var results []Result
	// var f1, f2, f3 float32

	var awards, fines, incomeFrozen, incomeFree, withdraw float32

	var balance Balance

	u64, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		panic(err)
	}

	db.DB.Model(Balance{}).Where("user_id = ?", userID).Assign(Balance{UserID: uint(u64)}).FirstOrInit(&balance)

	fmt.Println("balance is ....")
	fmt.Println(balance)

	// select state, direction, sum(amount) as total from settlements where user_id = 1 group by state, direction;
	db.DB.Table("settlements").Select("state, direction, sum(amount) as total").Group("state, direction").Where("user_id = ?", userID).Scan(&results)
	for _, i := range results {
		// fmt.Println(i.State)
		// fmt.Println(i.Total)

		if i.Direction == "收入" && i.State == "frozen" {
			incomeFrozen = i.Total
		}

		if i.Direction == "收入" && i.State == "free" {
			incomeFree = i.Total
		}

		if i.Direction == "奖励" {
			awards = i.Total
		}

		if i.Direction == "罚款" {
			fines = i.Total
		}

		if i.State == "withdrawed" {
			withdraw = i.Total
		}

	}

	balance.TotalAmount = incomeFrozen + incomeFree
	balance.WithdrawAmount = withdraw
	balance.FineAmount = fines
	balance.AwardAmount = awards

	balance.FrozenAmount = incomeFrozen
	balance.FreeAmount = incomeFrozen + incomeFree + awards + (fines + withdraw)

	db.DB.Save(&balance)
	return &balance
}

// Result 临时结果表
type Result struct {
	State     string
	Direction string
	Total     float32
}
