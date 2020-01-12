package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dfang/qor-demo/config/db"

	"github.com/contribsys/faktory/client"
	"github.com/dfang/qor-demo/models/orders"

	worker "github.com/contribsys/faktory_worker_go"
)

// just run `go startFaktoryWorker()` in main.go
func startFaktoryWorker() {
	// - FAKTORY_URL=tcp://:admin@faktory:7419
	if os.Getenv("FAKTORY_URL") == "" {
		panic("Please set FAKTORY_URL")
	}

	mgr := worker.NewManager()
	mgr.Use(func(perform worker.Handler) worker.Handler {
		return func(ctx worker.Context, job *client.Job) error {
			log.Printf("Starting work on job %s of type %s with custom %v\n", ctx.Jid(), ctx.JobType(), job.Custom)
			err := perform(ctx, job)
			log.Printf("Finished work on job %s with error %v\n", ctx.Jid(), err)
			return err
		}
	})

	// register job types and the function to execute them
	mgr.Register("upsert_order", UpsertOrder)

	// use up to N goroutines to execute jobs
	mgr.Concurrency = 10

	// pull jobs from these queues, in this order of precedence
	mgr.ProcessStrictPriorityQueues("critical", "upsert_order", "default", "bulk")

	// alternatively you can use weights to avoid starvation
	//mgr.ProcessWeightedPriorityQueues(map[string]int{"critical":3, "default":2, "bulk":1})

	// Start processing jobs, this method does not return
	mgr.Run()
}

func UpsertOrder(ctx worker.Context, args ...interface{}) error {
	// log.Printf("Working on job %s\n", ctx.Jid())
	// log.Printf("Context %v\n", ctx)
	// log.Printf("Args %v\n", args)
	var s string
	switch v := args[0].(type) {
	case string:
		fmt.Println("args[0] is", v) // here v has type string
		s = args[0].(string)
	default:
		return errors.New("参数正确")
	}

	// payload := strings.NewReader(s)
	// url := os.Getenv("DESTINATION_URL")
	// // 第三个参数不是string, 以下几种类型都支持
	// // https://github.com/hashicorp/go-retryablehttp/blob/master/client.go#L122
	// resp, err := retryablehttp.Post(url, "application/json", payload)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(resp.Status)
	// fmt.Println(resp.StatusCode)

	var obj OrderPayload
	if err := json.Unmarshal([]byte(s), &obj); err != nil {
		panic(err)
	}

	filename := fmt.Sprintf("%s/%s/%s.json", "./public/data", obj.Order.CreatedAt.Format("2006/01/02"), obj.Order.OrderNo)
	content, _ := json.Marshal(obj)
	err := saveFile(filename, content)
	if err != nil {
		return errors.New("save json to file failed")
	}

	normalized := normalizeOrder(obj)
	receivables, _ := normalized.Order.Receivables.Float64()
	o := orders.Order{
		OrderNo:              normalized.Order.OrderNo,
		CustomerAddress:      normalized.Order.CustomerAddress,
		CustomerName:         normalized.Order.CustomerName,
		CustomerPhone:        normalized.Order.CustomerPhone,
		ReservedDeliveryTime: normalized.Order.ReservedDeliveryTime,
		ReservedSetupTime:    normalized.Order.ReservedSetupTime,
		IsDeliveryAndSetup:   normalized.Order.IsDeliveryAndSetup,
		// CreatedAt: obj.Order.CreatedAt,
		// UpdatedAt: obj.Order.UpdatedAt,
		Receivables: float32(receivables),
	}
	var orderItems []orders.OrderItem
	for _, item := range obj.Order.OrderItems {
		q, _ := item.Quantity.Int64()
		oi := orders.OrderItem{
			OrderNo:   item.OrderNo,
			ProductNo: item.ProductNo,
			ItemName:  item.ProductName,
			Quantity:  uint(q),
			// Volume:    item.Volume,
			// Weight:    item.Weight,
			Install: item.Install,
		}
		orderItems = append(orderItems, oi)
	}

	var order orders.Order
	db.DB.Where("order_no = ?", obj.Order.OrderNo).First(&order)
	if order.OrderNo == "" {
		o.OrderItems = orderItems
	}

	db.DB.Where("order_no = ?", obj.Order.OrderNo).Assign(o).FirstOrInit(&o)
	if err := db.DB.Save(&o).Error; err != nil {
		log.Println(err)
		// return errors.New("uncomsued")
		return err
	}

	fmt.Printf("%+v\n", o)
	return nil
}

// 清洗数据
func normalizeOrder(payload OrderPayload) OrderPayload {
	item := payload.Order
	// 预处理
	var phone string
	phones := strings.Split(strings.TrimSpace(item.CustomerPhone), "/")
	if len(phones) > 1 && phones[0] == phones[1] {
		phone = phones[0]
	} else {
		phone = item.CustomerPhone
	}
	item.OrderNo = strings.TrimSpace(item.OrderNo)
	// 去掉首尾空格
	item.CustomerAddress = strings.TrimSpace(item.CustomerAddress)
	// 去掉"江西 九江市" 中间的空格
	item.CustomerAddress = strings.Replace(item.CustomerAddress, " ", "", -1)
	// 去掉 江西省,九江市,修水县, 中的逗号
	item.CustomerAddress = strings.Replace(item.CustomerAddress, ",", "", -1)
	item.CustomerAddress = strings.Replace(item.CustomerAddress, "，", "", -1)
	// 江西省.九江市.修水县
	item.CustomerAddress = strings.Replace(item.CustomerAddress, ".", "", -1)
	// 去掉客户姓名的首尾空格
	item.CustomerName = strings.TrimSpace(item.CustomerName)
	// 15618903080/15618903080 这种格式的电话号码，如果/左右两边一样，取一个
	item.CustomerPhone = phone
	// if Receivables not convertable (eg. " "), set it to 0,
	_, err := strconv.Atoi(item.Receivables.String())
	if err != nil {
		item.Receivables = "0"
	}

	payload.Order = item
	return payload
}

func saveFile(filename string, content []byte) error {
	dir := filepath.Dir(filename)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}

	err := ioutil.WriteFile(filename, content, 0644)
	if err != nil {
		fmt.Println(err)
		return errors.New("saveFile failed")
	}

	return nil
}

// OrderPayload 函数用来接收并发送消息到Pubsub的结构体
type OrderPayload struct {
	UUID  string `json:"uuid"`
	Order struct {
		OrderNo              string      `json:"order_no"`
		CustomerAddress      string      `json:"customer_address"`
		CustomerName         string      `json:"customer_name"`
		CustomerPhone        string      `json:"customer_phone"`
		ReservedDeliveryTime string      `json:"reserved_delivery_time"`
		ReservedSetupTime    string      `json:"reserved_setup_time"`
		IsDeliveryAndSetup   string      `json:"is_delivery_and_setup"`
		CreatedAt            time.Time   `json:"created_at"`
		UpdatedAt            time.Time   `json:"updated_at"`
		Receivables          json.Number `json:"receivables"`
		OrderItems           []struct {
			OrderNo     string      `json:"order_no"`
			ProductNo   string      `json:"product_no"`
			ProductName string      `json:"product_name"`
			Quantity    json.Number `json:"quantity"`
			Install     string      `json:"install"`
			Volume      json.Number `json:"volume"`
			Weight      json.Number `json:"weight"`
			CreatedAt   time.Time   `json:"created_at"`
			UpdatedAt   time.Time   `json:"updated_at"`
		} `json:"order_items"`
	} `json:"order"`
}
