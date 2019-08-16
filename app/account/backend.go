package account

import (
	"strings"

	"github.com/dfang/qor-demo/models/users"
	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/qor/qor/resource"
	qorutils "github.com/qor/qor/utils"
	"github.com/qor/validations"
	"golang.org/x/crypto/bcrypt"
)

// ConfigureAdmin configure admin interface
func (App) ConfigureAdmin(Admin *admin.Admin) {
	Admin.AddMenu(&admin.Menu{Name: "User Management", Priority: 3})
	user := Admin.AddResource(&users.User{}, &admin.Config{Menu: []string{"User Management"}})

	user.SearchAttrs("name", "mobile_phone")

	user.Meta(&admin.Meta{
		Name:   "Gender",
		Type:   "string",
		Config: &admin.SelectOneConfig{Collection: []string{"男", "女"}},
		FormattedValuer: func(record interface{}, context *qor.Context) (value interface{}) {
			user := record.(*users.User)
			switch strings.ToLower(user.Gender) {
			case "male":
				return "男"
			case "female":
				return "女"
			default:
				return "男"
			}
		},
		Valuer: func(record interface{}, context *qor.Context) (value interface{}) {
			user := record.(*users.User)
			switch strings.ToLower(user.Gender) {
			case "male":
				return "男"
			case "female":
				return "女"
			default:
				return "男"
			}
		},
		Setter: func(record interface{}, metaValue *resource.MetaValue, context *qor.Context) {
			// fmt.Println(metaValue.Value)
			mv := qorutils.ToString(metaValue.Value)
			var m string
			switch mv {
			case "男":
				m = "male"
			case "女":
				m = "female"
			default:
				m = "male"
			}
			record.(*users.User).Gender = m
		},
	})

	user.Meta(&admin.Meta{Name: "Birthday", Type: "date"})

	// user.Meta(&admin.Meta{Name: "Role", Config: &admin.SelectOneConfig{Collection: []string{"admin", "operator", "setup_man", "delivery_man"}}})
	user.Meta(&admin.Meta{
		Name:   "Role",
		Type:   "string",
		Config: &admin.SelectOneConfig{Collection: []string{"管理员", "调度员", "安装师傅", "配送师傅"}},
		Valuer: func(record interface{}, context *qor.Context) (value interface{}) {
			user := record.(*users.User)
			switch user.Role {
			case "operator":
				return "调度员"
			case "delivery_man":
				return "配送师傅"
			case "setup_man":
				return "安装师傅"
			case "Admin":
				return "管理员"
			default:
				return "调度员"
			}
		},
		Setter: func(record interface{}, metaValue *resource.MetaValue, context *qor.Context) {
			// fmt.Println(metaValue.Value)
			mv := qorutils.ToString(metaValue.Value)
			var m string
			switch mv {
			case "调度员":
				m = "operator"
			case "配送师傅":
				m = "delivery_man"
			case "安装师傅":
				m = "setup_man"
			case "管理员":
				m = "admin"
			default:
				m = "operator"
			}
			record.(*users.User).Role = m
		},
		FormattedValuer: func(record interface{}, _ *qor.Context) (result interface{}) {
			user := record.(*users.User)
			// return user.Role
			switch user.Role {
			case "operator":
				return "调度员"
			case "delivery_man":
				return "配送师傅"
			case "setup_man":
				return "安装师傅"
			case "Admin":
				return "管理员"
			default:
				return "调度员"
			}
		},
	})

	user.Meta(&admin.Meta{Name: "Password",
		Type:   "password",
		Valuer: func(interface{}, *qor.Context) interface{} { return "" },
		Setter: func(resource interface{}, metaValue *resource.MetaValue, context *qor.Context) {
			if newPassword := qorutils.ToString(metaValue.Value); newPassword != "" {
				bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
				if err != nil {
					context.DB.AddError(validations.NewError(user, "Password", "Can't encrpt password"))
					return
				}
				u := resource.(*users.User)
				u.Password = string(bcryptPassword)
			}
		},
	})
	user.Meta(&admin.Meta{Name: "Confirmed", Valuer: func(user interface{}, ctx *qor.Context) interface{} {
		if user.(*users.User).ID == 0 {
			return true
		}
		return user.(*users.User).Confirmed
	}})
	user.Meta(&admin.Meta{Name: "DefaultBillingAddress", Config: &admin.SelectOneConfig{Collection: userAddressesCollection}})
	user.Meta(&admin.Meta{Name: "DefaultShippingAddress", Config: &admin.SelectOneConfig{Collection: userAddressesCollection}})

	for _, role := range []string{"admin", "operator", "setup_man", "delivery_man"} {
		var role = role
		user.Scope(&admin.Scope{
			Name: role,
			// Label: strings.Title(strings.Replace(role, "_", " ", -1)),
			Label: users.V2T(users.ROLES_VALUES, users.ROLES_TEXTS, role),
			Group: "角色",
			Handler: func(db *gorm.DB, context *qor.Context) *gorm.DB {
				return db.Where(users.User{Role: role})
			},
		})
	}

	for _, item := range []string{"male", "female"} {
		var gender = item
		user.Scope(&admin.Scope{
			Name: gender,
			// Label: strings.Title(strings.Replace(role, "_", " ", -1)),
			Label: users.V2T(users.GENDERS_VALUES, users.GENDERS_TEXTS, gender),
			Group: "性别",
			Handler: func(db *gorm.DB, context *qor.Context) *gorm.DB {
				return db.Where(users.User{Gender: gender})
			},
		})
	}

	// user.Filter(&admin.Filter{
	// 	Name: "Role",
	// 	Config: &admin.SelectOneConfig{
	// 		Collection: []string{"admin", "operator", "setup_man", "delivery_man"},
	// 	},
	// })

	user.IndexAttrs("ID", "Name", "MobilePhone", "Gender", "Role")
	user.NewAttrs("ID", "Name", "Gender", "Role", "MobilePhone", "IdentityCardNum")
	user.ShowAttrs(
		&admin.Section{
			Title: "Basic Information",
			Rows: [][]string{
				{"Name", "Gender"},
				{"MobilePhone"},
				{"Role"},
			},
		},
		// &admin.Section{
		// 	Title: "Default Addresses",
		// 	Rows: [][]string{
		// 		{"DefaultBillingAddress"},
		// 		{"MobilePhone"},
		// 		{"Role"},
		// 	},
		// },
		"Addresses",
	)
	user.EditAttrs(user.ShowAttrs())

	// user.UseTheme("grid")

	// Add deliveryMen submenu
	deliveryMan := Admin.AddResource(&users.User{}, &admin.Config{Name: "Delivery Men", Menu: []string{"User Management"}})
	deliveryMan.Scope(&admin.Scope{
		Default: true,
		Handler: func(db *gorm.DB, context *qor.Context) *gorm.DB {
			return db.Where("role = ?", "delivery_man")
		},
	})

	// Add  submenu
	setupMan := Admin.AddResource(&users.User{}, &admin.Config{Name: "Setup Men", Menu: []string{"User Management"}})
	setupMan.Scope(&admin.Scope{
		Default: true,
		Handler: func(db *gorm.DB, context *qor.Context) *gorm.DB {
			return db.Where("role = ?", "setup_man")
		},
	})

	// Add  submenu
	operator := Admin.AddResource(&users.User{}, &admin.Config{Name: "Operator", Menu: []string{"User Management"}})
	operator.Scope(&admin.Scope{
		Default: true,
		Handler: func(db *gorm.DB, context *qor.Context) *gorm.DB {
			return db.Where("role = ?", "operator")
		},
	})

}
