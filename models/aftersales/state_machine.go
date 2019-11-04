package aftersales

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/dfang/qor-demo/config/db"
	"github.com/dfang/qor-demo/models/users"

	"github.com/jinzhu/gorm"
	"github.com/qor/transition"

	"github.com/gocraft/work"
)

var (
	// OrderStateMachine order's state machine
	OrderStateMachine = transition.New(&Aftersale{})

	// SettlementStateMachine 结算状态机
	SettlementStateMachine = transition.New(&Settlement{})
)

// State
// 括号里是Action

// (created)----> created 已接收状态 -----(inquire)---> inquired 未预约状态
// -----（schedule）----> scheduled 已预约  ---(confirm_schedule)--> schedule_confirmed 待处理
// -----(confirm_complete)----> to_be_audited 待审核 ---audit--> audited -----> finalized 锁定

// STATES  Dashboard for Operators
// 待预约   1 （aftersales.state == "created"）
// 待指派   1  (aftersales.state == "inquired"）
// 超时任务单 1  (aftersales.预约时间== "空" 需要重新指派）
// 待审核   1  (aftersales.state == "to_be_audited"）
var STATES = []string{

	"created", // 建单之后

	"inquired", // 信息员给用户打电话预约之后

	"scheduled", // 指派师傅之后

	"overdue", // 指派师傅之后, 师傅未给用户打电话预约超时了

	"processing", // 师傅给用户打过电话，确认了上门时间的状态

	"processed", // 师傅上传了照片，等待审核

	"audited", // 审核通过

	"audit_failed", // 审核不通过

	"frozen", // 冻结， 审核通过后，结算金额需冻结7天

	"completed", // 完成，解冻之后的状态
}

func init() {
	SettlementStateMachine.Initial("frozen")
	SettlementStateMachine.State("free")
	SettlementStateMachine.State("completed")
	SettlementStateMachine.Event("unfreeze").To("free").From("frozen")

	// Define Order's States
	OrderStateMachine.Initial("created")
	OrderStateMachine.State("inquired")
	OrderStateMachine.State("scheduled")
	OrderStateMachine.State("overdue")
	OrderStateMachine.State("processing")
	OrderStateMachine.State("processed")
	OrderStateMachine.State("audited")
	OrderStateMachine.State("audit_failed")
	OrderStateMachine.State("frozen")
	OrderStateMachine.State("completed")

	// 和用户预约大概时间
	OrderStateMachine.Event("inquire").To("inquired").From("created")

	// 指派师傅
	OrderStateMachine.Event("schedule").To("scheduled").From("inquired").After(func(value interface{}, tx *gorm.DB) (err error) {
		err = enqueueJob(value, tx, "schedule")
		if err != nil {
			panic("oooops")
		}

		return nil
	})

	// 重新指派师傅
	OrderStateMachine.Event("reschedule").To("scheduled").From("scheduled").After(func(value interface{}, tx *gorm.DB) (err error) {
		err = enqueueJob(value, tx, "schedule")
		if err != nil {
			panic("oooops")
		}

		return nil
	})

	// 重新指派师傅
	OrderStateMachine.Event("reschedule").To("scheduled").From("overdue").After(func(value interface{}, tx *gorm.DB) (err error) {
		err = enqueueJob(value, tx, "schedule")
		if err != nil {
			panic("oooops")
		}

		return nil
	})

	// 师傅接单
	OrderStateMachine.Event("take_order").To("processing").From("scheduled").After(func(value interface{}, tx *gorm.DB) (err error) {
		err = enqueueJob(value, tx, "take_order")
		if err != nil {
			panic("oooops")
		}

		return nil
	})

	// 超时未响应的订单自动过期
	OrderStateMachine.Event("expire").To("overdue").From("scheduled")

	// 根据师傅上传的照片 审核服务是否完成
	OrderStateMachine.Event("audit").To("audited").From("processed").After(func(value interface{}, tx *gorm.DB) (err error) {
		err = enqueueJob(value, tx, "audit")
		if err != nil {
			panic("oooops")
		}
		return nil
	})

	// 审核不通过
	OrderStateMachine.Event("audit_failed").To("audit_failed").From("processed").After(func(value interface{}, tx *gorm.DB) (err error) {
		err = enqueueJob(value, tx, "audit_failed")
		if err != nil {
			panic("oooops")
		}
		return nil
	})

	// 冻结
	OrderStateMachine.Event("freeze").To("frozen").From("audited").After(func(value interface{}, tx *gorm.DB) (err error) {
		item := value.(*Aftersale)
		// 结算里加一笔,
		if item.Fee > 0 {
			// 需要判断
			s := Settlement{
				UserID:      item.UserID,
				Amount:      item.Fee,
				Direction:   "收入",
				AftersaleID: item.ID,
			}

			tx.Model(Settlement{}).Save(&s)
		}

		return nil
	})

	// 解冻
	OrderStateMachine.Event("unfreeze").To("completed").From("frozen").After(func(value interface{}, tx *gorm.DB) (err error) {
		item := value.(*Aftersale)
		var settlement Settlement
		tx.Model(Settlement{}).Where("aftersale_id = ?", item.ID).Find(&settlement)

		SettlementStateMachine.Trigger("unfreeze", &settlement, tx, "unfreeze aftersale with id: "+fmt.Sprintf("%d", item.ID))
		tx.Model(Settlement{}).Save(&settlement)

		err = enqueueJob(value, tx, "unfreeze")
		if err != nil {
			panic("oooops")
		}

		return nil
	})
}

