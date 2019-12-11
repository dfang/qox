package orders

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/sorting"
)

type Rule struct {
	gorm.Model
	sorting.Sorting

	Name string

	// 针对哪个分类的
	Category string
	// 作用 (判断大小还是判断分类的、范围、定价)
	Effect string

	Description string
	Priority    int

	Conditions []Condition
	Executions []Execution
}

type Condition struct {
	gorm.Model

	RuleID uint
	Rule   Rule

	Name string
	// 运算符
	Operator string
	Value    string
}

type Execution struct {
	gorm.Model

	RuleID uint
	Rule   Rule

	Name  string
	Value string
}

// Rule
// id      name    description
// 1
// 2
// 3

// Condition
// id  rule_id  content
// 1     1

// Action

// id   name
// 1    设置类别
// 2    设置品类
// 3    设置大小
// 4    设置配送

// Rule_Actions

// id rule_id action_id remark
// 1   1       1
