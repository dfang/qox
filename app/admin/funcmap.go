package admin

import (
	"html/template"

	"github.com/dfang/qor-demo/models/aftersales"

	"github.com/qor/admin"
)

func initFuncMap(Admin *admin.Admin) {
	// Admin.RegisterFuncMap("render_latest_order", renderLatestOrder)
	// Admin.RegisterFuncMap("render_latest_products", renderLatestProduct)
	Admin.RegisterFuncMap("render_latest_aftersales", renderLatestAftersales)
	Admin.RegisterFuncMap("render_today", renderToday)
}

func renderLatestOrder(context *admin.Context) template.HTML {
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

func renderToday(context *admin.Context) template.HTML {
	var afterSaleContext = context.NewResourceContext("Aftersale")
	t := Today{}

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

	return afterSaleContext.Render("today", t)

	// return template.HTML("")
}

type Today struct {
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
