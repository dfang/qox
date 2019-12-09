package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/dfang/qor-demo/config"
	"github.com/dfang/qor-demo/config/db"
	"github.com/dfang/qor-demo/config/db/seeds"
	"github.com/dfang/qor-demo/models/aftersales"
	"github.com/dfang/qor-demo/models/users"
	"github.com/gocraft/work"
	"github.com/joho/sqltocsv"

	"github.com/rs/zerolog/log"
)

// just run `go startWorkerPool()` in main.go
// run workwebui -redis="redis:6379" -ns="qor" -listen=":5040"
// open localhost:5040 to view jobs ui
// https://crontab.guru/
// https://crontab.guru/examples.html
func startWorkerPool() {

	// Periodic Enqueueing (Cron)
	pool := work.NewWorkerPool(Context{}, 10, "qor", db.RedisPool)
	// pool.PeriodicallyEnqueue("30 * * * * *", "expire_aftersales") // This will enqueue a "expire_aftersales" job every minutes
	// pool.PeriodicallyEnqueue("30 * * * * *", "freeze_audited_aftersales")
	// pool.PeriodicallyEnqueue("5 * * * *", "unfreeze_aftersales")
	// pool.PeriodicallyEnqueue("30 * * * * *", "update_balances")
	pool.Middleware(Log)

	cronCfg := config.Config.Cron
	fmt.Println(cronCfg)

	pool.PeriodicallyEnqueue(cronCfg.ExpireAftersales, "expire_aftersales")
	pool.PeriodicallyEnqueue(cronCfg.FreezeAuditedAftersales, "freeze_audited_aftersales")
	pool.PeriodicallyEnqueue(cronCfg.UnfreezeAftersales, "unfreeze_aftersales")
	pool.PeriodicallyEnqueue(cronCfg.UpdateBalances, "update_balances")
	pool.PeriodicallyEnqueue(cronCfg.AutoExportMobilePhones, "export_mobile_phones")
	pool.PeriodicallyEnqueue(cronCfg.AutoExportOrderDetails, "export_order_details")
	pool.PeriodicallyEnqueue(cronCfg.AutoExportOrderFollowUps, "export_order_followups")

	if os.Getenv("QOR_ENV") != "production" && os.Getenv("DEMO_MODE") == "true" {
		pool.PeriodicallyEnqueue(cronCfg.AutoInquire, "auto_inquire")
		pool.PeriodicallyEnqueue(cronCfg.AutoSchedule, "auto_schedule")
		pool.PeriodicallyEnqueue(cronCfg.AutoProcess, "auto_process")
		pool.PeriodicallyEnqueue(cronCfg.AutoFinish, "auto_finish")
		pool.PeriodicallyEnqueue(cronCfg.AutoAudit, "auto_audit")
		pool.PeriodicallyEnqueue(cronCfg.AutoWithdraw, "auto_withdraw")
		pool.PeriodicallyEnqueue(cronCfg.AutoAward, "auto_award")
		pool.PeriodicallyEnqueue(cronCfg.AutoFine, "auto_fine")
		pool.PeriodicallyEnqueue(cronCfg.AutoGenerateAftersales, "auto_generate_aftersales")

		pool.Job("auto_inquire", AutoInquire)
		pool.Job("auto_schedule", AutoSchedule)
		pool.Job("auto_process", AutoProcess)
		pool.Job("auto_finish", AutoFinish)
		pool.Job("auto_audit", AutoAudit)
		pool.Job("auto_withdraw", AutoWithdraw)
		pool.Job("auto_award", AutoAward)
		pool.Job("auto_fine", AutoFine)
		pool.Job("auto_generate_aftersales", AutoGenerateAftersales)
	}

	pool.Job("expire_aftersales", ExpireAftersales)
	pool.Job("freeze_audited_aftersales", FreezeAftersales)
	pool.Job("unfreeze_aftersales", UnfreezeAftersales)
	pool.Job("update_balances", UpdateBalances)
	pool.Job("update_balance", UpdateBalance)

	pool.Job("send_wechat_template_msg", SendWechatTemplateMsg)
	pool.Job("export_mobile_phones", ExportMobilePhones)
	pool.Job("export_order_details", ExportOrderDetails)
	pool.Job("export_order_followups", ExportOrderFollowUps)

	// Start processing jobs
	pool.Start()
	// // Wait for a signal to quit:
	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt, os.Kill)
	// <-signalChan

	// // Stop the pool
	// pool.Stop()
}

