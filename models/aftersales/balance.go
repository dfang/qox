package aftersales

import (
	"fmt"

	"github.com/dfang/qor-demo/models/users"
	"github.com/jinzhu/gorm"
)

// Balance 用户余额统计
type Balance struct {
	gorm.Model

	UserID uint
	User   users.User

	// 总冻结金额
	FrozenAmount float32

	// 总可提现金额
	FreeAmount float32

	// 奖励金
	AwardAmount float32

	// 罚款
	FineAmount float32

	// 历史总收入
	TotalAmount float32

	// 历史总提现金额
	WithdrawAmount float32

	// transition.Transition
}

func (item *Balance) BeforeSave(scope *gorm.Scope) error {
	if item.UserID <= 0 {
		return fmt.Errorf("balance user_id 不能为空")
	}
	return nil
}
