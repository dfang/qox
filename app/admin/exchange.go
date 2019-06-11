package admin

import (
	"github.com/dfang/qor-demo/models/orders"
	"github.com/dfang/qor-demo/models/products"
	"github.com/qor/exchange"
	"github.com/qor/qor"
	"github.com/qor/qor/resource"
	"github.com/qor/qor/utils"
	"github.com/qor/validations"
)

// ProductExchange product exchange
var ProductExchange = exchange.NewResource(&products.Product{}, exchange.Config{PrimaryField: "Code"})

// OrderExchange order exchange
var OrderExchange = exchange.NewResource(&orders.Order{}, exchange.Config{PrimaryField: "order_no"})

func init() {
	OrderExchange.Meta(&exchange.Meta{Name: "source"})
	OrderExchange.Meta(&exchange.Meta{Name: "order_no"})
	OrderExchange.Meta(&exchange.Meta{Name: "state"})
	OrderExchange.Meta(&exchange.Meta{Name: "order_type"})
	OrderExchange.Meta(&exchange.Meta{Name: "customer_name"})
	OrderExchange.Meta(&exchange.Meta{Name: "customer_address"})
	OrderExchange.Meta(&exchange.Meta{Name: "customer_phone"})
	OrderExchange.Meta(&exchange.Meta{Name: "receivables"})
	OrderExchange.Meta(&exchange.Meta{Name: "is_delivery_and_setup"})
	OrderExchange.Meta(&exchange.Meta{Name: "reserverd_delivery_time"})
	OrderExchange.Meta(&exchange.Meta{Name: "reserverd_setup_time"})
	OrderExchange.Meta(&exchange.Meta{Name: "man_to_deliver_id"})
	OrderExchange.Meta(&exchange.Meta{Name: "man_to_setup_id"})
	OrderExchange.Meta(&exchange.Meta{Name: "man_to_pickup_id"})
	OrderExchange.Meta(&exchange.Meta{Name: "shipping_fee"})
	OrderExchange.Meta(&exchange.Meta{Name: "setup_fee"})
	OrderExchange.Meta(&exchange.Meta{Name: "pickup_fee"})
	OrderExchange.Meta(&exchange.Meta{Name: "created_at"})
	OrderExchange.Meta(&exchange.Meta{Name: "updated_at"})

	ProductExchange.Meta(&exchange.Meta{Name: "Code"})
	ProductExchange.Meta(&exchange.Meta{Name: "Name"})
	ProductExchange.Meta(&exchange.Meta{Name: "Price"})

	ProductExchange.AddValidator(&resource.Validator{
		Handler: func(record interface{}, metaValues *resource.MetaValues, context *qor.Context) error {
			if utils.ToInt(metaValues.Get("Price").Value) < 100 {
				return validations.NewError(record, "Price", "price can't less than 100")
			}
			return nil
		},
	})
}
