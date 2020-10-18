package orders

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/dfang/qor-demo/config/db"
	"github.com/dfang/qor-demo/models/products"
	"github.com/gocraft/work"
	"github.com/jinzhu/gorm"
	"github.com/qor/transition"
)

type OrderItem struct {
	gorm.Model
	OrderID          uint `json: "order_id`
	Order            Order
	SizeVariationID  uint `cartitem:"SizeVariationID"`
	SizeVariation    *products.SizeVariation
	ColorVariationID uint `cartitem:"ColorVariationID"`
	ColorVariation   *products.ColorVariation
	Quantity         uint    `cartitem:"Quantity"`
	Price            float32 `json:"price"`
	DiscountRate     uint

	ProductNo string `json:"product_no"`
	OrderNo   string `json:"order_no"`
	ItemName  string `json:"product_name"`
	Install   string `json:"install"`

	Range     int    `json:"range"`
	Category  string `json:"category"`
	Dimension string `json:"dimension"`

	// 单件商品的配送费 根据规则推断出来的
	DeliveryFee float64 `json:"delivery_fee"`

	transition.Transition
}

// IsCart order item's state is cart
func (item OrderItem) IsCart() bool {
	return item.State == DraftState || item.State == ""
}

func (item *OrderItem) loadSizeVariation() {
	// fmt.Println("loadSizeVariation")
	// fmt.Println(item.SizeVariation)
	// if item.SizeVariation == nil {
	// 	item.SizeVariation = &products.SizeVariation{}
	// 	// db.DB.Model(item).Preload("Size").Preload("ColorVariation.Product").Preload("ColorVariation.Color").Association("SizeVariation").Find(item.SizeVariation)
	// 	fmt.Println("load sizeVariation")
	// 	db.DB.Model(item).Preload("Size").Association("SizeVariation.Product").Find(item.SizeVariation)
	// }
}

func (item *OrderItem) Size() products.Size {
	var sizeVariation products.SizeVariation
	var size products.Size
	fmt.Printf("%+v\n", item)
	db.DB.Model(item).Where("id = ?", item.SizeVariationID).Find(&sizeVariation)
	fmt.Printf("%+v\n", sizeVariation)
	db.DB.Where("id = ?", sizeVariation.SizeID).Find(&size)
	fmt.Println(size.Name)
	fmt.Println(size.Code)
	return size
}

func (item *OrderItem) Color() products.Color {
	// db.DB.Model(item).Preload("Color").Association("ColorVariation").Find(item.ColorVariation)
	var colorVariation products.ColorVariation
	var color products.Color
	db.DB.Model(item).Where("id = ?", item.ColorVariationID).Find(&colorVariation)
	db.DB.Where("id = ?", colorVariation.ColorID).Find(&color)
	return color
}

// // ProductImageURL get product image
// func (item *OrderItem) ProductImageURL() string {
// 	item.loadSizeVariation()
// 	// return item.SizeVariation.ColorVariation.MainImageURL()
// }

// SellingPrice order item's selling price
func (item *OrderItem) SellingPrice() float32 {
	// if item.IsCart() {
	// 	item.loadSizeVariation()
	//  return item.SizeVariation.ColorVariation.Product.Price
	// }
	// return item.Price
	return 0
}

// ProductName order item's color name
func (item *OrderItem) ProductName() string {
	// item.loadSizeVariation()
	// fmt.Println("######")
	// fmt.Println(item.SizeVariation)
	// fmt.Println(item.SizeVariation.ID)
	// fmt.Println(item.SizeVariation.Product)
	// fmt.Println(item.SizeVariation.Product.Name)
	var product products.Product
	db.DB.Model(item).Find(&product)
	return product.Name
}

// // ColorName order item's color name
// func (item *OrderItem) ColorName() string {
// 	item.loadSizeVariation()
// 	// return item.SizeVariation.ColorVariation.Color.Name
// }

// // SizeName order item's size name
// func (item *OrderItem) SizeName() string {
// 	item.loadSizeVariation()
// 	return item.SizeVariation.Size.Name
// }

// ProductPath order item's product name
func (item *OrderItem) ProductPath() string {
	// item.loadSizeVariation()
	return item.SizeVariation.ViewPath()
}

// Amount order item's amount
func (item OrderItem) Amount() float32 {
	// amount := item.SellingPrice() * float32(item.Quantity)
	// if item.DiscountRate > 0 && item.DiscountRate <= 100 {
	// 	amount = amount * float32(100-item.DiscountRate) / 100
	// }
	// return amount
	return 0
}