// Context For gocraft/work
type Context struct {
	userID int64
}

// Log 开始执行任务的时候输出日志
func Log(job *work.Job, next work.NextMiddlewareFunc) error {
	log.Info().Msgf("Starting job: %s", job.Name)
	return next()
}

// ExpireAftersales 任务指派后 after_sale的状态为scheduled， 如果师傅20分钟之内没有响应，自动变为overdue状态
func ExpireAftersales(job *work.Job) error {
	// time.Sleep(10 * time.Second)
	log.Debug().Msgf("now is %s", time.Now().Format("2006-01-02 15:04:05"))
	log.Debug().Msg("expires all scheduled aftersales that idle for 20 minutes ......")
	var items []aftersales.Aftersale

	if os.Getenv("QOR_ENV") != "production" {
		db.DB.Model(aftersales.Aftersale{}).Where("state = ?", "scheduled").Where("updated_at <= NOW() - INTERVAL '2 minutes'").Find(&items)
	} else {
		db.DB.Model(aftersales.Aftersale{}).Where("state = ?", "scheduled").Where("updated_at <= NOW() - INTERVAL '20 minutes'").Find(&items)
	}
	// .Update("state", "overdue")
	for _, item := range items {
		log.Debug().Msgf("before expire: %s", item.State)
		aftersales.OrderStateMachine.Trigger("expire", &item, db.DB, "expires aftersale with id: "+fmt.Sprintf("%d", item.ID))
		log.Debug().Msgf("after expire: %s", item.State)
		// db.DB.Model(&item).Update("state", "overdue")
		db.DB.Save(&item)
	}

	log.Debug().Msgf("now is %s", time.Now().Format("2006-01-02 15:04:05"))
	log.Debug().Msg("expires aftersales done ")
	return nil
}

// FreezeAftersales 已审核的服务单冻结7天才能结算
func FreezeAftersales(job *work.Job) error {
	log.Debug().Msgf("now is %s", time.Now().Format("2006-01-02 15:04:05"))
	log.Debug().Msg("freeze aftersales ......")

	var items []aftersales.Aftersale
	// db.DB.Model(aftersales.Aftersale{}).Where("state = ?", "audited").Update("state", "frozen")
	db.DB.Model(aftersales.Aftersale{}).Where("state = ?", "audited").Find(&items)
	for _, item := range items {
		aftersales.OrderStateMachine.Trigger("freeze", &item, db.DB, "freeze aftersale with id: "+fmt.Sprintf("%d", item.ID))
		db.DB.Save(&item)
	}

	log.Debug().Msgf("now is %s", time.Now().Format("2006-01-02 15:04:05"))
	log.Debug().Msg("freeze aftersales done ......")
	return nil
}

// UnfreezeAftersales 解冻超过7天的，自动结算，金额算到师傅名下
func UnfreezeAftersales(job *work.Job) error {
	log.Debug().Msgf("now is %s", time.Now().Format("2006-01-02 15:04:05"))
	log.Debug().Msg("unfreeze aftersales ......")

	var items []aftersales.Aftersale
	if os.Getenv("QOR_ENV") != "production" {
		db.DB.Model(aftersales.Aftersale{}).Where("state = ?", "frozen").Where("updated_at <= NOW() - INTERVAL '2 minutes'").Find(&items)
	} else {
		db.DB.Model(aftersales.Aftersale{}).Where("state = ?", "frozen").Where("updated_at <= NOW() - INTERVAL '7 days'").Find(&items)
	}
	for _, item := range items {
		aftersales.OrderStateMachine.Trigger("unfreeze", &item, db.DB, "unfreeze aftersale with id: "+fmt.Sprintf("%d", item.ID))
		db.DB.Save(&item)
	}

	log.Debug().Msgf("now is %s", time.Now().Format("2006-01-02 15:04:05"))
	log.Debug().Msg("unfreeze aftersales done ......")
	return nil
}

