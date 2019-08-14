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
	user.Meta(&admin.Meta{Name: "Gender", Config: &admin.SelectOneConfig{Collection: []string{"Male", "Female", "Unknown"}}})
	user.Meta(&admin.Meta{Name: "Birthday", Type: "date"})
	user.Meta(&admin.Meta{Name: "Role", Config: &admin.SelectOneConfig{Collection: []string{"admin", "operator", "setup_man", "delivery_man"}}})
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
			Name:  role,
			Label: strings.Title(strings.Replace(role, "_", " ", -1)),
			Group: "Filter By Role",
			Handler: func(db *gorm.DB, context *qor.Context) *gorm.DB {
				return db.Where(users.User{Role: role})
			},
		})
	}

	user.Filter(&admin.Filter{
		Name: "Role",
		Config: &admin.SelectOneConfig{
			Collection: []string{"admin", "operator", "setup_man", "delivery_man"},
		},
	})

	user.IndexAttrs("ID", "Email", "Name", "Gender", "Role", "Balance")
	user.ShowAttrs(
		&admin.Section{
			Title: "Basic Information",
			Rows: [][]string{
				{"Name"},
				{"Email", "Password"},
				{"Avatar"},
				{"Gender", "Role", "Birthday"},
				{"Confirmed"},
			},
		},
		&admin.Section{
			Title: "Credit Information",
			Rows: [][]string{
				{"Balance"},
			},
		},
		&admin.Section{
			Title: "Accepts",
			Rows: [][]string{
				{"AcceptPrivate", "AcceptLicense", "AcceptNews"},
			},
		},
		&admin.Section{
			Title: "Default Addresses",
			Rows: [][]string{
				{"DefaultBillingAddress"},
				{"DefaultShippingAddress"},
			},
		},
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