func (item *OrderItem) SetRange(i int64) {
	item.Range = int(i)
}

func (item *OrderItem) SetCategory(c string) {
	item.Category = c
}

func (item *OrderItem) SetDimension(d string) {
	item.Dimension = d
}

func (item *OrderItem) SetDeliveryFee(f float64) {
	item.DeliveryFee = f
}

// GetLengthUnit 获取长度单位（x英寸）中的数字
// 电视机
func (item *OrderItem) GetLengthUnit() int64 {
	// match, _ := regexp.MatchString("\d+英寸", item.ItemName)
	re := regexp.MustCompile(`\d+英寸`)
	re2 := regexp.MustCompile(`\d+`)
	// fmt.Printf("%q\n", re.FindString(item.ItemName))
	if re.FindString(item.ItemName) != "" {
		d := re2.FindString(re.FindString(item.ItemName))
		n, err := strconv.ParseInt(d, 10, 64)
		if err != nil {
			fmt.Printf("%d of type %T", n, n)
		}
		return n
	}
	return 0
}

// GetWeightUnit 获取重量单位（x公斤/kg）中的数字
// 洗衣机
func (item *OrderItem) GetWeightUnit() int64 {
	// match, _ := regexp.MatchString("\d+英寸", item.ItemName)
	re := regexp.MustCompile(`\d+(?:[(KG)|(公斤)|(kg)]+)`)
	re2 := regexp.MustCompile(`\d+`)
	// fmt.Printf("%q\n", re.FindString(item.ItemName))
	if re.FindString(item.ItemName) != "" {
		d := re2.FindString(re.FindString(item.ItemName))
		n, err := strconv.ParseInt(d, 10, 64)
		if err != nil {
			fmt.Printf("%d of type %T", n, n)
		}
		return n
	}
	return 0
}

// GetVolumeUnit 获取重量单位（x升/L）中的数字
// 冰箱 冰柜 热水器
func (item *OrderItem) GetVolumeUnit() int64 {
	re := regexp.MustCompile(`\d+(?:[(L)|(升)|(l)]+)`)
	re2 := regexp.MustCompile(`\d+`)
	// fmt.Printf("%q\n", re.FindString(item.ItemName))
	if re.FindString(item.ItemName) != "" {
		d := re2.FindString(re.FindString(item.ItemName))
		n, err := strconv.ParseInt(d, 10, 64)
		if err != nil {
			fmt.Printf("%d of type %T", n, n)
		}
		return n
	}
	return 0
}

// GetPowerUnit 获取功率单位（x匹）中的数字
// 空调
func (item *OrderItem) GetPowerUnit() int64 {
	re := regexp.MustCompile(`\d+(?:[(L)|(升)|(l)]+)`)
	re2 := regexp.MustCompile(`\d+`)
	// fmt.Printf("%q\n", re.FindString(item.ItemName))
	if re.FindString(item.ItemName) != "" {
		d := re2.FindString(re.FindString(item.ItemName))
		n, err := strconv.ParseInt(d, 10, 64)
		if err != nil {
			fmt.Printf("%d of type %T", n, n)
		}
		return n
	}
	return 0
}

// GetUnit 获取单位
func (item *OrderItem) GetUnit() int64 {

	if strings.Contains(item.ItemName, "电视") {
		return item.GetLengthUnit()
	}

	return 0
}

// IsService 是否是增值服务
func (item *OrderItem) IsService() bool {
	s1 := item.ItemName
	s2 := []string{"延保", "服务", "安心享", "延误补贴", "礼包", "只换不修"}
	for _, s := range s2 {
		if strings.Contains(s1, s) {
			return true
		}
	}
	return false
}

// AfterCreate 初始状态
func (item *OrderItem) AfterCreate(scope *gorm.Scope) error {
	if strings.Contains(item.OrderNo, "Q") {
		scope.SetColumn("reserved_pickup_time", item.Order.ReservedDeliveryTime)
	}
	var enqueuer = work.NewEnqueuer("qor", db.RedisPool)
	enqueuer.Enqueue("update_order_items", work.Q{})
	// enqueuer.EnqueueIn("create_after_sale", 60, work.Q{"order_no": item.OrderNo})

	// a := aftersales.Aftersale{}
	// a.CustomerAddress = item.Order.CustomerAddress
	// a.CustomerName = item.Order.CustomerName
	// a.CustomerPhone = item.Order.CustomerPhone
	// a.ReservedServiceTime = item.Order.ReservedSetupTime
	// a.Source = "JD"
	// a.Remark = item.Order.OrderNo

	return nil
}