// UpdateBalances 统计每个师傅的冻结金额和可结算金额并更新到Balances表
func UpdateBalances(job *work.Job) error {
	var workmen []users.User
	db.DB.Select("name, id").Where("role = ?", "workman").Find(&workmen)
	log.Debug().Msg("update balances ......")

	for _, item := range workmen {
		id := strconv.FormatUint(uint64(item.ID), 10)
		// balance := aftersales.UpdateBalanceFor(fmt.Sprint(item.ID))
		aftersales.UpdateBalanceFor(id)
		// balance := aftersales.UpdateBalanceFor(id)
		// db.DB.Save(&balance)
	}

	log.Debug().Msgf("now is %s", time.Now().Format("2006-01-02 15:04:05"))
	log.Debug().Msg("update balances done ......")
	return nil
}

// UpdateBalance 计算并更新某个师傅的账户额度
func UpdateBalance(job *work.Job) error {
	userID := job.ArgString("user_id")
	balance := aftersales.UpdateBalanceFor(userID)
	db.DB.Save(&balance)
	log.Debug().Msg("update balance done ......")

	return nil
}

// SendWechatTemplateMsg 发送微信模版消息（当任务指派给师傅或者订单解冻了，需要给师傅推送一条微信模板消息)
func SendWechatTemplateMsg(job *work.Job) error {
	c := job.ArgString("contents")
	contents := []byte(c)
	tkn := job.ArgString("token")
	// https://mp.weixin.qq.com/advanced/tmplmsg?action=faq&token=1545649248&lang=zh_CN
	// tkn := getToken()
	url := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + tkn
	fmt.Println("URL:>", url)

	// var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	// req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(contents))
	fmt.Println("Request Body:> ", string(contents))

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))
	var rsp TemplateMsgResp
	if err := json.NewDecoder(resp.Body).Decode(&rsp); err != nil {
		log.Printf("Error decoding body: %s", err)
	}

	fmt.Printf("%+v\n", rsp)

	if rsp.ErrCode != 0 || rsp.ErrMsg != "ok" {
		// return errors.New(fmt.Sprintf("发送模板消息失败, errcode: %d, errmsg: %s", rsp.ErrCode, rsp.ErrMsg))
		return fmt.Errorf("发送模板消息失败, errcode: %d, errmsg: %s", rsp.ErrCode, rsp.ErrMsg)
	}

	return nil
}

// UnclutterOldNotifications 干掉太久的已读通知
func UnclutterOldNotifications(job *work.Job) error {
	// TODO: to implement
	return nil
}

// AlertOverdueAftersales 超时告警
func AlertOverdueAftersales(job *work.Job) error {
	// TODO: chrome web push notifications
	return nil
}

// AlertToAudit 待审核告警
func AlertToAudit(job *work.Job) error {
	// TODO: chrome web push notifications
	return nil
}

// TemplateMsgResp 发送模板消息返回结果
type TemplateMsgResp struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	MsgID   int64  `json:"msgid"`
}

/* DEMO_MODE=true 下自动化任务 */

// AutoInquire demo模式下自动预约
func AutoInquire(job *work.Job) error {
	log.Debug().Msg("demo模式下自动预约 .........")

	if os.Getenv("QOR_ENV") != "production" && os.Getenv("DEMO_MODE") == "true" {
		var a aftersales.Aftersale
		db.DB.Where("state = ?", "created").Order("random()").First(&a)

		if a.ID > 0 {
			aftersales.OrderStateMachine.Trigger("inquire", &a, db.DB, "auto inquire aftersale with id: "+fmt.Sprintf("%d", a.ID))
			a.Remark = "客户很急"
			db.DB.Save(&a)
		}
	}

	return nil
}

// AutoSchedule demo模式下自动派单
func AutoSchedule(job *work.Job) error {
	log.Debug().Msg("demo模式下自动派单 .........")

	if os.Getenv("QOR_ENV") != "production" && os.Getenv("DEMO_MODE") == "true" {
		var a aftersales.Aftersale
		var w users.User
		db.DB.Model(users.User{}).Where("role = ?", "workman").Order("random()").First(&w)
		db.DB.Where("state = ? or state = ?", "overdue", "inquired").Order("random()").First(&a)

		if a.ID > 0 && w.ID > 0 {
			a.UserID = w.ID
			if a.State == "inquired" {
				aftersales.OrderStateMachine.Trigger("schedule", &a, db.DB, "auto schedule aftersale with id: "+fmt.Sprintf("%d", a.ID))
			} else {
				aftersales.OrderStateMachine.Trigger("reschedule", &a, db.DB, "auto reschedule aftersale with id: "+fmt.Sprintf("%d", a.ID))
			}
			db.DB.Save(&a)
		}
	}

	return nil
}

