package orders

type PaymentMethod = string

const (
	COD        PaymentMethod = "COD"
	AmazonPay  PaymentMethod = "AmazonPay"
	CreditCard PaymentMethod = "CreditCard"
)

type Source string

const (
	JD     Source = "京东物流"
	Suning Source = "苏宁帮客"
	Midea  Source = "美的安得"
)

var (
	SOURCES = []string{
		"京东物流", "苏宁帮客", "美的安得",
	}
)

type OrderType string

const (
	Delivery         OrderType = "配送"
	Logistics        OrderType = "物流"
	Setup            OrderType = "安装"
	Repair           OrderType = "维修"
	Clean            OrderType = "清洗"
	Sale             OrderType = "销售"
	DeliveryAndSetup OrderType = "送装一体"
)

var (
	ORDER_TYPES = []string{"配送", "安装", "送装一体", "维修", "清洗", "销售", "物流"}
)