func enqueueJob(value interface{}, tx *gorm.DB, jobType string) error {
	// 首先找出openid 等需要填充到模板消息变量的信息
	item := value.(*Aftersale)
	u := users.User{}
	tx.Where("id = ?", item.UserID).First(&u)
	mobilePhone := u.MobilePhone
	wp := users.WechatProfile{}
	tx.Where("mobile_phone = ?", mobilePhone).First(&wp)
	fmt.Println("user_id is > ", item.UserID)
	fmt.Println("mobile_phone is > ", mobilePhone)
	fmt.Println("openid is > ", wp.Openid)

	if wp.Openid == "" {
		return fmt.Errorf("openid 不能为空，否则无法发送模板消息")
	}

	t, _ := TimeIn(time.Now(), "Asia/Shanghai")
	format := t.Format("2006-01-02 15:04:05")

	switch jobType {
	case "schedule":
		enqueueScheduleJob(wp.Openid, item.ID, format)
		return nil
	case "audit":
		enqueueAuditJob(wp.Openid, item.ID, format)
		return nil
	case "audit_failed":
		enqueueAuditFailedJob(wp.Openid, item.ID, format)
		return nil
	case "unfreeze":
		enqueueUnfreezeJob(wp.Openid, item.ID, format)
		return nil
	case "take_order":
		// enqueueTakeOrderJob(wp.Openid, item.ID, format)
		return nil
	default:
		fmt.Println("Unknown job type")
		return nil
	}
}

func enqueueScheduleJob(openid string, id uint, format string) error {
	// 推送微信模版消息到师傅微信
	// Make an enqueuer with a particular namespace
	var enqueuer = work.NewEnqueuer("qor", db.RedisPool)
	// _, err = enqueuer.Enqueue("send_wechat_template_msg", work.Q{"address": "test@example.com", "subject": "hello world", "customer_id": 4})

	// 填充模板
	m := ModelForSchedule{
		OpenID: openid,
		ID:     strconv.FormatUint(uint64(id), 10),
		URL:    "http://mp.xsjd123.com/",
		Now:    format,
	}
	s := executeTpl(SCHEDULE_TPL, m)
	// fmt.Println("模板消息插入变量后是: ", s)

	tkn := getToken()
	// sendTemplateMsg(b)
	fmt.Println("enqueueing send_wechat_template_msg .....")
	j, err := enqueuer.Enqueue("send_wechat_template_msg", work.Q{"contents": s, "token": tkn})
	if err != nil {
		return err
	}
	fmt.Printf("Job ID: %s, Name: %s, EnqueueAt: %d\n", j.ID, j.Name, j.EnqueuedAt)
	return nil
}

