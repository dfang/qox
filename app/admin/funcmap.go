package admin

import (
	"html/template"
	"time"

	"github.com/dfang/qor-demo/models/aftersales"
	"github.com/dfang/qor-demo/models/orders"

	"github.com/jinzhu/now"
	"github.com/qor/admin"
)

func initFuncMap(Admin *admin.Admin) {
	// Admin.RegisterFuncMap("render_latest_order", renderLatestOrder)
	// Admin.RegisterFuncMap("render_latest_products", renderLatestProduct)
	Admin.RegisterFuncMap("render_latest_aftersales", renderLatestAftersales)
	Admin.RegisterFuncMap("render_latest_orders", renderLatestOrders)
	Admin.RegisterFuncMap("render_today_aftersales", renderTodayAftersales)
	Admin.RegisterFuncMap("render_today_orders", renderTodayOrders)
}

func renderLatestOrders(context *admin.Context) template.HTML {
	var orderContext = context.NewResourceContext("Order")
	orderContext.Searcher.Pagination.PerPage = 5
	// orderContext.SetDB(orderContext.GetDB().Where("state in (?)", []string{"paid"}))

	if orders, err := orderContext.FindMany(); err == nil {
		return orderContext.Render("index/table", orders)
	}
	return template.HTML("")
}

func renderLatestProduct(context *admin.Context) template.HTML {
	var productContext = context.NewResourceContext("Product")
	productContext.Searcher.Pagination.PerPage = 5
	// productContext.SetDB(productContext.GetDB().Where("state in (?)", []string{"paid"}))

	if products, err := productContext.FindMany(); err == nil {
		return productContext.Render("index/table", products)
	}
	return template.HTML("")
}

func renderLatestAftersales(context *admin.Context) template.HTML {
	var productContext = context.NewResourceContext("Aftersale")
	productContext.Searcher.Pagination.PerPage = 10
	// productContext.SetDB(productContext.GetDB().Where("state in (?)", []string{"paid"}))

	if products, err := productContext.FindMany(); err == nil {
		return productContext.Render("index/table", products)
	}
	return template.HTML("")
}

func renderTodayAftersales(context *admin.Context) template.HTML {
	var afterSaleContext = context.NewResourceContext("Aftersale")
	t := TodayAfterSalesCount{}

	// var count1 int
	// var count2 int

	afterSaleContext.GetDB().Model(&aftersales.Aftersale{}).Where("state = ?", "created").Count(&t.ToReserve)
	afterSaleContext.GetDB().Model(&aftersales.Aftersale{}).Where("state = ?", "inquired").Count(&t.ToSchedule)

	afterSaleContext.GetDB().Model(&aftersales.Aftersale{}).Where("state = ?", "scheduled").Count(&t.Scheduled)
	afterSaleContext.GetDB().Model(&aftersales.Aftersale{}).Where("state = ?", "processing").Count(&t.ToProcess)
	afterSaleContext.GetDB().Model(&aftersales.Aftersale{}).Where("state = ?", "processed").Count(&t.ToAudit)

	afterSaleContext.GetDB().Model(&aftersales.Aftersale{}).Where("state = ?", "audited").Count(&t.Audited)
	afterSaleContext.GetDB().Model(&aftersales.Aftersale{}).Where("state = ?", "frozen").Count(&t.Frozen)

	// 已指派的状态超过了20分钟就算超时了需要重新调度
	afterSaleContext.GetDB().Model(&aftersales.Aftersale{}).Where("state = ?", "overdue").Count(&t.Overdue)

	// t.Overdue = "0"
	t.FailedToAudit = "0"

	// fmt.Println(t.ToReserve)
	// fmt.Println(t.ToSchedule)

	return afterSaleContext.Render("today_aftersales", t)

	// return template.HTML("")
}

func renderTodayOrders(context *admin.Context) template.HTML {
	var ctx = context.NewResourceContext("Order")
	t := TodayOrdersCount{}
	// var count1 int
	// var count2 int

	ctx.GetDB().Model(&orders.Order{}).Where("created_at >= ?", now.BeginningOfDay()).Where("created_at <=? ", time.Now()).Where("order_no like ?", "Q%").Count(&t.ToPickUpTomorrow)
	ctx.GetDB().Model(&orders.Order{}).Where("created_at >= ?", now.BeginningOfDay().AddDate(0, 0, -1)).Where("created_at <=? ", now.EndOfDay().AddDate(0, 0, -1)).Where("order_no like ?", "Q%").Count(&t.ToPickUpToday)

	// 今日预约了
	ctx.GetDB().Model(&orders.Order{}).Where("created_at >= ?", now.BeginningOfDay()).Where("created_at <=? ", now.EndOfDay()).Count(&t.Reserved)

	// 今日需妥投
	ctx.GetDB().Model(&orders.Order{}).Where("reserverd_delivery_time = ?", now.BeginningOfDay().Format("2006-01-02")).Count(&t.ToDeliver)

	ctx.GetDB().Model(&orders.Order{}).Where("created_at >= ?", now.BeginningOfDay().AddDate(0, 0, -2)).Where("created_at <=? ", now.EndOfDay().AddDate(0, 0, -2)).Count(&t.YesterdayToDeliver)

	ctx.GetDB().Model(&orders.Order{}).Where("created_at >= ?", now.BeginningOfDay().AddDate(0, 0, -3)).Where("created_at <=? ", now.EndOfDay().AddDate(0, 0, -3)).Count(&t.TheDayBeforeYesterdayToDeliver)

	t.ToDeclare = "0"

	// fmt.Println(t.ToReserve)
	// fmt.Println(t.ToSchedule)

	return ctx.Render("today_orders", t)

	// return template.HTML("")
}

type TodayOrdersCount struct {
	// 待取件
	ToPickUpToday string

	ToPickUpTomorrow string

	// 待妥投
	ToDeliver string

	// 待报单
	ToDeclare string

	YesterdayToDeliver string

	TheDayBeforeYesterdayToDeliver string

	Reserved string
}

type TodayAfterSalesCount struct {
	// 待预约
	ToReserve string

	// 待指派
	ToSchedule string

	// 已超时
	Overdue string

	// 已指派
	Scheduled string

	// 审核不通过的
	FailedToAudit string

	// 已审核
	Audited string

	//已冻结
	Frozen string

	// 待上门
	ToProcess string

	// 已处理 待提交服务完成证明
	Processed string

	// 待审核
	ToAudit string
}
