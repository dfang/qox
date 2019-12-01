package orders

import (
	"github.com/dfang/qor-demo/models/orders"
	"github.com/qor/admin"
)

// ConfigureAdmin configure admin interface
func (App) ConfigureOrderFollowUpsAdmin(Admin *admin.Admin) {
	// Admin.AddMenu(&admin.Menu{Name: "Order Management", Priority: 1})

	// Add Order
	followup := Admin.AddResource(&orders.OrderFollowUp{}, &admin.Config{Menu: []string{"Order Management"}})
	configureVisibleFieldsForOrderFollowUps(followup)

	followup.Action(&admin.Action{
		Name:        "导出",
		URLOpenType: "slideout",
		URL: func(record interface{}, context *admin.Context) string {
			return "/admin/workers/new?job=Export_FollowUps"
		},
		Modes: []string{"collection"},
	})

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
}

func configureVisibleFieldsForOrderFollowUps(item *admin.Resource) {
	item.IndexAttrs("-CreatedBy", "-UpdatedBy", "-State")
	item.ShowAttrs("-CreatedBy", "-UpdatedBy")
	item.NewAttrs("-CreatedBy", "-UpdatedBy")
	item.EditAttrs("-CreatedBy", "-UpdatedBy")
}
