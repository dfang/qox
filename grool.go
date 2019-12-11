package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/dfang/qor-demo/config/db"
	"github.com/dfang/qor-demo/models/orders"
	"github.com/newm4n/grool/builder"
	"github.com/newm4n/grool/context"
	"github.com/newm4n/grool/engine"
	"github.com/newm4n/grool/model"
	"github.com/newm4n/grool/pkg"
)

type Item struct {
	OrderID uint

	Quantity uint
	Price    float32

	// 单件商品的配送费 根据规则推断出来的
	DeliveryFee float32

	ProductNo string `json:"product_no"`
	OrderNo   string `json:"order_no"`
	ItemName  string `json:"product_name"`

	CustomerAddress string
}

type Helper struct {
}

func (h *Helper) Contains(s1 string, s2 string) bool {
	return strings.Contains(s1, s2)
}

func (h *Helper) NotContains(s1 string, s2 string) bool {
	return !strings.Contains(s1, s2)
}

func (h *Helper) ContainsAny(s1 string, s2 string) bool {
	v := strings.Replace(s2, "｜", "|", -1)
	if strings.Contains(v, "|") {
		s3 := strings.Split(v, "|")
		for _, s := range s3 {
			if strings.Contains(s1, s) {
				return true
			}
		}
	} else {
		return strings.Contains(s1, s2)
	}
	return false
}

func (h *Helper) ContainsAll(s1 string, s2 []string) bool {
	for _, s := range s2 {
		if !strings.Contains(s1, s) {
			return false
		}
		return true
	}
	return false
}

func (h *Helper) GetCategory(s1 string) string {
	// 冰箱（冰柜）
	// 洗衣机
	// 彩电
	// 电脑
	// 小家电
	// 厨卫
	var cats = []string{
		"空调",
		"冰箱",
		"冰柜",
		"洗衣机",
		"彩电",
		"电脑",
		"小家电",
		"厨卫",
	}

	for _, v := range cats {
		if strings.Contains(s1, v) {
			return v
		}
	}
	return ""
}

func (h *Helper) RandomRuleName() string {
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 5)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func Compile() {
	// resultB := template.ProcessFile("templates/got.tmpl", vars)
	tmpl, err := template.ParseFiles("./rules.grl.tmpl")
	if err != nil {
		panic(err)
	}
	// t := template.Must(template.New("rules").
	// 	Parse("Dot:{{.}}\n"))

	f, err := os.Create("./rules.grl")
	if err != nil {
		log.Println("create file: ", err)
		return
	}
	defer f.Close()

	fmt.Printf("模版开始:\n\n")
	// tmpl.Execute(os.Stdout, nil)
	tmpl.Execute(f, nil)
	fmt.Printf("模版结束\n\n")
}

func Evaluate() {
	fmt.Println("start evaluating ......")
	// Loading GRL on to Knowledge
	kb := model.NewKnowledgeBase()
	ruleBuilder := builder.NewRuleBuilder(kb)
	fileRes := pkg.NewFileResource("./rules.grl")
	err := ruleBuilder.BuildRuleFromResource(fileRes)
	if err != nil {
		panic(err)
	}

	// Prepare the engine
	eng := engine.NewGroolEngine()

	var ordersToEvaluate []orders.Order
	// db.DB.Order("ID DESC").Limit(100).Find(&ordersToEvaluate)
	// db.DB.Order("random()").Limit(100).Find(&ordersToEvaluate)
	// db.DB.LogMode(true)
	db.DB.Order("random()").Find(&ordersToEvaluate)
	// db.DB.Order("random()").Limit(10).Find(&ordersToEvaluate)

	// db.DB.Where("item_name LIKE $1", "%电视%").Find(&ordersToEvaluate)

	// fmt.Println(len(ordersToEvaluate))

	for _, o := range ordersToEvaluate {
		EvaluateOrder(o, kb, eng)
	}

	// fmt.Println("-----------------")
	fmt.Println("执行完毕!")
}

func EvaluateOrder(order orders.Order, knowledgeBase *model.KnowledgeBase, engine *engine.Grool) {
	// Preparing Facts
	// add facts to data context
	// dataContext := context.NewDataContext()
	// var order orders.Order
	var orderItems []orders.OrderItem
	// db.DB.Model(orders.Order{}).First(&order)
	// fmt.Println(order.ID)
	db.DB.Model(&order).Related(&orderItems)
	var items []orders.OrderItem

	for _, oi := range orderItems {
		oi.Order = order
		// fmt.Println(oi.Order)
		items = append(items, oi)
	}

	for _, v := range items {
		dctx := context.NewDataContext()
		err := dctx.Add("Pogo", &Helper{})
		if err != nil {
			panic(err)
		}

		err = dctx.Add("Item", &v)
		if err != nil {
			panic(err)
		}
		// fmt.Println("......... 应用规则引擎 ........ ")
		err = engine.Execute(dctx, knowledgeBase)
		if err != nil {
			panic(err)
		}

		fmt.Println("输入如下: ")
		fmt.Println("客户地址: " + v.Order.CustomerAddress)
		fmt.Println("商品名: " + v.ItemName)
		fmt.Println("-------------------------------")
		fmt.Println("判断结果如下:")
		if v.Range == 1 {
			fmt.Println("配送范围: ", "县城")
		} else {
			fmt.Println("配送范围: ", "乡村")
		}
		fmt.Println("分类: ", v.Category)
		fmt.Println("大小: ", v.Dimension)
		fmt.Println("推算出来的配送价格为: ", v.DeliveryFee)

		db.DB.Save(&v)
	}
	// Executing A Knowledge On Facts and get result
}