func enqueueAuditJob(openid string, id uint, format string) error {
	// 推送微信模版消息到师傅微信
	// Make an enqueuer with a particular namespace
	var enqueuer = work.NewEnqueuer("qor", db.RedisPool)

	// 填充模板
	m := ModelForAudit{
		OpenID: openid,
		ID:     strconv.FormatUint(uint64(id), 10),
		URL:    "http://mp.xsjd123.com/",
		Now:    format,
	}
	s := executeTpl(AUDIT_OK_TPL, m)
	// fmt.Println("模板消息插入变量后是: ", s)

	tkn := getToken()
	// sendTemplateMsg(b)
	fmt.Println("enqueueing send_wechat_template_msg .....")
	j, err := enqueuer.Enqueue("send_wechat_template_msg", work.Q{"contents": s, "token": tkn})
	if err != nil {
		return err
	}
	fmt.Printf("Job ID: %s, Name: %s, EnqueueAt: %d\n", j.ID, j.Name, j.EnqueuedAt)
	return nil
}

func enqueueUnfreezeJob(openid string, id uint, format string) error {
	// 推送微信模版消息到师傅微信
	var enqueuer = work.NewEnqueuer("qor", db.RedisPool)
	// 填充模板
	m := ModelForAudit{
		OpenID: openid,
		ID:     strconv.FormatUint(uint64(id), 10),
		URL:    "http://mp.xsjd123.com/",
		Now:    format,
	}
	s := executeTpl(UNFREEZE_TPL, m)
	// fmt.Println("模板消息插入变量后是: ", s)

	tkn := getToken()
	// sendTemplateMsg(b)
	fmt.Println("enqueueing send_wechat_template_msg .....")
	j, err := enqueuer.Enqueue("send_wechat_template_msg", work.Q{"contents": s, "token": tkn})
	if err != nil {
		return err
	}
	fmt.Printf("Job ID: %s, Name: %s, EnqueueAt: %d\n", j.ID, j.Name, j.EnqueuedAt)
	return nil
}

func enqueueAuditFailedJob(openid string, id uint, format string) error {
	// 推送微信模版消息到师傅微信
	var enqueuer = work.NewEnqueuer("qor", db.RedisPool)
	// 填充模板
	m := ModelForAudit{
		OpenID: openid,
		ID:     strconv.FormatUint(uint64(id), 10),
		URL:    "http://mp.xsjd123.com/",
		Now:    format,
	}
	s := executeTpl(AUDIT_FAILED_TPL, m)
	// fmt.Println("模板消息插入变量后是: ", s)

	tkn := getToken()
	// sendTemplateMsg(b)
	fmt.Println("enqueueing send_wechat_template_msg .....")
	j, err := enqueuer.Enqueue("send_wechat_template_msg", work.Q{"contents": s, "token": tkn})
	if err != nil {
		return err
	}
	fmt.Printf("Job ID: %s, Name: %s, EnqueueAt: %d\n", j.ID, j.Name, j.EnqueuedAt)
	return nil
}

func enqueueTakeOrderJob(value interface{}, tx *gorm.DB) error {
	fmt.Println("师傅接单了， 给用户发通知(模板消息或短信")
	// TODO: 发送通知
	return nil
}

// 从中控服务器获取access_token
func getToken() string {
	var tkn Token
	resp, err := http.Get("http://wx.xsjd123.com/access_token")
	if err != nil {
		log.Fatal(err)
	}

	json.NewDecoder(resp.Body).Decode(&tkn)
	return tkn.Token
}

// 用变量填充模板
func executeTpl(tpl string, model interface{}) string {
	tmpl, err := template.New("test").Parse(tpl)
	if err != nil {
		panic(err)
	}

	var tmplBytes bytes.Buffer
	// err = tmpl.Execute(os.Stdout, sweaters)
	err = tmpl.Execute(&tmplBytes, model)
	if err != nil {
		panic(err)
	}

	// return tmplBytes.Bytes()
	return tmplBytes.String()
}

// TimeIn returns the time in UTC if the name is "" or "UTC".
// It returns the local time if the name is "Local".
// Otherwise, the name is taken to be a location name in
// the IANA Time Zone database, such as "Africa/Lagos".
func TimeIn(t time.Time, name string) (time.Time, error) {
	loc, err := time.LoadLocation(name)
	if err == nil {
		t = t.In(loc)
	}
	return t, err
}
