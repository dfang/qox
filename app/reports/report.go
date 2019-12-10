package reports

import (
	"github.com/qor/admin"
	"github.com/qor/application"
)

// AftersaleReportByName 服务单完成报告
type AftersaleReportByName struct {
	Time  string `json:"time"`
	Name  string `json:"name"`
	Count string `json:"count"`
	Sum   string `json:"sum"`
}

// AftersaleReportByBrand 服务单完成报告
type AftersaleReportByBrand struct {
	Time  string `json:"time"`
	Brand string `json:"brand"`
	Count string `json:"count"`
	Sum   string `json:"sum"`
}

// AftersaleReportBySource 服务单完成报告
type AftersaleReportBySource struct {
	Time   string `json:"time"`
	Source string `json:"source"`
	Count  string `json:"count"`
	Sum    string `json:"sum"`
}

// OrdersCount 服务单完成报告
type OrdersCount struct {
	Time  string `json:"time"`
	Count string `json:"count"`
}

// New new home app
func New(config *Config) *App {
	if config.Prefix == "" {
		config.Prefix = "/admin"
	}
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
	Prefix string
}

func (app App) ConfigureApplication(application *application.Application) {

	Admin := application.Admin

	// Admin.AddMenu(&admin.Menu{Name: "Today", Link: "/today", Priority: 1})
	Admin.AddMenu(&admin.Menu{Name: "统计品牌", Link: "/admin/reports/by_brands", Ancestors: []string{"Reports Management"}})
	Admin.AddMenu(&admin.Menu{Name: "统计来源", Link: "/admin/reports/by_sources", Ancestors: []string{"Reports Management"}})
	Admin.AddMenu(&admin.Menu{Name: "统计师傅", Link: "/admin/reports/by_names", Ancestors: []string{"Reports Management"}})

	Admin.AddMenu(&admin.Menu{Name: "统计订单", Link: "/admin/reports/orders", Ancestors: []string{"Reports Management"}})

	// Admin.AddResource(AftersaleReportByName{}, &admin.Config{Menu: []string{"Reports Management"}, Priority: 1})
	// Admin.AddResource(AftersaleReportByBrand{}, &admin.Config{Menu: []string{"Reports Management"}, Priority: 2})
	// Admin.AddResource(AftersaleReportBySource{}, &admin.Config{Menu: []string{"Reports Management"}, Priority: 3})

	Admin.GetRouter().Get("/reports/by_sources", func(context *admin.Context) {
		// 	// do something here
		// 	// context.Request.URL.Query().Get(":name")
		// 	var bySourceContext = context.NewResourceContext("AftersaleReportBySource")
		// 	var result AftersaleReportBySource
		// 	db.DB.Raw("select to_char(date_trunc('day', aftersales.updated_at), 'YYYY') as time, source, count(*) as count, sum(fee) as sum from aftersales where state = 'completed' group by 1, source").Scan(&result)

		// 	bySourceContext.Render("index/table", result)
		context.Execute("by_sources", nil)
	})

	Admin.GetRouter().Get("/reports/by_brands", func(context *admin.Context) {
		context.Execute("by_brands", nil)
	})

	Admin.GetRouter().Get("/reports/by_names", func(context *admin.Context) {
		context.Execute("by_names", nil)
	})

	Admin.GetRouter().Get("/reports/orders", func(context *admin.Context) {
		context.Execute("orders", nil)
	})
	// bySource.FindManyHandler = func(results interface{}, context *qor.Context) error {
	// 	// find records and decode them to results

	// 	// var bySourceContext = context.NewResourceContext("AftersaleReportBySource")
	// 	var result AftersaleReportBySource
	// 	// db.Raw("SELECT name, age FROM users WHERE name = ?", 3).Scan(&result)
	// 	return db.DB.Raw("select to_char(date_trunc('day', aftersales.updated_at), 'YYYY') as time, source, count(*) as count, sum(fee) as sum from aftersales where state = 'completed' group by 1, source").Find(&result).Error

	// 	// return bySourceContext.Render("index/table", result)

	// }

	initFuncMap(Admin)
}
