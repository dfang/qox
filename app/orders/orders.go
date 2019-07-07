package orders

import (
	"fmt"
	"html/template"
	"time"

	// "net/http"

	"strings"

	"github.com/dfang/qor-demo/models/orders"
	"github.com/dfang/qor-demo/models/users"
	"github.com/dfang/qor-demo/utils/funcmapmaker"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/now"
	"github.com/qor/activity"
	"github.com/qor/admin"
	"github.com/qor/application"
	"github.com/qor/exchange"
	"github.com/qor/qor"
	"github.com/qor/render"
	"github.com/qor/transition"
)

// New new home app
func New(config *Config) *App {
	return &App{Config: config}
}

func NewWithDefault() *App {
	return &App{Config: &Config{}}
}

// App home app
type App struct {
	Config *Config
}

// Config home config struct
type Config struct {
}

var Genders = []string{"Men", "Women", "Kids"}

// ConfigureApplication configure application
func (app App) ConfigureApplication(application *application.Application) {
	controller := &Controller{View: render.New(&render.Config{AssetFileSystem: application.AssetFS.NameSpace("orders")}, "app/orders/views")}

	funcmapmaker.AddFuncMapMaker(controller.View)
	app.ConfigureAdmin(application.Admin)

	application.Router.Get("/cart", controller.Cart)
	application.Router.Put("/cart", controller.UpdateCart)
	application.Router.Post("/cart", controller.UpdateCart)
	application.Router.Get("/cart/checkout", controller.Checkout)
	application.Router.Put("/cart/checkout", controller.Checkout)
	application.Router.Post("/cart/complete", controller.Complete)
	application.Router.Post("/cart/complete/creditcard", controller.CompleteCreditCard)
	application.Router.Get("/cart/success", controller.CheckoutSuccess)
	// application.Router.Post("/order/callback/amazon", controller.AmazonCallback)
}

