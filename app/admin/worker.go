package admin

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dfang/qor-demo/config/db"
	"github.com/dfang/qor-demo/config/i18n"
	"github.com/joho/sqltocsv"
	"github.com/qor/admin"
	"github.com/qor/i18n/exchange_actions"
	"github.com/qor/worker"
)

// SetupWorker setup worker
func SetupWorker(Admin *admin.Admin) {
	Worker := worker.New()

	type sendNewsletterArgument struct {
		Subject      string
		Content      string `sql:"size:65532"`
		SendPassword string
		worker.Schedule
	}

	Worker.RegisterJob(&worker.Job{
		Name: "Send Newsletter",
		Handler: func(argument interface{}, qorJob worker.QorJobInterface) error {
			qorJob.AddLog("Started sending newsletters...")
			qorJob.AddLog(fmt.Sprintf("Argument: %+v", argument.(*sendNewsletterArgument)))
			for i := 1; i <= 100; i++ {
				time.Sleep(100 * time.Millisecond)
				qorJob.AddLog(fmt.Sprintf("Sending newsletter %v...", i))
				qorJob.SetProgress(uint(i))
			}
			qorJob.AddLog("Finished send newsletters")
			return nil
		},
		Resource: Admin.NewResource(&sendNewsletterArgument{}),
	})

	// type importProductArgument struct {
	// 	File oss.OSS
	// }

	// Worker.RegisterJob(&worker.Job{
	// 	Name:  "Import Products",
	// 	Group: "Products Management",
	// 	Handler: func(arg interface{}, qorJob worker.QorJobInterface) error {
	// 		argument := arg.(*importProductArgument)

	// 		context := &qor.Context{DB: db.DB}

	// 		var errorCount uint

	// 		if err := ProductExchange.Import(
	// 			csv.New(filepath.Join("public", argument.File.URL())),
	// 			context,
	// 			func(progress exchange.Progress) error {
	// 				var cells = []worker.TableCell{
	// 					{Value: fmt.Sprint(progress.Current)},
	// 				}

	// 				var hasError bool
	// 				for _, cell := range progress.Cells {
	// 					var tableCell = worker.TableCell{
	// 						Value: fmt.Sprint(cell.Value),
	// 					}

	// 					if cell.Error != nil {
	// 						hasError = true
	// 						errorCount++
	// 						tableCell.Error = cell.Error.Error()
	// 					}

	// 					cells = append(cells, tableCell)
	// 				}

	// 				if hasError {
	// 					if errorCount == 1 {
	// 						var headerCells = []worker.TableCell{
	// 							{Value: "Line No."},
	// 						}
	// 						for _, cell := range progress.Cells {
	// 							headerCells = append(headerCells, worker.TableCell{
	// 								Value: cell.Header,
	// 							})
	// 						}
	// 						qorJob.AddResultsRow(headerCells...)
	// 					}

	// 					qorJob.AddResultsRow(cells...)
	// 				}

	// 				qorJob.SetProgress(uint(float32(progress.Current) / float32(progress.Total) * 100))
	// 				qorJob.AddLog(fmt.Sprintf("%d/%d Importing product %v", progress.Current, progress.Total, progress.Value.(*products.Product).Code))
	// 				return nil
	// 			},
	// 		); err != nil {
	// 			qorJob.AddLog(err.Error())
	// 		}

	// 		return nil
	// 	},
	// 	Resource: Admin.NewResource(&importProductArgument{}),
	// })

	// Worker.RegisterJob(&worker.Job{
	// 	Name:  "Export Products",
	// 	Group: "Products Management",
	// 	Handler: func(arg interface{}, qorJob worker.QorJobInterface) error {
	// 		qorJob.AddLog("Exporting products...")

	// 		context := &qor.Context{DB: db.DB}
	// 		fileName := fmt.Sprintf("/downloads/products/%v.xlsx", time.Now().UnixNano())
	// 		if err := ProductExchange.Export(

	// 			csv.New(filepath.Join("public", fileName)),

	// 			context,
	// 			func(progress exchange.Progress) error {
	// 				qorJob.AddLog(fmt.Sprintf("%v/%v Exporting product %v", progress.Current, progress.Total, progress.Value.(*products.Product).Code))
	// 				return nil
	// 			},
	// 		); err != nil {
	// 			qorJob.AddLog(err.Error())
	// 		}

	// 		qorJob.SetProgressText(fmt.Sprintf("<a href='%v'>Download exported products</a>", fileName))
	// 		return nil
	// 	},
	// })

	Worker.RegisterJob(&worker.Job{
		Name:  "Export Orders",
		Group: "Orders Management",
		Handler: func(arg interface{}, qorJob worker.QorJobInterface) error {
			fmt.Println("fuck")
			qorJob.AddLog("Exporting orders...")
			// context := &qor.Context{DB: db.DB}

			// 导出csv文件 中文乱码问题
			// https://forum.golangbridge.org/t/how-to-write-csv-file-with-bom-utf8/9434
			// https://www.zhihu.com/question/21869078
			// https://blog.csdn.net/wodatoucai/article/details/46970347
			// https://pathbox.github.io/2017/01/20/csv-operation-in-GO/
			// https://stackoverflow.com/questions/21371673/reading-files-with-a-bom-in-go
			// fileName := fmt.Sprintf("/downloads/orders/%v.csv", time.Now().UnixNano())
			fileName := fmt.Sprintf("/downloads/orders/%v.csv", time.Now().AddDate(0, 0, -1).Format("20060102"))
			// bomUtf8 := []byte{0xEF, 0xBB, 0xBF}
			f, err := os.Create(filepath.Join("public", fileName))
			defer f.Close()
			// f.Write(bomUtf8)
			if err != nil {
				panic(err)
			}

			// // dec := encoding.Encoding.UTF8
			// dec := unicode.UTF8.NewDecoder()
			// dec := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder()
			// // transform.NewWriter()
			// writer := transform.NewWriter(f, dec.Transformer)

			// if err := OrderExchange.Export(
			// 	// csv.New(filepath.Join("public", fileName)),
			// 	csv.New(f),
			// 	context,
			// 	func(progress exchange.Progress) error {
			// 		qorJob.AddLog(fmt.Sprintf("%v/%v Exporting order %v", progress.Current, progress.Total, progress.Value.(*orders.Order).OrderNo))
			// 		return nil
			// 	},
			// ); err != nil {
			// 	qorJob.AddLog(err.Error())
			// }

			rows, err := db.DB.DB().Query("SELECT * FROM orders_view where DATE(created_at) = DATE(timestamp 'yesterday');")
			if err != nil {
			}

			defer rows.Close()

			// context.Writer.Header().Set("Content-type", "text/csv")
			// context.Writer.Header().Set("Content-Disposition", "attachment; filename=\"report.csv\"")
			sqltocsv.Write(f, rows)

			// err = sqltocsv.WriteFile(fileName, rows)
			// if err != nil {
			// 	panic(err)
			// }

			qorJob.SetProgressText(fmt.Sprintf("<a href='%v'>Download exported orders</a>", fileName))
			return nil
		},
	})

	Worker.RegisterJob(&worker.Job{
		Name:  "Export_FollowUps",
		Group: "Orders Management",
		Handler: func(arg interface{}, qorJob worker.QorJobInterface) error {
			qorJob.AddLog("导出订单回访记录...")
			// context := &qor.Context{DB: db.DB}
			// 导出csv文件 中文乱码问题
			// https://forum.golangbridge.org/t/how-to-write-csv-file-with-bom-utf8/9434
			// https://www.zhihu.com/question/21869078
			// https://blog.csdn.net/wodatoucai/article/details/46970347
			// https://pathbox.github.io/2017/01/20/csv-operation-in-GO/
			// https://stackoverflow.com/questions/21371673/reading-files-with-a-bom-in-go
			// fileName := fmt.Sprintf("/downloads/orders/%v.csv", time.Now().UnixNano())
			fileName := fmt.Sprintf("/downloads/orders/%v.csv", time.Now().AddDate(0, 0, -1).Format("20060102"))
			bomUtf8 := []byte{0xEF, 0xBB, 0xBF}
			f, err := os.Create(filepath.Join("public", fileName))
			defer f.Close()
			f.Write(bomUtf8)
			if err != nil {
				panic(err)
			}

			// rows, err := db.DB.DB().Query("SELECT * FROM order_follow_ups where DATE(created_at) = DATE(timestamp 'yesterday');")
			rows, err := db.DB.DB().Query("SELECT order_follow_ups.order_no AS 订单号, order_items.item_name AS 产品信息,orders.customer_address AS 客户地址, orders.customer_phone AS 客户电话, satisfaction_of_timeliness AS 对时效是否满意, satisfaction_of_services AS  对服务是否满意, inspect_the_goods AS 是否开箱验货, request_feedback AS  是否邀评, leave_contact_infomation AS 是否留下联系方式, introduce_warranty_extension AS 是否有介绍延保, position_properly AS 是否把商品放到指定位置, feedback AS 有无问题反馈, exception_handling AS 异常处理结果,  order_follow_ups.created_at AS 创建时间 FROM order_follow_ups INNER JOIN orders ON orders.order_no = order_follow_ups.order_no INNER JOIN order_items ON orders.order_no = order_items.order_no;")
			if err != nil {
			}

			defer rows.Close()
			// context.Writer.Header().Set("Content-type", "text/csv")
			// context.Writer.Header().Set("Content-Disposition", "attachment; filename=\"report.csv\"")
			sqltocsv.Write(f, rows)

			// err = sqltocsv.WriteFile(fileName, rows)
			// if err != nil {
			// 	panic(err)
			// }

			qorJob.SetProgressText(fmt.Sprintf("<a href='%v'>下载导出的订单回访记录</a>", fileName))
			return nil
		},
	})

	exchange_actions.RegisterExchangeJobs(i18n.I18n, Worker)
	Admin.AddResource(Worker, &admin.Config{Menu: []string{"Site Management"}, Priority: 3})
}
