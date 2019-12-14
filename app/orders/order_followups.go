package orders

import (
	"time"

	"github.com/dfang/qor-demo/models/orders"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/now"
	"github.com/qor/admin"
	"github.com/qor/qor"
)

// ConfigureAdmin configure admin interface
func (App) ConfigureOrderFollowUpsAdmin(Admin *admin.Admin) {
	// Admin.AddMenu(&admin.Menu{Name: "Order Management", Priority: 1})

	// Add Order
	followup := Admin.AddResource(&orders.OrderFollowUp{}, &admin.Config{Menu: []string{"Order Management"}})
	configureVisibleFieldsForOrderFollowUps(followup)

	// followup.Action(&admin.Action{
	// 	Name:        "导出",
	// 	URLOpenType: "slideout",
	// 	URL: func(record interface{}, context *admin.Context) string {
	// 		return "/admin/workers/new?job=Export_FollowUps"
	// 	},
	// 	Modes: []string{"collection"},
	// })

	followup.Meta(&admin.Meta{
		Name:       "SatisfactionOfTimeliness",
		Type:       "select_one",
		Collection: []string{"是", "否"},
	})

	followup.Meta(&admin.Meta{
		Name:       "SatisfactionOfServices",
		Type:       "select_one",
		Collection: []string{"是", "否"},
	})

	followup.Meta(&admin.Meta{
		Name:       "InspectTheGoods",
		Type:       "select_one",
		Collection: []string{"是", "否"},
	})

	followup.Meta(&admin.Meta{
		Name:       "RequestFeedback",
		Type:       "select_one",
		Collection: []string{"是", "否"},
	})

	followup.Meta(&admin.Meta{
		Name:       "LeaveContactInfomation",
		Type:       "select_one",
		Collection: []string{"是", "否"},
	})

	followup.Meta(&admin.Meta{
		Name:       "IntroduceWarrantyExtension",
		Type:       "select_one",
		Collection: []string{"是", "否"},
	})

	followup.Meta(&admin.Meta{
		Name:       "PositionProperly",
		Type:       "select_one",
		Collection: []string{"是", "否"},
	})

	followup.Meta(&admin.Meta{Name: "OrderNo", Type: "followup_orderno_field", FormattedValuer: func(record interface{}, _ *qor.Context) (result interface{}) {
		m := record.(*orders.OrderFollowUp)
		return m
	}})

	followup.Scope(&admin.Scope{
		Name:  "noanswer",
		Label: "无人接听",
		Handler: func(db *gorm.DB, context *qor.Context) *gorm.DB {
			now.WeekStartDay = time.Monday
			// select order_no, customer_name, item_name::varchar(20), quantity, created_at
			// from orders_view
			// where created_at between now() - interval '2 day' and  now() - interval '1 day';
			// return db.Where("created_at between now() - interval '2 day' and  now() - interval '1 day'")
			return db.Where("Feedback like ?", "%无人接听")
		},
	})
}

func configureVisibleFieldsForOrderFollowUps(item *admin.Resource) {
	// item.IndexAttrs("-CreatedBy", "-UpdatedBy", "-State")
	item.IndexAttrs("OrderNo", "SatisfactionOfTimeliness", "SatisfactionOfServices", "InspectTheGoods", "RequestFeedback", "LeaveContactInfomation", "IntroduceWarrantyExtension", "PositionProperly", "CreatedAt")
	item.ShowAttrs("-CreatedBy", "-UpdatedBy", "-State")
	item.NewAttrs("-CreatedBy", "-UpdatedBy", "-State")
	item.EditAttrs("-CreatedBy", "-UpdatedBy", "-State")
}
