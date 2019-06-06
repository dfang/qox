package orders

import (
	"fmt"

	"github.com/dfang/qor-demo/config/db"
	"github.com/dfang/qor-demo/models/products"
	"github.com/jinzhu/gorm"
	"github.com/qor/transition"
)

type OrderItem struct {
	gorm.Model
	OrderID          uint
	SizeVariationID  uint `cartitem:"SizeVariationID"`
	SizeVariation    *products.SizeVariation
	ColorVariationID uint `cartitem:"ColorVariationID"`
	ColorVariation   *products.ColorVariation
	Quantity         uint `cartitem:"Quantity"`
	Price            float32
	DiscountRate     uint
	transition.Transition

	ProductNo   string `json:"product_no"`
	OrderNo     string `json:"order_no"`
	ProductName string `json:"product_name"`
	Install     string `json:"install"`
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
