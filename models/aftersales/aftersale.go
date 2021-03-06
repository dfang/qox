package aftersales

import (
	"fmt"

	"github.com/dfang/qor-demo/models/users"
	"github.com/jinzhu/gorm"
	"github.com/qor/audited"
	"github.com/qor/media/oss"
	"github.com/qor/transition"
)

// Aftersale 售后管理
type Aftersale struct {
	gorm.Model

	// 客户信息, 从京东后台导入或者扫描枪扫入的
	CustomerName    string
	CustomerPhone   string
	CustomerAddress string

	// -- ORDER_TYPE starts with Q 退货的取件单
	ServiceType    string
	ServiceContent string

	// 预约安装时间
	ReservedServiceTime string

	Source string

	// 品牌
	Brand string

	Fee float32

	// 从物流单来的
	// 商品名 order_items.item_name
	// 数量 order_items.quantity
	// 单价 单个商品售后价

	ItemName     string  `json:"item_name"`
	Quantity     uint    `json:"quantity"`
	PricePerUnit float32 `json:"price_per_unit"`
	FromOrderNo  string  `json:"from_order_no"`

	// 备注
	Remark string

	UserID uint
	User   users.User

	Images []AftersaleImage

	// AftersaleItems []AftersaleItem

	transition.Transition
	audited.AuditedModel
}

// BeforeCreate 初始状态
func (item *Aftersale) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("State", "created")
	return nil
}

// BeforeSave 验证费用
func (item *Aftersale) BeforeSave(scope *gorm.Scope) error {
	if item.Fee < 0 {
		return fmt.Errorf("费用不能小于或这等于0")
	}

	return nil
}

// Aftersale
// 姓名
// 地址
// 电话
// 服务类型 安装 清洗 维修
// 服务内容
// 备注

type AftersaleImage struct {
	gorm.Model
	Aftersale   Aftersale
	AftersaleID uint
	Image       oss.OSS `sql:"size:4294967295;" media_library:"url:/backend/{{class}}/{{primary_key}}/{{column}}.{{extension}};path:./private"`
}
