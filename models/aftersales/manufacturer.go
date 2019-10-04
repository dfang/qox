package aftersales

import (
	"github.com/jinzhu/gorm"
)

// Manufacturer AfterSale 售后厂家管理
type Manufacturer struct {
	gorm.Model

	// 客户信息, 从京东后台导入或者扫描枪扫入的
	ManufacturerName string
	URL              string
}
