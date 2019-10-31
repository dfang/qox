package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dfang/qor-demo/config"
	"github.com/dfang/qor-demo/config/db"
	"github.com/dfang/qor-demo/models/aftersales"
	"github.com/dfang/qor-demo/models/users"
	"github.com/gocraft/work"

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

	pool.PeriodicallyEnqueue(config.Config.Cron.ExpireAftersales, "expire_aftersales")
	pool.PeriodicallyEnqueue(config.Config.Cron.FreezeAuditedAftersales, "freeze_audited_aftersales")
	pool.PeriodicallyEnqueue(config.Config.Cron.UnfreezeAftersales, "unfreeze_aftersales")
	pool.PeriodicallyEnqueue(config.Config.Cron.UpdateBalances, "update_balances")

	pool.Job("expire_aftersales", ExpireAftersales)
	pool.Job("freeze_audited_aftersales", FreezeAftersales)
	pool.Job("unfreeze_aftersales", UnfreezeAftersales)
	pool.Job("update_balances", UpdateBalances)

	pool.Job("send_wechat_template_msg", SendWechatTemplateMsg)

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
	db.DB.Model(aftersales.Aftersale{}).Where("state = ?", "scheduled").Where("updated_at <= NOW() - INTERVAL '20 minutes'").Find(&items)
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
	// time.Sleep(70 * time.Second)
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
	// time.Sleep(55 * time.Second)
	log.Debug().Msg("unfreeze aftersales ......")

	var items []aftersales.Aftersale
	db.DB.Model(aftersales.Aftersale{}).Where("state = ?", "frozen").Find(&items)
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
		// 计算frozen_amount
		// 计算free_amount
		// update balance by user_id
		var balance aftersales.Balance
		db.DB.Model(aftersales.Balance{}).Where("user_id = ?", item.ID).Assign(aftersales.Balance{UserID: item.ID}).FirstOrInit(&balance)

		// select sum(amount) from settlements where user_id = 73 and state='frozen';

		// var frozenResult float32
		// var freeResult float32
		// db.DB.Table("settlements").Select("sum(amount)").Where("state = 'frozen'").Where("user_id = ?", item.ID).Take(&frozenResult)
		// db.DB.Table("settlements").Select("sum(amount)").Where("state = 'free'").Where("user_id = ?", item.ID).Take(&freeResult)
		type Result struct {
			State string
			Total float32
		}
		// rows, err :=
		var results []Result
		var f1, f2, f3 float32

		db.DB.Table("settlements").Select("state, sum(amount) as total").Group("state").Where("user_id = ?", item.ID).Scan(&results)
		for _, i := range results {
			// fmt.Println(i.State)
			// fmt.Println(i.Total)
			if i.State == "frozen" {
				f1 = i.Total
			}

			if i.State == "free" {
				f2 = i.Total
			}

			if i.State == "withdrawed" {
				f3 = i.Total
			}
		}

		balance.FrozenAmount = f1
		balance.FreeAmount = f2 + f3
		balance.WithdrawAmount = f3
		balance.TotalAmount = f2 + f1

		// balance.FrozenAmount = balance.FrozenAmount + balance.FreeAmount + balance.WithdrawAmount

		// balance.UserID = item.ID
		// balance.FrozenAmount = frozenResult
		// balance.FreeAmount = freeResult
		db.DB.Save(&balance)
	}

	log.Debug().Msgf("now is %s", time.Now().Format("2006-01-02 15:04:05"))
	log.Debug().Msg("update balances done ......")
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
