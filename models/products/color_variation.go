package products

import (
  "fmt"
	"github.com/jinzhu/gorm"
	"github.com/qor/media/media_library"
  "github.com/qor/publish2"
  "github.com/dfang/qor-demo/config/db"
)

type ColorVariation struct {
	gorm.Model
	ProductID uint
	Product   Product
	ColorID   uint
	Color     Color
	ColorCode string
	Images    media_library.MediaBox
	// SizeVariations []SizeVariation
	publish2.SharedVersion
}

// ViewPath view path of color variation
func (colorVariation ColorVariation) ViewPath() string {
	defaultPath := ""
	var product Product
	if !db.DB.First(&product, "id = ?", colorVariation.ProductID).RecordNotFound() {
		defaultPath = fmt.Sprintf("/products/%s_%s", product.Code, colorVariation.ColorCode)
	}
	return defaultPath
}
