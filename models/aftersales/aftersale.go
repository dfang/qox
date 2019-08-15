package aftersales

import (
	"github.com/jinzhu/gorm"
)

// AfterSale 售后管理
type AfterSale struct {
	gorm.Model

	// 客户信息, 从京东后台导入或者扫描枪扫入的
	CustomerName    string
	CustomerPhone   string
	CustomerAddress string

	// -- ORDER_TYPE starts with Q 退货的取件单
	ServiceType    string
	ServiceContent string
	// 预约安装时间
	ReserverdServiceTime string
	Remark               string
}

// AfterSale
// 姓名
// 地址
// 电话
// 服务类型 安装 清洗 维修
// 服务内容
// 备注
