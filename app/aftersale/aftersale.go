package aftersale

// "net/http"
import (
	"fmt"
	"github.com/jinzhu/gorm"

	"github.com/dfang/qor-demo/models/aftersales"
	"github.com/dfang/qor-demo/models/users"
	"github.com/dfang/qor-demo/models/settings"
	"github.com/qor/admin"
	"github.com/qor/application"
	"github.com/qor/qor"
	"github.com/dfang/qor-demo/config/db"
)

// New new home app
func New(config *Config) *App {
	return &App{Config: config}
}

// NewWithDefault new home app
func NewWithDefault() *App {
	return &App{Config: &Config{}}
}

// App home app
type App struct {
	Config *Config
}

// Config home config struct
type Config struct {
}

var brands []settings.Brand
var service_types []settings.ServiceType
 
// ConfigureApplication configure application
func (app App) ConfigureApplication(application *application.Application) {
	// 售后后台
	app.ConfigureAdmin(application.Admin)
}

// ConfigureAdmin configure admin interface
func (App) ConfigureAdmin(Admin *admin.Admin) {
	Admin.AddMenu(&admin.Menu{Name: "Aftersale Management", Priority: 6})

	db.DB.Select("name, id").Find(&brands)
	db.DB.Select("name, id").Find(&service_types)

	// Add Aftersale
	aftersale := Admin.AddResource(&aftersales.AfterSale{}, &admin.Config{Menu: []string{"Aftersale Management"}, Priority: 1})
	manufacturer := Admin.AddResource(&aftersales.Manufacturer{}, &admin.Config{Menu: []string{"Aftersale Management"}, Priority: 4})
	Admin.AddResource(&settings.Brand{}, &admin.Config{Name: "Brand", Menu: []string{"Aftersale Management"}, Priority: 3})
	Admin.AddResource(&settings.ServiceType{}, &admin.Config{Name: "ServiceType", Menu: []string{"Aftersale Management"}, Priority: 2})

	aftersale.Meta(&admin.Meta{
		Name:       "ServiceType",
		Type:       "select_one",
		// Collection: []string{"安装", "维修", "清洗"},
		Collection: func(value interface{}, context *qor.Context) (options [][]string) {
			for _, m := range service_types {
				idStr := fmt.Sprintf("%s", m.Name)
				var option = []string{idStr, m.Name}
				options = append(options, option)
			}
			return options
		},
	})

	aftersale.Meta(&admin.Meta{
		Name:       "Source",
		Type:       "select_one",
		Collection: func(value interface{}, context *qor.Context) (options [][]string) {
			for _, m := range brands {
				idStr := fmt.Sprintf("%s", m.Name)
				var option = []string{idStr, m.Name}
				options = append(options, option)
			}
			return options
		},
	})

	// https://doc.getqor.com/admin/metas/select-one.html
	// Generate options by data from the database
	aftersale.Meta(&admin.Meta{
		Name:  "UserID",
		Type:  "select_one",
		Label: "分配",
		Config: &admin.SelectOneConfig{
			Collection: func(_ interface{}, context *admin.Context) (options [][]string) {
				var users []users.User
				context.GetDB().Where("role = ?", "setup_man").Find(&users)
				for _, n := range users {
					idStr := fmt.Sprintf("%d", n.ID)
					var option = []string{idStr, n.Name}
					options = append(options, option)
				}
				return options
			},
			AllowBlank: true,
		}})

	manufacturer.Action(&admin.Action{
		Name: "打开厂家后台网站",
		URL: func(record interface{}, context *admin.Context) string {
			if item, ok := record.(*aftersales.Manufacturer); ok {
				return fmt.Sprintf("%s", item.URL)
			}
			return "#"
		},
		URLOpenType: "_blank",
		Modes:       []string{"menu_item", "edit", "show"},
	})
	

	configureMetas(aftersale)
	configureActions(Admin, aftersale)
	configureScopes(aftersale)

	
	// aftersale.UseTheme("grid")
	// aftersale.UseTheme("publish2")
	aftersale.UseTheme("fancy")


	// aftersale.FindManyHandler = func(results interface{}, context *qor.Context) error {
	// 	db             = context.GetDB()
	// 	scope          = db.NewScope(record)

	// 	// find records and decode them to results
	// 	return nil
	// }
}

