package products

import (
	"fmt"

	"github.com/dfang/qor-demo/config/db"
	"github.com/jinzhu/gorm"
	"github.com/qor/publish2"
)

type SizeVariation struct {
	gorm.Model
	ProductID uint
	Product   Product

	// ColorVariationID  uint
	// ColorVariation    ColorVariation
	SizeID uint
	Size   Size
	// AvailableQuantity uint
	publish2.SharedVersion
}

// ViewPath view path of color variation
func (sizeVariation SizeVariation) ViewPath() string {
	defaultPath := ""
	var product Product
	if !db.DB.First(&product, "id = ?", sizeVariation.ProductID).RecordNotFound() {
		defaultPath = fmt.Sprintf("/products/%s", product.Code)
	}
	return defaultPath
}
