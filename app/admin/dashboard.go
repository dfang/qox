package admin

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dfang/qor-demo/config/db"
	"github.com/jinzhu/now"
	"github.com/qor/admin"
)

type Chart struct {
	Total string
	Date  string
}

/*
date format 2015-01-23
*/
func GetChartDataByDay(table, start, end string) (res []Chart) {
	startdate, err := now.Parse(start)
	if err != nil {
		return
	}

	enddate, err := now.Parse(end)
	if err != nil || enddate.UnixNano() < startdate.UnixNano() {
		enddate = now.EndOfDay()
	} else {
		enddate = enddate.AddDate(0, 0, 1)
	}

	db.DB.Table(table).
		Where("created_at > ? AND created_at < ?", startdate, enddate).
		Select("date(created_at) as date, count(*) as total").
		Group("date(created_at)").
		Order("date(created_at)").
		Scan(&res)
	return
}

func GetChartDataByMonth(table, start, end string) (res []Chart) {
	startdate, err := now.Parse(start)
	if err != nil {
		return
	}

	enddate, err := now.Parse(end)
	if err != nil || enddate.UnixNano() < startdate.UnixNano() {
		panic("endDate must > startDate")
	}

	// select
	//   to_char(date_trunc('month', created_at), 'YYYY-MM'),
	//   count(1)
	// from orders
	// group by 1
	// order by 1;

	// select
	// state,
	// count(1)
	// from after_sales
	// where DATE(created_at) = DATE(timestamp 'yesterday')
	// group by state;

	db.DB.Table(table).
		Where("created_at > ? AND created_at < ?", startdate, enddate).
		Select("to_char(date_trunc('month', created_at), 'YYYY-MM') as date, count(1) as total").
		Group("date").
		Order("date").
		Scan(&res)
	return
}

type Charts struct {
	Orders []Chart
	Users  []Chart
}

func ReportsDataHandler(context *admin.Context) {
	charts := &Charts{}
	startDate := context.Request.URL.Query().Get("startDate")
	endDate := context.Request.URL.Query().Get("endDate")

	// startDate -> 2019-08-01 按天
	// startDate -> 2019-08 按月
	if len(strings.Split(startDate, "-")) == 3 {
		charts.Orders = GetChartDataByDay("orders", startDate, endDate)
		// charts.Users = GetChartDataByDay("users", startDate, endDate)
	} else {
		charts.Orders = GetChartDataByMonth("orders", startDate, endDate)
	}

	b, _ := json.Marshal(charts)
	context.Writer.Write(b)

	return
}

// SetupDashboard setup dashboard
func SetupDashboard(Admin *admin.Admin) {
	// Add Dashboard
	Admin.AddMenu(&admin.Menu{Name: "Dashboard", Link: "/admin", Priority: 1})

	// Admin.AddMenu(&admin.Menu{Name: "Today", Link: "/today", Priority: 1})
	// Admin.AddMenu(&admin.Menu{Name: "Today", Link: "/admin/today", Ancestors: []string{"Reports"}})

	// Admin.GetRouter().Get("/reports", ReportsDataHandler)
	Admin.GetRouter().Get("/reports", ReportsHandler)
	initFuncMap(Admin)
}

func ReportsHandler(context *admin.Context) {
	http.Redirect(context.Writer, context.Request, "/admin/reports/by_brands", http.StatusFound)
}
