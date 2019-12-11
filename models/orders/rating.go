package orders

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/audited"
)

// Rating 记录好评差评模块
type Rating struct {
	gorm.Model

	OrderID uint
	Order   Order

	// 好评还是差评
	Type string

	// 具体原因
	Reason string

	// 奖励或罚款金额
	Amount int

	Remark string

	audited.AuditedModel
}
