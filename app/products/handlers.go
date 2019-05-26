package products

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dfang/qor-demo/models/products"
	"github.com/dfang/qor-demo/utils"
	"github.com/qor/render"
)

// Controller products controller
type Controller struct {
	View *render.Render
}

// Index products index page
func (ctrl Controller) Index(w http.ResponseWriter, req *http.Request) {
	var (
		Products []products.Product
		tx       = utils.GetDB(req)
	)

	tx.Preload("Category").Preload("ColorVariations").Find(&Products)
	ctrl.View.Execute("index", map[string]interface{}{"Products": Products}, req, w)
}

// Gender products gender page
func (ctrl Controller) Gender(w http.ResponseWriter, req *http.Request) {
	var (
		Products []products.Product
		tx       = utils.GetDB(req)
	)

	tx.Where(&products.Product{Gender: strings.Title(utils.URLParam("gender", req))}).Preload("Category").Preload("ColorVariations").Find(&Products)
	ctrl.View.Execute("gender", map[string]interface{}{"Products": Products}, req, w)
}

// Show product show page
func (ctrl Controller) Show(w http.ResponseWriter, req *http.Request) {
	var (
		product         products.Product
		variations      []products.ProductVariation
		colorVariations []products.ColorVariation
		sizeVariations  []products.SizeVariation
		colorVariation  products.ColorVariation
		codes           = strings.Split(utils.URLParam("code", req), "_")
		productCode     = codes[0]
		// colorCode      string
		tx = utils.GetDB(req)
	)

	fmt.Println(codes)

	if len(codes) > 1 {
		// colorCode = codes[1]
		fmt.Println(codes[1])
	}

	if tx.Preload("Category").Where(&products.Product{Code: productCode}).First(&product).RecordNotFound() {
		http.Redirect(w, req, "/", http.StatusFound)
	}

	// tx.Preload("Product").Preload("Color").Preload("SizeVariations.Size").Where(&products.ColorVariation{ProductID: product.ID, ColorCode: colorCode}).First(&colorVariation)

	// tx.Preload("Product").Preload("Variations").Where(&products.Product{Code: productCode})
	tx.Where(&products.Product{Code: productCode}).Find(&product)

	// SizeVariants
	tx.Where(&products.SizeVariation{ProductID: product.ID}).Preload("Size").Find(&sizeVariations)

	// ColorVariants
	tx.Where(&products.ColorVariation{ProductID: product.ID}).Preload("Color").Find(&colorVariations)

	tx.Where(&products.ProductVariation{ProductID: &product.ID}).Find(&variations)

	ctrl.View.Execute("show", map[string]interface{}{"CurrentColorVariation": colorVariation, "CurrentProduct": product, "Variations": variations, "SizeVariations": sizeVariations, "ColorVariations": colorVariations}, req, w)
}

// Category category show page
func (ctrl Controller) Category(w http.ResponseWriter, req *http.Request) {
	var (
		category products.Category
		Products []products.Product
		tx       = utils.GetDB(req)
	)

	if tx.Where("code = ?", utils.URLParam("code", req)).First(&category).RecordNotFound() {
		http.Redirect(w, req, "/", http.StatusFound)
	}

	tx.Where(&products.Product{CategoryID: category.ID}).Preload("ColorVariations").Find(&Products)
	ctrl.View.Execute("category", map[string]interface{}{"CategoryName": category.Name, "Products": Products}, req, w)
}
