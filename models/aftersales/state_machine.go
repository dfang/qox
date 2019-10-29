package aftersales

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/dfang/qor-demo/models/users"
	"github.com/jinzhu/gorm"
	"github.com/qor/transition"
)

var (
	// OrderStateMachine order's state machine
	OrderStateMachine = transition.New(&Aftersale{})

	// SettlementStateMachine 结算状态机
	SettlementStateMachine = transition.New(&Settlement{})
)

var (
	// DraftState draft state
	DraftState = "created"
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
	OrderStateMachine.Event("inquire").To("inquired").From("created").After(func(value interface{}, tx *gorm.DB) (err error) {
		// order := value.(*Aftersale)
		// tx.Model(order).Association("OrderItems").Find(&order.OrderItems)
		// for _, item := range order.OrderItems {
		// }
		return nil
	})

	// 指派师傅
	OrderStateMachine.Event("schedule").To("scheduled").From("inquired").After(func(value interface{}, tx *gorm.DB) (err error) {
		// 推送微信模版消息到师傅微信
		item := value.(*Aftersale)
		u := users.User{}
		tx.Where("id = ?", item.UserID).First(&u)
		mobilePhone := u.MobilePhone
		wp := users.WechatProfile{}
		tx.Where("mobile_phone = ?", mobilePhone).First(&wp)

		fmt.Println("user_id is > ", item.UserID)
		fmt.Println("mobile_phone is > ", mobilePhone)
		fmt.Println("openid is > ", wp.Openid)

		m := ModelForTpl1{
			OpenID: wp.Openid,
			ID:     strconv.FormatUint(uint64(item.ID), 10),
			URL:    "http://mp.xsjd123.com/",
			Date:   time.Now().Format("2006-01-02 15:04"),
		}

		b := executeTpl1(tpl1, m)
		fmt.Println("模板消息插入变量后是: ", string(b))

		sendTemplateMsg(b)

		return nil
	})

	OrderStateMachine.Event("expire").To("overdue").From("scheduled").After(func(value interface{}, tx *gorm.DB) (err error) {
		return nil
	})

	// 根据师傅上传的照片 审核服务是否完成
	OrderStateMachine.Event("audit").To("audited").From("processed").After(func(value interface{}, tx *gorm.DB) (err error) {
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

		return nil
	})
}

// 模板消息
func sendTemplateMsg(contents []byte) {
	// https://mp.weixin.qq.com/advanced/tmplmsg?action=faq&token=1545649248&lang=zh_CN
	tkn := getToken()
	url := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + tkn
	fmt.Println("URL:>", url)

	// var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	// req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(contents))
	fmt.Println("Request Body:> ", string(contents))

	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
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
func executeTpl1(tpl string, model ModelForTpl1) []byte {
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

	return tmplBytes.Bytes()
}