func configureMetas(model *admin.Resource) {
	model.EditAttrs("-UserID", "-User", "-CreatedAt", "-UpdatedAt", "-CreatedBy", "-UpdatedBy", "-State")
	model.NewAttrs("-UserID",  "-User", "-CreatedAt", "-UpdatedAt", "-CreatedBy", "-UpdatedBy", "-State")
	model.IndexAttrs("-UserID", "-CreatedAt", "-UpdatedAt", "-CreatedBy", "-UpdatedBy", "-Fee", "-Remark")
}

func configureScopes(model *admin.Resource) {
	// filter by order source
	// var brands = []settings.Brand{
	// 	settings.Brand{
	// 		Name: "海尔",
	// 	},
	// 	settings.Brand{
	// 		Name: "格力",
	// 	},
	// }
	for _, item := range brands {
		var item = item
		model.Scope(&admin.Scope{
			Name:  item.Name,
			Label: item.Name,
			Group: "Filter By Source",
			Handler: func(db *gorm.DB, context *qor.Context) *gorm.DB {
				return db.Where("source = ?", item.Name)
			},
		})
	}

	// filter by order state
	for _, item := range aftersales.STATES {
		var item = item // 这句必须有否则会报错，永远都是最后一个值
		model.Scope(&admin.Scope{
			Name:  item,
			Label: item,
			Group: "Filter By State",
			Handler: func(db *gorm.DB, context *qor.Context) *gorm.DB {
				// 两种写法都可以
				// return db.Where(aftersales.AfterSale{State: item})
				return db.Where("state = ?", item)
			},
		})
	}

	// filter by order type
	for _, item := range service_types {
		var item = item // 这句必须有否则会报错，永远都是最后一个值
		model.Scope(&admin.Scope{
			Name:  item.Name,
			Label: item.Name,
			Group: "Filter By Type",
			Handler: func(db *gorm.DB, context *qor.Context) *gorm.DB {
				// 两种写法都可以
				return db.Where("service_type = ?", item.Name)
			},
		})
	}
}

func configureActions(Admin *admin.Admin, aftersale *admin.Resource) {
	// 安排
	type setupActionArgument struct {
		UserID uint
	}

	setupActionArgumentResource := Admin.NewResource(&setupActionArgument{})
	setupActionArgumentResource.Meta(&admin.Meta{
		Name: "UserID",
		Type: "select_one",
		Valuer: func(record interface{}, context *qor.Context) interface{} {
			// return record.(*users.User).ID
			return ""
		},
		Collection: func(value interface{}, context *qor.Context) (options [][]string) {
			var setupMen []users.User
			// context.GetDB().Where("role = ?", "setup_man").Find(&setupMen)
			context.GetDB().Find(&setupMen)
			for _, m := range setupMen {
				idStr := fmt.Sprintf("%d", m.ID)
				var option = []string{idStr, m.Name}
				options = append(options, option)
			}
			return options
		},
		// Collection: []string{"Male", "Female", "Unknown"},
	})

	// 安排安装
	aftersale.Action(&admin.Action{
		Name: "安排",
		Handler: func(argument *admin.ActionArgument) error {
			var (
				tx  = argument.Context.GetDB().Begin()
				arg = argument.Argument.(*setupActionArgument)
			)
			for _, record := range argument.FindSelectedRecords() {
				item := record.(*aftersales.AfterSale)
				item.UserID = &arg.UserID
				// orders.OrderState.Trigger("schedule_setup", order, tx, "man to setup: "+arg.ManToSetup)
				if err := tx.Save(item).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
			tx.Commit()
			return nil
		},
		Visible: func(record interface{}, context *admin.Context) bool {
			// if order, ok := record.(*orders.Order); ok {
			// 	return order.State == "processing"
			// }
			// return false
			return true
		},
		Resource: setupActionArgumentResource,
		Modes:    []string{"edit", "show", "menu_item"},
	})

}
