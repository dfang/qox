package aftersales

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/transition"
)
var (
	// OrderState order's state machine
	OrderState = transition.New(&AfterSale{})
)


var (
	// DraftState draft state
	DraftState = "created"
)

// State
// 括号里是Action

// (created)----> created 已接收状态 -----(inquire)---> inquired 未预约状态
// -----（schedule）----> scheduled 已预约  ---(confirm_schedule)--> schedule_confirmed 待处理
// -----(confirm_complete)----> to_be_audited 待审核 ---audit--> audited -----> finalized 锁定


// Dashboard for Operators
// 待预约   1 （aftersales.state == "created"）
// 待指派   1  (aftersales.state == "inquired"）
// 超时任务单 1  (aftersales.预约时间== "空" 需要重新指派）
// 待审核   1  (aftersales.state == "to_be_audited"）


var STATES = [] string{

	"created",      // 建单之后

	"inquired",     // 信息员给用户打电话预约之后

	"scheduled",    // 指派师傅之后

	"overduled",    // 指派师傅之后, 师傅未给用户打电话预约超时了

	"processing",   // 师傅给用户打过电话，确认了上门时间的状态

	"processed",    // 师傅上传了照片，等待审核

	"audited",      // 审核通过

	"audit_failed",      // 审核不通过
}

func init() {
	// Define Order's States
	OrderState.Initial("created")

	// 和用户预约大概时间
	OrderState.Event("inquire").To("inquired").From("created").After(func(value interface{}, tx *gorm.DB) (err error) {
		// order := value.(*AfterSale)
		// tx.Model(order).Association("OrderItems").Find(&order.OrderItems)
		// for _, item := range order.OrderItems {
		// }
		return nil
	})

	// 指派师傅
	OrderState.Event("schedule").To("scheduled").From("inquired").After(func(value interface{}, tx *gorm.DB) (err error) {
		return nil
	})

	// 根据师傅上传的照片 审核服务是否完成
	OrderState.Event("audit").To("audited").From("processed").After(func(value interface{}, tx *gorm.DB) (err error) {
		return nil
	})

}