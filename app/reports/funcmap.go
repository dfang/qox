package reports

import (
	"fmt"
	"html/template"

	"github.com/qor/admin"
)

func initFuncMap(Admin *admin.Admin) {
	// Admin.RegisterFuncMap("render_latest_order", renderLatestOrder)
	// Admin.RegisterFuncMap("render_latest_products", renderLatestProduct)
	Admin.RegisterFuncMap("render_by_sources", renderBySources)
	Admin.RegisterFuncMap("render_by_brands", renderByBrands)
	Admin.RegisterFuncMap("render_by_names", renderByNames)

	// Admin.RegisterFuncMap("render_latest_aftersales", renderLatestAftersales)
	// Admin.RegisterFuncMap("render_today", renderToday)

	Admin.RegisterFuncMap("toAftersaleReportBySource", toAftersaleReportBySource)
	Admin.RegisterFuncMap("toAftersaleReportByBrand", toAftersaleReportByBrand)
	Admin.RegisterFuncMap("toAftersaleReportByName", toAftersaleReportByName)

}

func renderBySources(context *admin.Context) template.HTML {
	// var productContext = context.NewResourceContext("Product")
	// productContext.Searcher.Pagination.PerPage = 5
	// productContext.SetDB(productContext.GetDB().Where("state in (?)", []string{"paid"}))
	db := context.GetDB()

	t := context.Request.URL.Query().Get("type")
	var sqlStr string
	switch t {
	case "year":
		sqlStr = "select to_char(date_trunc('day', aftersales.updated_at), 'YYYY') as time, source, count(*) as count, sum(fee) as sum from aftersales where state = 'completed' group by 1, source"
	case "month":
		sqlStr = "select to_char(date_trunc('day', aftersales.updated_at), 'YYYY-MM') as time, source, count(*) as count, sum(fee) as sum from aftersales where state = 'completed' group by 1, source order by time desc"
	case "day":
		sqlStr = "select to_char(date_trunc('day', aftersales.updated_at), 'YYYY-MM-DD') as time, source, count(*) as count, sum(fee) as sum from aftersales where state = 'completed' group by 1, source order by time desc"
	default:
		sqlStr = "select to_char(date_trunc('day', aftersales.updated_at), 'YYYY-MM') as time, source, count(*) as count, sum(fee) as sum from aftersales where state = 'completed' group by 1, source order by time desc"
	}

	var bySourceContext = context.NewResourceContext("AftersaleReportBySource")
	var result []AftersaleReportBySource
	// db.Raw("select to_char(date_trunc('day', aftersales.updated_at), 'YYYY-MM-DD') as time, source, count(*) as count, sum(fee) as sum from aftersales where state = 'completed' group by 1, source").Scan(&result)
	db.Raw(sqlStr).Scan(&result)

	fmt.Println("result is .....")
	fmt.Println(result)

	if len(result) > 0 {
		return bySourceContext.Render("index/table2", result)
	}
	return template.HTML("")

	// if products, err := Context.FindMany(); err == nil {
	// 	return Context.Render("index/table", products)
	// }
	// return template.HTML("")
}