// AutoProcess demo模式下自动接单
func AutoProcess(job *work.Job) error {
	log.Debug().Msg("demo模式下自动接单 .........")

	if os.Getenv("QOR_ENV") != "production" && os.Getenv("DEMO_MODE") == "true" {
		var a aftersales.Aftersale
		db.DB.Where("state = ?", "scheduled").Order("random()").First(&a)
		if a.ID > 0 {
			a.State = "processing"
			db.DB.Save(&a)
		}
	}
	return nil
}

// AutoFinish demo模式下自动完成
func AutoFinish(job *work.Job) error {
	log.Debug().Msg("demo模式下自动完成 .........")

	if os.Getenv("QOR_ENV") != "production" && os.Getenv("DEMO_MODE") == "true" {
		var a aftersales.Aftersale
		db.DB.Where("state = ?", "processing").Order("random()").First(&a)
		if a.ID > 0 {
			a.State = "processed"
			db.DB.Save(&a)
		}
	}
	return nil
}

// AutoAudit demo模式下自动审批
func AutoAudit(job *work.Job) error {
	log.Debug().Msg("demo模式下自动审批 .........")

	if os.Getenv("QOR_ENV") != "production" && os.Getenv("DEMO_MODE") == "true" {
		var a aftersales.Aftersale
		db.DB.Where("state = ?", "processed").Order("random()").First(&a)
		if a.ID > 0 {
			aftersales.OrderStateMachine.Trigger("audit", &a, db.DB, "auto audit aftersale with id: "+fmt.Sprintf("%d", a.ID))
			db.DB.Save(&a)
		}
	}

	return nil
}

// AutoWithdraw demo模式下自动提现
func AutoWithdraw(job *work.Job) error {
	log.Debug().Msg("demo模式下自动提现 .........")

	if os.Getenv("QOR_ENV") != "production" && os.Getenv("DEMO_MODE") == "true" {
		var w users.User
		db.DB.Model(users.User{}).Where("role = ?", "workman").Order("random()").First(&w)

		if w.ID > 0 {
			s := aftersales.Settlement{
				Direction: "提现",
				Amount:    -5,
				UserID:    w.ID,
			}

			if err := db.DB.Save(&s); err != nil {
				fmt.Println("提现失败！！！！")
			}
		}
	}
	return nil
}

// AutoAward demo模式下自动增加奖金
func AutoAward(job *work.Job) error {
	log.Debug().Msg("demo模式下自动增加奖金 .........")
	if os.Getenv("QOR_ENV") != "production" && os.Getenv("DEMO_MODE") == "true" {
		var w users.User
		db.DB.Model(users.User{}).Where("role = ?", "workman").Order("random()").First(&w)

		if w.ID > 0 {
			s := aftersales.Settlement{
				Direction: "奖励",
				Amount:    20,
				UserID:    w.ID,
			}

			if err := db.DB.Save(&s); err != nil {
				fmt.Println("奖励失败！！！！")
			}
		}
	}
	return nil
}

// AutoFine demo模式下自动罚款
func AutoFine(job *work.Job) error {
	log.Debug().Msg("demo模式下自动罚款 .........")
	if os.Getenv("QOR_ENV") != "production" && os.Getenv("DEMO_MODE") == "true" {
		var w users.User
		db.DB.Model(users.User{}).Where("role = ?", "workman").Order("random()").First(&w)

		if w.ID > 0 {
			s := aftersales.Settlement{
				Direction: "罚款",
				Amount:    -15,
				UserID:    w.ID,
			}

			if err := db.DB.Save(&s); err != nil {
				fmt.Println("罚款失败！！！！")
			}
		}
	}
	return nil
}

// AutoGenerateAftersales demo模式下自动创建服务单
func AutoGenerateAftersales(job *work.Job) error {
	log.Debug().Msg("demo模式下自动创建服务单 .........")
	if os.Getenv("QOR_ENV") != "production" && os.Getenv("DEMO_MODE") == "true" {
		seeds.CreateAftersales()
	}
	return nil
}

