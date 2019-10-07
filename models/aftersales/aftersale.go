package aftersales

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/audited"
	"github.com/qor/transition"
	"time"
	"strings"

	"github.com/dfang/qor-demo/models/users"
)

// AfterSale 售后管理
type AfterSale struct {
	gorm.Model

	// 客户信息, 从京东后台导入或者扫描枪扫入的
	CustomerName    string
	CustomerPhone   string
	CustomerAddress string

	// -- ORDER_TYPE starts with Q 退货的取件单
	ServiceType    string
	ServiceContent string

	// 预约安装时间
	ReserverdServiceTime string

	// 品牌
	Source  string

	Fee float32

	// 备注
	Remark string

	UserID *uint
	User   users.User

	audited.AuditedModel

	transition.Transition
}

func (item *AfterSale) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("State", "created")
  return nil
}

// AfterSale
// 姓名
// 地址
// 电话
// 服务类型 安装 清洗 维修
// 服务内容
// 备注

const (
	CREATED_FROM = "created_at_from:query"
	CREATED_TO = "created_at_to:query"
) 

func RegisterCallbacks(db *gorm.DB) {
	if db.Callback().Query().Get(CREATED_FROM) == nil {
		db.Callback().Query().Before("gorm:query").Register(CREATED_FROM, queryCallback)
	}

	// if db.Callback().RowQuery().Get(CREATED_TO) == nil {
	// 	db.Callback().RowQuery().Before("gorm:row_query").Register(CREATED_TO, queryCallback)
	// }
}

func queryCallback(scope *gorm.Scope) {
	var			conditions         []string
	var			conditionValues    []interface{}
	// var			scheduledStartTime, scheduledEndTime, scheduledCurrentTime *time.Time
	var			scheduledStartTime, scheduledEndTime *time.Time

	if v, ok := scope.Get(CREATED_FROM); ok {
		if t, ok := v.(*time.Time); ok {
			scheduledStartTime = t
		} else if t, ok := v.(time.Time); ok {
			scheduledStartTime = &t
		}

		if scheduledStartTime != nil {
			conditions = append(conditions, "created_at >= ?")
			conditionValues = append(conditionValues, scheduledStartTime)
		}
	}

	if v, ok := scope.Get(CREATED_TO); ok {
		if t, ok := v.(*time.Time); ok {
			scheduledEndTime = t
		} else if t, ok := v.(time.Time); ok {
			scheduledEndTime = &t
		}

		if scheduledEndTime != nil {
			conditions = append(conditions, "created_at <= ?")
			conditionValues = append(conditionValues, scheduledEndTime)
		}
	}

	scope.Search.Where(strings.Join(conditions, " AND "), conditionValues...)
}

// func init() {
// 	RegisterCallbacks	
// }