func renderByBrands(context *admin.Context) template.HTML {
	// var productContext = context.NewResourceContext("Product")
	// productContext.Searcher.Pagination.PerPage = 5
	// productContext.SetDB(productContext.GetDB().Where("state in (?)", []string{"paid"}))
	db := context.GetDB()

	t := context.Request.URL.Query().Get("type")
	var sqlStr string
	switch t {
	case "year":
		sqlStr = "select to_char(date_trunc('day', aftersales.updated_at), 'YYYY') as time, brand, count(*) as count, sum(fee) as sum from aftersales where state = 'completed' group by 1, brand order by time desc"
	case "month":
		sqlStr = "select to_char(date_trunc('day', aftersales.updated_at), 'YYYY-MM') as time, brand, count(*) as count, sum(fee) as sum from aftersales where state = 'completed' group by 1, brand order by time desc"
	case "day":
		sqlStr = "select to_char(date_trunc('day', aftersales.updated_at), 'YYYY-MM-DD') as time, brand, count(*) as count, sum(fee) as sum from aftersales where state = 'completed' group by 1, brand order by time desc"
	default:
		sqlStr = "select to_char(date_trunc('day', aftersales.updated_at), 'YYYY-MM') as time, brand, count(*) as count, sum(fee) as sum from aftersales where state = 'completed' group by 1, brand order by time desc"
	}

	var ctx = context.NewResourceContext("AftersaleReportByBrand")
	var result []AftersaleReportByBrand
	// db.Raw("select to_char(date_trunc('day', aftersales.updated_at), 'YYYY-MM-DD') as time, source, count(*) as count, sum(fee) as sum from aftersales where state = 'completed' group by 1, source").Scan(&result)
	db.Raw(sqlStr).Scan(&result)

	fmt.Println("result is .....")
	fmt.Println(result)

	if len(result) > 0 {
		return ctx.Render("index/table3", result)
	}
	return template.HTML("")

	// if products, err := Context.FindMany(); err == nil {
	// 	return Context.Render("index/table", products)
	// }
	// return template.HTML("")
}

func renderByNames(context *admin.Context) template.HTML {
	// var productContext = context.NewResourceContext("Product")
	// productContext.Searcher.Pagination.PerPage = 5
	// productContext.SetDB(productContext.GetDB().Where("state in (?)", []string{"paid"}))
	db := context.GetDB()

	t := context.Request.URL.Query().Get("type")
	var sqlStr string
	switch t {
	case "year":
		sqlStr = "select to_char(date_trunc('day', aftersales.updated_at), 'YYYY') as time, users.name, count(*), sum(fee) from aftersales inner join users on aftersales.user_id = users.id where state = 'completed' group by users.name, 1 order by time desc"
	case "month":
		sqlStr = "select to_char(date_trunc('day', aftersales.updated_at), 'YYYY-MM') as time, users.name, count(*), sum(fee) from aftersales inner join users on aftersales.user_id = users.id where state = 'completed' group by users.name, 1 order by time desc"
	case "day":
		sqlStr = "select to_char(date_trunc('day', aftersales.updated_at), 'YYYY-MM-DD') as time, users.name, count(*), sum(fee) from aftersales inner join users on aftersales.user_id = users.id where state = 'completed' group by users.name, 1 order by time desc"
	default:
		sqlStr = "select to_char(date_trunc('day', aftersales.updated_at), 'YYYY-MM') as time, users.name, count(*), sum(fee) from aftersales inner join users on aftersales.user_id = users.id where state = 'completed' group by users.name, 1 order by time desc"
	}

	var ctx = context.NewResourceContext("AftersaleReportByName")
	var result []AftersaleReportByName
	// db.Raw("select to_char(date_trunc('day', aftersales.updated_at), 'YYYY-MM-DD') as time, source, count(*) as count, sum(fee) as sum from aftersales where state = 'completed' group by 1, source").Scan(&result)
	db.Raw(sqlStr).Scan(&result)

	fmt.Println("result is .....")
	fmt.Println(result)

	if len(result) > 0 {
		return ctx.Render("index/table4", result)
	}
	return template.HTML("")

	// if products, err := Context.FindMany(); err == nil {
	// 	return Context.Render("index/table", products)
	// }
	// return template.HTML("")
}

func toAftersaleReportBySource(f interface{}) []AftersaleReportBySource {
	result, ok := f.([]AftersaleReportBySource)
	if ok {
		return result
	}
	return []AftersaleReportBySource{}
}

func toAftersaleReportByBrand(f interface{}) []AftersaleReportByBrand {
	result, ok := f.([]AftersaleReportByBrand)
	if ok {
		return result
	}
	return []AftersaleReportByBrand{}
}

func toAftersaleReportByName(f interface{}) []AftersaleReportByName {
	result, ok := f.([]AftersaleReportByName)
	if ok {
		return result
	}
	return []AftersaleReportByName{}
}