// ExportMobilePhones 定时导出昨日电话号码
func ExportMobilePhones(job *work.Job) error {
	for i := -7; i < 0; i++ {
		fileName := fmt.Sprintf("public/downloads/电话号码/%v.phones.csv", time.Now().AddDate(0, 0, i).Format("20060102"))
		fmt.Println(filepath.Dir(fileName))
		err := os.MkdirAll(filepath.Dir(fileName), 0777)
		if err != nil {
			fmt.Printf("mkdirall err : %v\n", err)
		}
		// bomUtf8 := []byte{0xEF, 0xBB, 0xBF}
		f, err := os.Create(fileName)
		defer f.Close()
		// f.Write(bomUtf8)
		if err != nil {
			panic(err)
		}
		rows, err := db.DB.DB().Query(`select distinct customer_phone from orders where DATE(created_at) = DATE(timestamp 'yesterday');`)
		if err != nil {
			panic(err)
		}

		defer rows.Close()

		sqltocsv.Write(f, rows)
	}

	return nil
}

// ExportOrderDetails 定时导出订单详情
func ExportOrderDetails(job *work.Job) error {

	for i := -60; i < 0; i++ {
		n := time.Now()
		m := n.AddDate(0, 0, i)
		fmt.Println(m)

		fileName := fmt.Sprintf("public/downloads/订单/%v.orders.csv", m.Format("20060102"))
		err := os.MkdirAll(filepath.Dir(fileName), 0777)
		if err != nil {
			fmt.Printf("mkdirall err : %v\n", err)
		}
		bomUtf8 := []byte{0xEF, 0xBB, 0xBF}
		f, err := os.Create(fileName)
		defer f.Close()
		f.Write(bomUtf8)
		if err != nil {
			panic(err)
		}
		rows, err := db.DB.DB().Query(`
    SELECT to_char(orders.created_at, 'YYYY-MM-DD') AS 预约日, orders.order_no AS 订单号, customer_name AS 客户姓名, customer_phone AS 客户手机, customer_address AS 客户地址, order_items.item_name AS 商品, order_items.quantity AS 数量
    FROM orders
    INNER JOIN order_items
    ON orders.order_no = order_items.order_no
    WHERE to_char(orders.created_at, 'YYYY-MM-DD') = $1
    ORDER BY orders.created_at`, m.Format("2006-01-02"))

		if err != nil {
			panic(err)
		}

		defer rows.Close()

		sqltocsv.Write(f, rows)
	}

	return nil
}

// ExportOrderFollowUps 定时导出订单回访
func ExportOrderFollowUps(job *work.Job) error {
	for i := -7; i < 0; i++ {
		n := time.Now()
		m := n.AddDate(0, 0, i)
		fmt.Println(m)

		fileName := fmt.Sprintf("public/downloads/订单回访/%v.order_followups.csv", m.Format("20060102"))
		err := os.MkdirAll(filepath.Dir(fileName), 0777)
		if err != nil {
			fmt.Printf("mkdirall err : %v\n", err)
		}

		bomUtf8 := []byte{0xEF, 0xBB, 0xBF}
		f, err := os.Create(fileName)
		defer f.Close()
		f.Write(bomUtf8)
		if err != nil {
			panic(err)
		}
		rows, err := db.DB.DB().Query(`
      SELECT order_follow_ups.order_no AS 订单号, order_items.item_name AS 产品信息, orders.customer_address AS 客户地址, orders.customer_phone AS 客户电话, satisfaction_of_timeliness AS 对时效是否满意, satisfaction_of_services AS  对服务是否满意, inspect_the_goods AS 是否开箱验货, request_feedback AS  是否邀评, leave_contact_infomation AS 是否留下联系方式, introduce_warranty_extension AS 是否有介绍延保, position_properly AS 是否把商品放到指定位置, feedback AS 有无问题反馈, exception_handling AS 异常处理结果,  order_follow_ups.created_at AS 创建时间 
      FROM order_follow_ups INNER JOIN orders 
      ON orders.order_no = order_follow_ups.order_no 
      INNER JOIN order_items 
      ON orders.order_no = order_items.order_no
      WHERE to_char(order_follow_ups.created_at, 'YYYY-MM-DD') = $1
      ORDER BY order_follow_ups.created_at`,
			m.Format("2006-01-02"))

		if err != nil {
			panic(err)
		}

		defer rows.Close()

		sqltocsv.Write(f, rows)
	}

	return nil
}
