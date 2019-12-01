package orders

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/audited"
	"github.com/qor/transition"
)

/*
对配送时效是否满意
对服务态度是否满意
是否有开箱验货
师傅是否邀评
是否有留下联系方式, 方便后期有问题联系
师傅是否有介绍延保
是否有把商品放到指定位置
是否现场指导厂家400电话安装
有无问题要反馈
异常处理结果
*/

// OrderFollowUp 订单回访
type OrderFollowUp struct {
	gorm.Model
	audited.AuditedModel
	transition.Transition

	// OrderID uint
	OrderNo string `gorm:"unique;not null" json:"order_no"`

	// 对配送时效是否满意
	SatisfactionOfTimeliness string `json:"satisfaction_of_timeliness"`

	// 对服务态度是否满意
	SatisfactionOfServices string `json:"satisfaction_of_services"`

	// 是否有开箱验货
	InspectTheGoods string `json:"inspect_the_goods"`

	// 师傅是否邀评
	RequestFeedback string `json:"request_feedback"`

	// 是否有留下联系方式 方便后期有问题联系
	LeaveContactInfomation string `json:"leave_contact_information"`

	// 师傅是否有介绍延保
	IntroduceWarrantyExtension string `json:"introduce_warranty_extension"`

	// 是否有把商品放到指定位置
	PositionProperly string `json:"position_properly"`

	// 有无问题要反馈
	Feedback string `json:"feedback"`

	// 异常处理结果
	ExceptionHandling string `json:"exception_handling"`
}
