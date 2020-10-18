package orders

import (
	"github.com/jinzhu/gorm"
)

// Pricing 运费定价
// https://cl.ly/c1c76e6a4776
// eg. 超小件的电脑配送到县城10元，配送到乡下15元
// 运费根据规则引擎计算
// grool.go
type Pricing struct {
	gorm.Model

	// 针对哪个分类的 (冰洗空)
	Category string
	// 件型大小
	VolumeType string
	// 配送范围
	DeliveryArea string
	// 费用
	Fee float64
}