// ConfigureAdmin configure admin interface
func (App) ConfigureAdmin(Admin *admin.Admin) {
	Admin.AddMenu(&admin.Menu{Name: "Order Management", Priority: 1})
	// Add Order
	order := Admin.AddResource(&orders.Order{}, &admin.Config{Menu: []string{"Order Management"}})
	// order.Meta(&admin.Meta{Name: "ShippingAddress", Type: "single_edit"})
	// order.Meta(&admin.Meta{Name: "BillingAddress", Type: "single_edit"})
	order.Meta(&admin.Meta{Name: "ShippedAt", Type: "date"})
	order.Meta(&admin.Meta{Name: "PaymentLog", Type: "readonly", FormattedValuer: func(record interface{}, _ *qor.Context) (result interface{}) {
		order := record.(*orders.Order)
		return template.HTML(strings.Replace(strings.TrimSpace(order.PaymentLog), "\n", "<br>", -1))
	}})
	order.Meta(&admin.Meta{Name: "DeliveryMethod", Type: "select_one",
		Config: &admin.SelectOneConfig{
			Collection: func(_ interface{}, context *admin.Context) (options [][]string) {
				var methods []orders.DeliveryMethod
				context.GetDB().Find(&methods)
				for _, m := range methods {
					idStr := fmt.Sprintf("%d", m.ID)
					var option = []string{idStr, fmt.Sprintf("%s (%0.2f) руб", m.Name, m.Price)}
					options = append(options, option)
				}
				return options
			},
		},
	})

	order.Meta(&admin.Meta{Name: "created_at", Type: "string", FormattedValuer: func(record interface{}, _ *qor.Context) (result interface{}) {
		order := record.(*orders.Order)
		return order.CreatedAt.Local().Format("2006-01-02 15:04:05")
	}})

	order.Meta(&admin.Meta{Name: "updated_at", Type: "string", FormattedValuer: func(record interface{}, _ *qor.Context) (result interface{}) {
		order := record.(*orders.Order)
		return order.UpdatedAt.Local().Format("2006-01-02 15:04:05")
	}})

	order.Meta(&admin.Meta{Name: "customer_address", Type: "string", FormattedValuer: func(record interface{}, _ *qor.Context) (result interface{}) {
		order := record.(*orders.Order)
		return strings.Replace(order.CustomerAddress, "江西九江市修水县", "", -1)
	}})

	order.Meta(&admin.Meta{Name: "customer_phone", Type: "string", FormattedValuer: func(record interface{}, _ *qor.Context) (result interface{}) {
		order := record.(*orders.Order)
		phones := strings.Split(order.CustomerPhone, "/")
		if phones[0] == phones[1] {
			return phones[0]
		}
		return order.CustomerPhone
	}})
	order.Meta(&admin.Meta{Name: "man_to_deliver_id", Type: "string", FormattedValuer: func(record interface{}, ctx *qor.Context) (result interface{}) {
		order := record.(*orders.Order)
		if order.ManToDeliverID != "" {
			var user users.User
			ctx.DB.Where("id = ?", order.ManToDeliverID).Find(&user)
			return user.Name
		}
		return ""
	}})
	order.Meta(&admin.Meta{Name: "man_to_setup_id", Type: "string", FormattedValuer: func(record interface{}, ctx *qor.Context) (result interface{}) {
		order := record.(*orders.Order)
		if order.ManToSetupID != "" {
			var user users.User
			ctx.DB.Where("id = ?", order.ManToSetupID).Find(&user)
			return user.Name
		}
		return ""
	}})
	order.Meta(&admin.Meta{Name: "man_to_pickup_id", Type: "string", FormattedValuer: func(record interface{}, ctx *qor.Context) (result interface{}) {
		order := record.(*orders.Order)
		if order.ManToPickupID != "" {
			var user users.User
			ctx.DB.Where("id = ?", order.ManToPickupID).Find(&user)
			return user.Name
		}
		return ""
	}})

	orderItemMeta := order.Meta(&admin.Meta{Name: "OrderItems"})
	orderItemMeta.Resource.Meta(&admin.Meta{Name: "SizeVariation", Config: &admin.SelectOneConfig{Collection: sizeVariationCollection}})

	// define scopes for Order
	order.Scope(&admin.Scope{
		Name:    "Today",
		Label:   "Today",
		Default: true,
		Group:   "Filter By Date",
		Handler: func(db *gorm.DB, context *qor.Context) *gorm.DB {
			return db.Where("created_at >= ?", now.BeginningOfDay()).Where("created_at <=? ", time.Now())
		},
	})
	order.Scope(&admin.Scope{
		Name:  "Yesterday",
		Label: "Yesterday",
		Group: "Filter By Date",
		Handler: func(db *gorm.DB, context *qor.Context) *gorm.DB {
			now.WeekStartDay = time.Monday

			return db.Where("created_at >= ?", now.BeginningOfDay().AddDate(0, 0, -1)).Where("created_at <=? ", now.EndOfDay().AddDate(0, 0, -1))
		},
	})
	order.Scope(&admin.Scope{
		Name:  "ThisWeek",
		Label: "This Week",
		Group: "Filter By Date",
		Handler: func(db *gorm.DB, context *qor.Context) *gorm.DB {
			now.WeekStartDay = time.Monday
			return db.Where("created_at >= ?", now.BeginningOfWeek()).Where("created_at <=? ", now.EndOfWeek())
		},
	})
	order.Scope(&admin.Scope{
		Name:  "ThisMonth",
		Label: "This Month",
		Group: "Filter By Date",
		Handler: func(db *gorm.DB, context *qor.Context) *gorm.DB {
			now.WeekStartDay = time.Monday
			return db.Where("created_at >= ?", now.BeginningOfMonth()).Where("created_at <=? ", now.EndOfMonth())
		},
	})
	order.Scope(&admin.Scope{
		Name:  "ThisQuarter",
		Label: "This Quarter",
		Group: "Filter By Date",
		Handler: func(db *gorm.DB, context *qor.Context) *gorm.DB {
			return db.Where("created_at >= ?", now.BeginningOfQuarter()).Where("created_at <=? ", now.EndOfQuarter())
		},
	})
	order.Scope(&admin.Scope{
		Name:  "ThisYear",
		Label: "This Year",
		Group: "Filter By Date",
		Handler: func(db *gorm.DB, context *qor.Context) *gorm.DB {
			return db.Where("created_at >= ?", now.BeginningOfYear()).Where("created_at <=? ", now.EndOfYear())
		},
	})

	// filter by order state
	for _, state := range []string{"pending", "processing", "delivery_scheduled", "setup_scheduled", "pickup_scheduled", "cancelled", "shipped", "paid_cancelled", "returned"} {
		var state = state
		order.Scope(&admin.Scope{
			Name:  state,
			Label: strings.Title(strings.Replace(state, "_", " ", -1)),
			Group: "Order Status",
			Handler: func(db *gorm.DB, context *qor.Context) *gorm.DB {
				return db.Where(orders.Order{Transition: transition.Transition{State: state}})
			},
		})
	}

	// filter by order source
	order.Scope(&admin.Scope{
		Name:  "JD",
		Label: "JD",
		Group: "Filter By Source",
		Handler: func(db *gorm.DB, context *qor.Context) *gorm.DB {
			return db.Where(orders.Order{Source: "JD"})
		},
	})

	// filter by order type
	for _, state := range []string{"delivery", "setup", "delivery_and_setup", "repair", "clean", "sales"} {
		var state = state
		order.Scope(&admin.Scope{
			Name:  state,
			Label: strings.Title(strings.Replace(state, "_", " ", -1)),
			Group: "Filter By Order Type",
			Handler: func(db *gorm.DB, context *qor.Context) *gorm.DB {
				return db.Where(orders.Order{Transition: transition.Transition{State: state}})
			},
		})
	}

	// define actions for Order
	type trackingNumberArgument struct {
		TrackingNumber string
	}

	type deliveryActionArgument struct {
		ManToDeliver string
	}

	type setupActionArgument struct {
		ManToSetup string
	}

	type processingActionArgument struct {
		ShippingFee float32
		SetupFee    float32
		PickupFee   float32
		OrderType   string
	}
	processingActionResource := Admin.NewResource(&processingActionArgument{})
	processingActionResource.Meta(&admin.Meta{
		Name: "ShippingFee",
		Type: "float",
	})
	processingActionResource.Meta(&admin.Meta{
		Name: "SetupFee",
		Type: "float",
	})
	processingActionResource.Meta(&admin.Meta{
		Name: "PickupFee",
		Type: "float",
	})
	processingActionResource.Meta(&admin.Meta{
		Name:       "OrderType",
		Type:       "select_one",
		Collection: []string{"配送", "安装", "配送一体", "维修", "清洗"},
	})
	order.Action(&admin.Action{
		Name: "Processing",
		Handler: func(argument *admin.ActionArgument) error {
			db := argument.Context.GetDB()
			var (
				tx  = argument.Context.GetDB().Begin()
				arg = argument.Argument.(*processingActionArgument)
			)
			for _, record := range argument.FindSelectedRecords() {
				order := record.(*orders.Order)
				order.ShippingFee = arg.ShippingFee
				order.SetupFee = arg.SetupFee
				order.PickupFee = arg.PickupFee
				order.OrderType = arg.OrderType
				if err := orders.OrderState.Trigger("process", order, db); err != nil {
					return err
				}
				if err := tx.Save(order).Error; err != nil {
					tx.Rollback()
					return err
				}
				tx.Commit()
				return nil
			}
			return nil
		},
		Visible: func(record interface{}, context *admin.Context) bool {
			if order, ok := record.(*orders.Order); ok {
				return order.State == "pending"
			}
			return false
		},
		Resource: processingActionResource,
		Modes:    []string{"show", "menu_item"},
	})
	// order.Action(&admin.Action{
	// 	Name: "Processing",
	// 	Handler: func(argument *admin.ActionArgument) error {
	// 		for _, order := range argument.FindSelectedRecords() {
	// 			db := argument.Context.GetDB()
	// 			if err := orders.OrderState.Trigger("process", order.(*orders.Order), db); err != nil {
	// 				return err
	// 			}
	// 			db.Save(order)
	// 		}
	// 		return nil
	// 	},
	// 	Visible: func(record interface{}, context *admin.Context) bool {
	// 		if order, ok := record.(*orders.Order); ok {
	// 			return order.State == "pending"
	// 		}
	// 		return false
	// 	},
	// 	Modes: []string{"show", "menu_item"},
	// })

	// order.Action(&admin.Action{
	// 	Name: "Ship",
	// 	Handler: func(argument *admin.ActionArgument) error {
	// 		var (
	// 			tx                     = argument.Context.GetDB().Begin()
	// 			trackingNumberArgument = argument.Argument.(*trackingNumberArgument)
	// 		)
	// 		if trackingNumberArgument.TrackingNumber != "" {
	// 			for _, record := range argument.FindSelectedRecords() {
	// 				order := record.(*orders.Order)
	// 				order.TrackingNumber = &trackingNumberArgument.TrackingNumber
	// 				orders.OrderState.Trigger("ship", order, tx, "tracking number "+trackingNumberArgument.TrackingNumber)
	// 				if err := tx.Save(order).Error; err != nil {
	// 					tx.Rollback()
	// 					return err
	// 				}
	// 			}
	// 		} else {
	// 			return errors.New("invalid shipment number")
	// 		}

	// 		tx.Commit()
	// 		return nil
	// 	},
	// 	Visible: func(record interface{}, context *admin.Context) bool {
	// 		if order, ok := record.(*orders.Order); ok {
	// 			return order.State == "processing"
	// 		}
	// 		return false
	// 	},
	// 	Resource: Admin.NewResource(&trackingNumberArgument{}),
	// 	Modes:    []string{"show", "menu_item"},
	// })

	deliveryActionArgumentResource := Admin.NewResource(&deliveryActionArgument{})
	deliveryActionArgumentResource.Meta(&admin.Meta{
		Name: "ManToDeliver",
		Type: "select_one",
		Valuer: func(record interface{}, context *qor.Context) interface{} {
			// return record.(*users.User).ID
			return ""
		},
		Collection: func(value interface{}, context *qor.Context) (options [][]string) {
			var setupMen []users.User
			context.GetDB().Where("role = ?", "delivery_man").Find(&setupMen)
			for _, m := range setupMen {
				idStr := fmt.Sprintf("%d", m.ID)
				var option = []string{idStr, m.Name}
				options = append(options, option)
			}
			return options
		},
		// Collection: []string{"Male", "Female", "Unknown"},
	})
	// 安排配送
	order.Action(&admin.Action{
		Name: "Schedule Delivery",
		Handler: func(argument *admin.ActionArgument) error {
			var (
				tx = argument.Context.GetDB().Begin()
				// deliveryActionArgument = argument.Argument.(*deliveryActionArgument)
				arg = argument.Argument.(*deliveryActionArgument)
			)
			// if deliveryActionArgument.ManToDeliver != "" {
			for _, record := range argument.FindSelectedRecords() {
				order := record.(*orders.Order)
				order.ManToDeliverID = arg.ManToDeliver
				orders.OrderState.Trigger("schedule_delivery", order, tx, "man to deliver: "+arg.ManToDeliver)
				if err := tx.Save(order).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
			// } else {
			// 	return errors.New("invalid man to deliver")
			// }

			tx.Commit()
			return nil
		},
		Visible: func(record interface{}, context *admin.Context) bool {
			if order, ok := record.(*orders.Order); ok {
				return order.State == "processing"
			}
			return false
		},
		// Resource: Admin.NewResource(&deliveryActionArgument{}),
		Resource: deliveryActionArgumentResource,
		Modes:    []string{"show", "edit"},
	})

	setupActionArgumentResource := Admin.NewResource(&setupActionArgument{})
	setupActionArgumentResource.Meta(&admin.Meta{
		Name: "ManToSetup",
		Type: "select_one",
		Valuer: func(record interface{}, context *qor.Context) interface{} {
			// return record.(*users.User).ID
			return ""
		},
		Collection: func(value interface{}, context *qor.Context) (options [][]string) {
			var setupMen []users.User
			context.GetDB().Where("role = ?", "setup_man").Find(&setupMen)
			for _, m := range setupMen {
				idStr := fmt.Sprintf("%d", m.ID)
				var option = []string{idStr, m.Name}
				options = append(options, option)
			}
			return options
		},
		// Collection: []string{"Male", "Female", "Unknown"},
	})
	// 安排安装
	order.Action(&admin.Action{
		Name: "Schedule Setup",
		Handler: func(argument *admin.ActionArgument) error {
			var (
				tx  = argument.Context.GetDB().Begin()
				arg = argument.Argument.(*setupActionArgument)
			)
			// if setupArgument.ManToSetup != "" {
			for _, record := range argument.FindSelectedRecords() {
				order := record.(*orders.Order)
				order.ManToSetupID = arg.ManToSetup
				orders.OrderState.Trigger("schedule_setup", order, tx, "man to setup: "+arg.ManToSetup)
				if err := tx.Save(order).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
			// } else {
			// 	return errors.New("invalid man to setup")
			// }
			tx.Commit()
			return nil
		},
		Visible: func(record interface{}, context *admin.Context) bool {
			if order, ok := record.(*orders.Order); ok {
				return order.State == "processing"
			}
			return false
		},
		// Resource: Admin.NewResource(&setupActionArgument{}),
		Resource: setupActionArgumentResource,
		Modes:    []string{"show", "menu_item"},
	})

	order.Action(&admin.Action{
		Name: "Cancel",
		Handler: func(argument *admin.ActionArgument) error {
			for _, order := range argument.FindSelectedRecords() {
				db := argument.Context.GetDB()
				order := order.(*orders.Order)
				if err := orders.OrderState.Trigger("cancel", order, db); err != nil {
					return err
				}
				db.Save(order)
			}
			return nil
		},
		Visible: func(record interface{}, context *admin.Context) bool {
			if order, ok := record.(*orders.Order); ok {
				for _, state := range []string{"draft", "pending", "processing", "shipped"} {
					if order.State == state {
						return true
					}
				}
			}
			return false
		},
		Modes: []string{"show", "menu_item"},
	})

	order.Action(&admin.Action{
		Name:        "Export",
		URLOpenType: "slideout",
		URL: func(record interface{}, context *admin.Context) string {
			return "/admin/workers/new?job=Export Orders"
		},
		Modes: []string{"collection"},
	})

	// order.IndexAttrs("ID", "User", "PaymentAmount", "ShippedAt", "CancelledAt", "State", "ShippingAddress")
	// order.IndexAttrs("ID", "source", "order_no", "state", "order_type", "customer_name", "customer_address", "customer_phone", "receivables", "is_delivery_and_setup", "reserverd_delivery_time", "reserverd_setup_time", "man_to_deliver_id", "man_to_setup_id", "man_to_pickup_id", "shipping_fee", "setup_fee", "pickup_fee")
	order.IndexAttrs("ID", "source", "order_no", "state", "order_type", "customer_name", "customer_address", "customer_phone", "receivables",
		"is_delivery_and_setup", "reserverd_delivery_time", "reserverd_setup_time", "man_to_deliver_id", "man_to_setup_id", "man_to_pickup_id",
		"shipping_fee", "setup_fee", "pickup_fee", "created_at", "updated_at")
	order.NewAttrs("-DiscountValue", "-AbandonedReason", "-CancelledAt", "-PaymentLog", "-AmazonOrderReferenceID", "-AmazonAddressAccessToken")
	order.EditAttrs("-DiscountValue", "-AbandonedReason", "-CancelledAt", "-State", "-PaymentLog", "-AmazonOrderReferenceID", "-AmazonAddressAccessToken")
	order.ShowAttrs("-DiscountValue", "-State", "-AmazonAddressAccessToken")
	order.SearchAttrs("customer_name", "customer_phone", "order_no")

	// https://doc.getqor.com/admin/metas/select-one.html
	// Generate options by data from the database
	order.Meta(&admin.Meta{
		Name:  "ManToSetup",
		Type:  "select_one",
		Label: "Man To Setup",
		Config: &admin.SelectOneConfig{
			Collection: func(_ interface{}, context *admin.Context) (options [][]string) {
				var users []users.User
				context.GetDB().Where("role = ?", "setup_man").Find(&users)
				for _, n := range users {
					idStr := fmt.Sprintf("%d", n.ID)
					var option = []string{idStr, n.Name}
					options = append(options, option)
				}
				return options
			},
			AllowBlank: true,
		}})

	order.Meta(&admin.Meta{
		Name: "ManToDeliver",
		Type: "select_one",
		Config: &admin.SelectOneConfig{
			Collection: func(_ interface{}, context *admin.Context) (options [][]string) {
				var users []users.User
				context.GetDB().Where("role = ?", "delivery_man").Find(&users)

				for _, n := range users {
					idStr := fmt.Sprintf("%d", n.ID)
					var option = []string{idStr, n.Name}
					options = append(options, option)
				}
				return options
			},
			AllowBlank: true,
		}})

	oldSearchHandler := order.SearchHandler
	order.SearchHandler = func(keyword string, context *qor.Context) *gorm.DB {
		return oldSearchHandler(keyword, context).Where("state <> ? AND state <> ?", "", orders.DraftState)
	}

	// Add activity for order
	activity.Register(order)

	// Define another resource for same model
	abandonedOrder := Admin.AddResource(&orders.Order{}, &admin.Config{Name: "Abandoned Order", Menu: []string{"Order Management"}})
	abandonedOrder.Meta(&admin.Meta{Name: "ShippingAddress", Type: "single_edit"})
	abandonedOrder.Meta(&admin.Meta{Name: "BillingAddress", Type: "single_edit"})

	// Define default scope for abandoned orders
	abandonedOrder.Scope(&admin.Scope{
		Default: true,
		Handler: func(db *gorm.DB, context *qor.Context) *gorm.DB {
			return db.Where("abandoned_reason IS NOT NULL AND abandoned_reason <> ?", "")
		},
	})

	// Define scopes for abandoned orders
	for _, amount := range []int{5000, 10000, 20000} {
		var amount = amount
		abandonedOrder.Scope(&admin.Scope{
			Name:  fmt.Sprint(amount),
			Group: "Amount Greater Than",
			Handler: func(db *gorm.DB, context *qor.Context) *gorm.DB {
				return db.Where("payment_amount > ?", amount)
			},
		})
	}

	abandonedOrder.IndexAttrs("-ShippingAddress", "-BillingAddress", "-DiscountValue", "-OrderItems")
	abandonedOrder.NewAttrs("-DiscountValue")
	abandonedOrder.EditAttrs("-DiscountValue")
	abandonedOrder.ShowAttrs("-DiscountValue")

	// Delivery Methods
	Admin.AddResource(&orders.DeliveryMethod{}, &admin.Config{Menu: []string{"Site Management"}})

	// installs := Admin.AddResource(&orders.Order{Source: "JD"}, &admin.Config{Name: "Installs", Menu: []string{"Order Management"}})
	// installs.Scope(&admin.Scope{
	// 	Default: true,
	// 	Handler: func(db *gorm.DB, context *qor.Context) *gorm.DB {
	// 		return db.Where("source IS NOT NULL")
	// 	},
	// })
	// // installs.IndexAttrs("ID", "source", "order_no", "customer_name", "customer_address", "customer_phone", "is_delivery_and_setup", "reserverd_delivery_time", "reserverd_setup_time", "man_to_deliver_id", "man_to_setup_id", "man_to_pickup_id", "state")

	// installs.IndexAttrs("ID", "source", "order_no", "customer_name", "customer_address", "customer_phone", "is_delivery_and_setup", "reserverd_delivery_time", "reserverd_setup_time", "man_to_deliver_id", "man_to_setup_id", "man_to_pickup_id", "state")

	// // define scopes for Order
	// for _, state := range []string{"pending", "processing", "delivery_scheduled", "setup_scheduled", "pickup_scheduled", "cancelled", "shipped", "paid_cancelled", "returned"} {
	// 	var state = state
	// 	installs.Scope(&admin.Scope{
	// 		Name:  state,
	// 		Label: strings.Title(strings.Replace(state, "_", " ", -1)),
	// 		Group: "Order Status",
	// 		Handler: func(db *gorm.DB, context *qor.Context) *gorm.DB {
	// 			return db.Where(orders.Order{Transition: transition.Transition{State: strings.Title(state)}})
	// 		},
	// 	})
	// }

	// Define Resource
	o1 := exchange.NewResource(&orders.Order{}, exchange.Config{PrimaryField: "customer_name"})
	// Define columns are exportable/importable
	o1.Meta(&exchange.Meta{Name: "customer_name"})
	o1.Meta(&exchange.Meta{Name: "customer_address"})
	o1.Meta(&exchange.Meta{Name: "customer_phone"})
	o1.Meta(&exchange.Meta{Name: "receivables"})
	o1.Meta(&exchange.Meta{Name: "reserverd_delivery_time"})
	o1.Meta(&exchange.Meta{Name: "reserverd_setup_time"})

}

func sizeVariationCollection(resource interface{}, context *qor.Context) (results [][]string) {
	// for _, sizeVariation := range products.SizeVariations() {
	// 	results = append(results, []string{strconv.Itoa(int(sizeVariation.ID)), sizeVariation.Stringify()})
	// }
	return
}
