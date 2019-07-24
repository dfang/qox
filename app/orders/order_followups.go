package orders

import (
	"github.com/dfang/qor-demo/models/orders"
	"github.com/qor/admin"
)

// ConfigureAdmin configure admin interface
func (App) ConfigureOrderFollowUpsAdmin(Admin *admin.Admin) {
	// Admin.AddMenu(&admin.Menu{Name: "Order Management", Priority: 1})

	// Add Order
	order_followup := Admin.AddResource(&orders.OrderFollowUp{}, &admin.Config{Menu: []string{"Order Management"}})
	configureVisibleFieldsForOrderFollowUps(order_followup)
}

func configureVisibleFieldsForOrderFollowUps(item *admin.Resource) {
	item.ShowAttrs("-CreatedBy", "-UpdatedBy")
	item.NewAttrs("-CreatedBy", "-UpdatedBy")
}
