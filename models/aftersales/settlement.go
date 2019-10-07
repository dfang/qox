package aftersales

import (
	"github.com/jinzhu/gorm"
)

// Settlement AfterSale 售后厂家管理
type Settlement struct {
	gorm.Model

	ManufacturerName 	string

	UserID *uint
	// User   users.User

	// deposit 存入（完成服务）
	// withdraw 提现 ()
	// Direction string

	Fee float32
}
