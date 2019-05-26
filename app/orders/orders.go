package orders

import (
	"errors"
	"fmt"
	"html/template"

	// "net/http"

	"strings"

	"github.com/dfang/qor-demo/config/application"
	"github.com/dfang/qor-demo/models/orders"
	"github.com/dfang/qor-demo/models/users"
	"github.com/dfang/qor-demo/utils/funcmapmaker"
	"github.com/jinzhu/gorm"
	"github.com/qor/activity"
	"github.com/qor/admin"
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
	Admin.AddMenu(&admin.Menu{Name: "Order Management", Priority: 2})
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

	orderItemMeta := order.Meta(&admin.Meta{Name: "OrderItems"})
	orderItemMeta.Resource.Meta(&admin.Meta{Name: "SizeVariation", Config: &admin.SelectOneConfig{Collection: sizeVariationCollection}})

	// define scopes for Order
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

	order.Scope(&admin.Scope{
		Name:  "Today",
		Label: "Today",
		Group: "Filter By Date",
		Handler: func(db *gorm.DB, context *qor.Context) *gorm.DB {
			return db.Where(orders.Order{})
		},
	})
	order.Scope(&admin.Scope{
		Name:  "ThisWeek",
		Label: "This Week",
		Group: "Filter By Date",
		Handler: func(db *gorm.DB, context *qor.Context) *gorm.DB {
			return db.Where(orders.Order{})
		},
	})

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

	order.Action(&admin.Action{
		Name: "Processing",
		Handler: func(argument *admin.ActionArgument) error {
			for _, order := range argument.FindSelectedRecords() {
				db := argument.Context.GetDB()
				if err := orders.OrderState.Trigger("process", order.(*orders.Order), db); err != nil {
					return err
				}
				db.Save(order)
			}
			return nil
		},
		Visible: func(record interface{}, context *admin.Context) bool {
			if order, ok := record.(*orders.Order); ok {
				return order.State == "pending"
			}
			return false
		},
		Modes: []string{"show", "menu_item"},
	})

	order.Action(&admin.Action{
		Name: "Ship",
		Handler: func(argument *admin.ActionArgument) error {
			var (
				tx                     = argument.Context.GetDB().Begin()
				trackingNumberArgument = argument.Argument.(*trackingNumberArgument)
			)

			if trackingNumberArgument.TrackingNumber != "" {
				for _, record := range argument.FindSelectedRecords() {
					order := record.(*orders.Order)
					order.TrackingNumber = &trackingNumberArgument.TrackingNumber
					orders.OrderState.Trigger("ship", order, tx, "tracking number "+trackingNumberArgument.TrackingNumber)
					if err := tx.Save(order).Error; err != nil {
						tx.Rollback()
						return err
					}
				}
			} else {
				return errors.New("invalid shipment number")
			}

			tx.Commit()
			return nil
		},
		Visible: func(record interface{}, context *admin.Context) bool {
			if order, ok := record.(*orders.Order); ok {
				return order.State == "processing"
			}
			return false
		},
		Resource: Admin.NewResource(&trackingNumberArgument{}),
		Modes:    []string{"show", "menu_item"},
	})

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
		Name: "配送",
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
				// orders.OrderState.Trigger("deliver", order, tx, "man to deliver "+orderActionArgument.ManToDeliver)
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
		Name: "安装",
		Handler: func(argument *admin.ActionArgument) error {
			var (
				tx  = argument.Context.GetDB().Begin()
				arg = argument.Argument.(*setupActionArgument)
				// setupArgument = argument.Argument.(*setupActionArgument)
			)
			// if setupArgument.ManToSetup != "" {
			for _, record := range argument.FindSelectedRecords() {
				order := record.(*orders.Order)
				order.ManToSetupID = arg.ManToSetup
				// 		// orders.OrderState.Trigger("ship", order, tx, "tracking number "+setupActionArgument.ManToSetup)
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

	order.IndexAttrs("ID", "User", "PaymentAmount", "ShippedAt", "CancelledAt", "State", "ShippingAddress")
	order.NewAttrs("-DiscountValue", "-AbandonedReason", "-CancelledAt", "-PaymentLog", "-AmazonOrderReferenceID", "-AmazonAddressAccessToken")
	order.EditAttrs("-DiscountValue", "-AbandonedReason", "-CancelledAt", "-State", "-PaymentLog", "-AmazonOrderReferenceID", "-AmazonAddressAccessToken")
	order.ShowAttrs("-DiscountValue", "-State", "-AmazonAddressAccessToken")
	order.SearchAttrs("ID", "User.Name", "User.Email", "AmazonOrderReferenceID", "ShippingAddress.Phone", "ShippingAddress.ContactName", "ShippingAddress.Address1", "ShippingAddress.Address2")

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
}

func sizeVariationCollection(resource interface{}, context *qor.Context) (results [][]string) {
	// for _, sizeVariation := range products.SizeVariations() {
	// 	results = append(results, []string{strconv.Itoa(int(sizeVariation.ID)), sizeVariation.Stringify()})
	// }
	return
}

// func deliveryManCollection(resource interface{}, context *qor.Context) (results []users.User) {
// 	// for _, sizeVariation := range products.SizeVariations() {
// 	// 	results = append(results, []string{strconv.Itoa(int(sizeVariation.ID)), sizeVariation.Stringify()})
// 	// }
// 	// return context.GetDB().Find(&users.User)
// }

// getAvailableLocales Get Available DevelieryMan
// func getAvailableLocales(req *http.Request, currentUser interface{}) []string {
// 	if user, ok := currentUser.(viewableLocalesInterface); ok {
// 		return user.ViewableLocales()
// 	}

// 	if user, ok := currentUser.(availableLocalesInterface); ok {
// 		return user.AvailableLocales()
// 	}
// 	return []string{Global}
// }
// 	Collection: func(_ interface{}, context *admin.Context) (options [][]string) {
// 				var methods []orders.DeliveryMethod
// 				context.GetDB().Find(&methods)

// 				for _, m := range methods {
// 					idStr := fmt.Sprintf("%d", m.ID)
// 					var option = []string{idStr, fmt.Sprintf("%s (%0.2f) руб", m.Name, m.Price)}
// 					options = append(options, option)
// 				}

// 				return options
// 			},
