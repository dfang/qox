package seeds

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dfang/qor-demo/app/admin"
	"github.com/dfang/qor-demo/config/auth"
	"github.com/dfang/qor-demo/config/db"
	"github.com/qor/auth/auth_identity"
	"github.com/qor/auth/providers/password"

	"github.com/dfang/qor-demo/models/aftersales"
	"github.com/dfang/qor-demo/models/blogs"
	"github.com/dfang/qor-demo/models/orders"
	"github.com/dfang/qor-demo/models/products"
	adminseo "github.com/dfang/qor-demo/models/seo"
	"github.com/dfang/qor-demo/models/settings"
	"github.com/dfang/qor-demo/models/stores"
	"github.com/dfang/qor-demo/models/users"
	"github.com/jinzhu/now"
	"github.com/qor/activity"
	qoradmin "github.com/qor/admin"
	"github.com/qor/banner_editor"
	"github.com/qor/help"
	i18n_database "github.com/qor/i18n/backends/database"
	"github.com/qor/media"
	"github.com/qor/media/asset_manager"
	"github.com/qor/media/media_library"
	"github.com/qor/media/oss"
	"github.com/qor/notification"
	"github.com/qor/notification/channels/database"
	"github.com/qor/publish2"
	"github.com/qor/qor"
	"github.com/qor/seo"
	"github.com/qor/slug"
	"github.com/qor/sorting"
	"github.com/qor/transition"
	"syreclabs.com/go/faker"
)

/* How to upload file
 * $ brew install s3cmd
 * $ s3cmd --configure (Refer https://github.com/theplant/qor-example)
 * $ s3cmd put local_file_path s3://qor3/
 */

var (
	AdminUser    *users.User
	Notification = notification.New(&notification.Config{})
	Tables       = []interface{}{
		// users
		&auth_identity.AuthIdentity{},
		&users.User{},
		&users.Address{},

		// orders
		&orders.Order{},
		&orders.OrderItem{},
		&orders.DeliveryMethod{},

		// notification, activity, i18n, state machine
		&activity.QorActivity{},
		&i18n_database.Translation{},
		&transition.StateChangeLog{},
		&notification.QorNotification{},

		// stores
		&stores.Store{},

		// settings
		&settings.Setting{},
		&settings.MediaLibrary{},
		&admin.QorWidgetSetting{},
		&adminseo.MySEOSetting{},

		&asset_manager.AssetManager{},
		&banner_editor.QorBannerEditorSetting{},
		&blogs.Page{}, &blogs.Article{},
		&help.QorHelpEntry{},

		// products
		&products.Category{},
		&products.Color{},
		&products.Size{},
		&products.Material{},
		&products.Collection{},
		&products.Product{},
		&products.ProductImage{},
		&products.ColorVariation{},
		&products.SizeVariation{},
		&products.ColorVariationImage{},
	}
)

func Run() {
	Notification.RegisterChannel(database.New(&database.Config{}))

	// fmt.Println("Truncate tables .....")
	// TruncateTables(Tables...)

	// fmt.Println("Run migrations .....")
	// migrations.Migrate()

	fmt.Println("Create root user .....")
	CreateRootUser()

	fmt.Println("Create aftersales .....")
	CreateAftersales()

	fmt.Println("Create sources .....")
	CreateSources()

	fmt.Println("Create brands .....")
	CreateBrands()

	fmt.Println("Create serviceTypes .....")
	CreateServiceTypes()

	fmt.Println("Create workman .....")
	CreateWorkman()
	// createRecords()
}

func CreateSources() {
	sources := []string{
		"京东",
		"苏宁",
		"天猫",
		"线下",
	}

	for i := 0; i < len(sources); i++ {
		a := settings.Source{
			Name: sources[i],
		}

		if err := DraftDB.Create(&a).Error; err != nil {
			log.Fatalf("create source (%v) failure, got err %v", a, err)
		}
	}
}

func CreateServiceTypes() {
	serviceTypes := []string{
		"安装调试服务",
		"售后维修服务",
		"清洗养护服务",
		"家装保洁服务",
	}

	// var afs []settings.ServiceType
	for i := 0; i < len(serviceTypes); i++ {
		a := settings.ServiceType{
			Name: serviceTypes[i],
		}

		if err := DraftDB.Create(&a).Error; err != nil {
			log.Fatalf("create service type (%v) failure, got err %v", a, err)
		}
	}
}

func CreateBrands() {
	brands := []string{
		"海尔",
		"格力",
		"奥克斯",
		"小米电视",
		"春兰空调",
		"长虹",
		"康佳",
	}

	for i := 0; i < len(brands); i++ {
		a := settings.Brand{
			Name: brands[i],
		}

		if err := DraftDB.Create(&a).Error; err != nil {
			log.Fatalf("create brand (%v) failure, got err %v", a, err)
		}
	}
}

func CreateAftersales() {
	serviceTypes := []string{
		"安装调试服务",
		"售后维修服务",
	}
	sources := []string{
		"京东",
		"苏宁",
		"天猫",
		"线下",
	}

	brands := []string{
		"格力",
		"海尔",
		"TCL",
		"海信",
		"创维",
		"长虹",
		"美的",
		"小米",
		"奥克斯",
		"康佳",
	}

	serviceContents := []string{
		"电视安装",
		"电视维修",
		"冰箱维修",
		"冰箱安装",
		"空调安装",
		"空调维修",
		"洗衣机安装",
		"洗衣机维修",
		"油烟机维修",
	}

	fees := []int{
		50, 60, 80, 100, 120, 150, 180, 200,
	}

	rand.Seed(time.Now().UnixNano())

	var afs []aftersales.Aftersale
	for i := 0; i < 100; i++ {
		a := aftersales.Aftersale{
			CustomerName:    faker.Name().Name(),
			CustomerPhone:   faker.PhoneNumber().CellPhone(),
			CustomerAddress: faker.Address().StreetAddress(),

			ServiceType:    serviceTypes[rand.Intn(len(serviceTypes)-1)],
			ServiceContent: serviceContents[rand.Intn(len(serviceContents)-1)],
			Source:         sources[rand.Intn(len(sources)-1)],
			Brand:          brands[rand.Intn(len(brands)-1)],
			Fee:            float32(fees[rand.Intn(len(fees)-1)]),
			// Fee:            float32(rand.Intn(101)),
		}

		afs = append(afs, a)
	}

	for _, s := range afs {
		if err := DraftDB.Create(&s).Error; err != nil {
			log.Fatalf("create aftersale (%v) failure, got err %v", s, err)
		}
	}
}

func CreateWorkman() {
	a := users.User{
		Name:        "段访",
		MobilePhone: "15618903080",
		Role:        "workman",
	}
	b := users.User{
		Name:        "朱大平",
		MobilePhone: "18970278113",
		Role:        "workman",
	}

	var names []users.User
	names = append(names, a)
	names = append(names, b)

	for i := 0; i < len(names); i++ {
		if err := DraftDB.Create(&names[i]).Error; err != nil {
			log.Fatalf("create workman (%v) failure, got err %v", a, err)
		}
	}

	wp1 := users.WechatProfile{
		Openid:      "oLROJs729qr09CqFLFx03eGAAHU8",
		Unionid:     "opDDlslN3CL8zeH8AuW_LW_pNBoM",
		Nickname:    "Fang",
		MobilePhone: "15618903080",
	}

	if err := DraftDB.Create(&wp1).Error; err != nil {
		log.Fatalf("create wechat profiles (%v) failure, got err %v", a, err)
	}

	wp2 := users.WechatProfile{
		Openid:      "oLROJs4Zr8rn-Azsbs7_7fr2MLdU",
		Unionid:     "opDDlstLHxfJU3budoZjto1WDR3Y",
		Nickname:    "修水京东18070421816朱大平",
		MobilePhone: "18970278113",
	}

	if err := DraftDB.Create(&wp2).Error; err != nil {
		log.Fatalf("create wechat profiles (%v) failure, got err %v", a, err)
	}
}

func importUsers() {
	for _, s := range Seeds.Users {
		user := users.User{}
		user.Name = s.Name
		user.Gender = s.Gender
		user.IdentityCardNum = s.IdentityCardNum
		user.MobilePhone = s.MobilePhone
		user.CarLicencePlateNum = s.CarLicencePlateNum
		user.CarType = s.CarType
		user.CarLicenseType = s.CarLicenseType
		user.IsCasual = s.IsCasual
		// user.Type = s.Type
		user.Role = s.Role
		user.JDAppUser = s.JDAppUser
		// user.HireDate =

		emailRegexp := regexp.MustCompile(".*(@.*)")
		user.Email = emailRegexp.ReplaceAllString(Fake.Email(), strings.Replace(strings.ToLower(Fake.Name()), " ", "_", -1)+"@example.com")

		if err := DraftDB.Create(&user).Error; err != nil {
			log.Fatalf("Save user (%v) failure, got err %v", user, err)
		}
		provider := auth.Auth.GetProvider("password").(*password.Provider)
		hashedPassword, _ := provider.Encryptor.Digest("password")
		now := time.Now()
		authIdentity := &auth_identity.AuthIdentity{}
		authIdentity.Provider = "password"
		authIdentity.UID = user.Email
		authIdentity.EncryptedPassword = hashedPassword
		authIdentity.UserID = fmt.Sprint(user.ID)
		authIdentity.ConfirmedAt = &now

		DraftDB.Create(authIdentity)

		// Send welcome notification
		Notification.Send(&notification.Message{
			From:        AdminUser,
			To:          AdminUser,
			Title:       "Welcome To QOR Admin",
			Body:        "Welcome To QOR Admin",
			MessageType: "info",
		}, &qor.Context{DB: DraftDB})

	}
}

func createRecords() {
	fmt.Println("Start create sample data...")

	createSetting()
	fmt.Println("--> Created setting.")

	// createSeo()
	// fmt.Println("--> Created seo.")

	CreateRootUser()
	fmt.Println("--> Created admin users.")

	importUsers()
	fmt.Println("--> Imported users from yaml.")

	createUsers()
	fmt.Println("--> Created users.")
	createAddresses()
	fmt.Println("--> Created addresses.")

	createCategories()
	fmt.Println("--> Created categories.")
	createCollections()
	fmt.Println("--> Created collections.")
	createColors()
	fmt.Println("--> Created colors.")
	createSizes()
	fmt.Println("--> Created sizes.")
	createMaterial()
	fmt.Println("--> Created material.")

	createMediaLibraries()
	fmt.Println("--> Created medialibraries.")

	createProducts()
	fmt.Println("--> Created products.")

	createStores()
	fmt.Println("--> Created stores.")

	// createOrders()
	// fmt.Println("--> Created orders.")

	createWidgets()
	fmt.Println("--> Created widgets.")

	createArticles()
	fmt.Println("--> Created articles.")

	createHelps()
	fmt.Println("--> Created helps.")

	fmt.Println("--> Done!")
}

// CreateRootUser the system should have at least one root user
func CreateRootUser() {
	AdminUser = &users.User{}
	AdminUser.Email = "admin@example.com"
	AdminUser.Confirmed = true
	AdminUser.Name = "QOR Admin"
	AdminUser.Role = "Admin"
	DraftDB.FirstOrCreate(AdminUser)

	provider := auth.Auth.GetProvider("password").(*password.Provider)
	hashedPassword, _ := provider.Encryptor.Digest("admin")
	now := time.Now()

	authIdentity := &auth_identity.AuthIdentity{}
	authIdentity.Provider = "password"
	authIdentity.UID = AdminUser.Email
	authIdentity.EncryptedPassword = hashedPassword
	authIdentity.UserID = fmt.Sprint(AdminUser.ID)
	authIdentity.ConfirmedAt = &now

	DraftDB.FirstOrCreate(authIdentity)
}

func createSetting() {
	setting := settings.Setting{}

	setting.ShippingFee = Seeds.Setting.ShippingFee
	setting.GiftWrappingFee = Seeds.Setting.GiftWrappingFee
	setting.CODFee = Seeds.Setting.CODFee
	setting.TaxRate = Seeds.Setting.TaxRate
	setting.Address = Seeds.Setting.Address
	setting.Region = Seeds.Setting.Region
	setting.City = Seeds.Setting.City
	setting.Country = Seeds.Setting.Country
	setting.Zip = Seeds.Setting.Zip
	setting.Latitude = Seeds.Setting.Latitude
	setting.Longitude = Seeds.Setting.Longitude

	if err := DraftDB.Create(&setting).Error; err != nil {
		log.Fatalf("create setting (%v) failure, got err %v", setting, err)
	}
}

func createSeo() {
	globalSeoSetting := adminseo.MySEOSetting{}
	globalSetting := make(map[string]string)
	globalSetting["SiteName"] = "Qor Demo"
	globalSeoSetting.Setting = seo.Setting{GlobalSetting: globalSetting}
	globalSeoSetting.Name = "QorSeoGlobalSettings"
	globalSeoSetting.LanguageCode = "en-US"
	globalSeoSetting.QorSEOSetting.SetIsGlobalSEO(true)

	if err := db.DB.Create(&globalSeoSetting).Error; err != nil {
		log.Fatalf("create seo (%v) failure, got err %v", globalSeoSetting, err)
	}

	defaultSeo := adminseo.MySEOSetting{}
	defaultSeo.Setting = seo.Setting{Title: "{{SiteName}}", Description: "{{SiteName}} - Default Description", Keywords: "{{SiteName}} - Default Keywords", Type: "Default Page"}
	defaultSeo.Name = "Default Page"
	defaultSeo.LanguageCode = "en-US"
	if err := db.DB.Create(&defaultSeo).Error; err != nil {
		log.Fatalf("create seo (%v) failure, got err %v", defaultSeo, err)
	}

	productSeo := adminseo.MySEOSetting{}
	productSeo.Setting = seo.Setting{Title: "{{SiteName}}", Description: "{{SiteName}} - {{Name}} - {{Code}}", Keywords: "{{SiteName}},{{Name}},{{Code}}", Type: "Product Page"}
	productSeo.Name = "Product Page"
	productSeo.LanguageCode = "en-US"
	if err := db.DB.Create(&productSeo).Error; err != nil {
		log.Fatalf("create seo (%v) failure, got err %v", productSeo, err)
	}

	// seoSetting := models.SEOSetting{}
	// seoSetting.SiteName = Seeds.Seo.SiteName
	// seoSetting.DefaultPage = seo.Setting{Title: Seeds.Seo.DefaultPage.Title, Description: Seeds.Seo.DefaultPage.Description, Keywords: Seeds.Seo.DefaultPage.Keywords}
	// seoSetting.HomePage = seo.Setting{Title: Seeds.Seo.HomePage.Title, Description: Seeds.Seo.HomePage.Description, Keywords: Seeds.Seo.HomePage.Keywords}
	// seoSetting.ProductPage = seo.Setting{Title: Seeds.Seo.ProductPage.Title, Description: Seeds.Seo.ProductPage.Description, Keywords: Seeds.Seo.ProductPage.Keywords}

	// if err := DraftDB.Create(&seoSetting).Error; err != nil {
	// 	log.Fatalf("create seo (%v) failure, got err %v", seoSetting, err)
	// }
}

func createUsers() {
	emailRegexp := regexp.MustCompile(".*(@.*)")
	totalCount := 60
	for i := 0; i < totalCount; i++ {
		user := users.User{}
		user.Name = Fake.Name()
		user.Email = emailRegexp.ReplaceAllString(Fake.Email(), strings.Replace(strings.ToLower(user.Name), " ", "_", -1)+"@example.com")
		user.Gender = []string{"Female", "Male"}[i%2]
		if err := DraftDB.Create(&user).Error; err != nil {
			log.Fatalf("create user (%v) failure, got err %v", user, err)
		}

		day := (-14 + i/45)
		user.CreatedAt = now.EndOfDay().Add(time.Duration(day*rand.Intn(24)) * time.Hour)
		if user.CreatedAt.After(time.Now()) {
			user.CreatedAt = time.Now()
		}
		if err := DraftDB.Save(&user).Error; err != nil {
			log.Fatalf("Save user (%v) failure, got err %v", user, err)
		}

		provider := auth.Auth.GetProvider("password").(*password.Provider)
		hashedPassword, _ := provider.Encryptor.Digest("admin")
		authIdentity := &auth_identity.AuthIdentity{}
		authIdentity.Provider = "password"
		authIdentity.UID = user.Email
		authIdentity.EncryptedPassword = hashedPassword
		authIdentity.UserID = fmt.Sprint(user.ID)
		authIdentity.ConfirmedAt = &user.CreatedAt

		DraftDB.Create(authIdentity)
	}
}

func createAddresses() {
	var Users []users.User
	if err := DraftDB.Find(&Users).Error; err != nil {
		log.Fatalf("query users (%v) failure, got err %v", Users, err)
	}

	for _, user := range Users {
		address := users.Address{}
		address.UserID = user.ID
		address.ContactName = user.Name
		address.Phone = Fake.PhoneNumber()
		address.City = Fake.City()
		address.Address1 = Fake.StreetAddress()
		address.Address2 = Fake.SecondaryAddress()
		if err := DraftDB.Create(&address).Error; err != nil {
			log.Fatalf("create address (%v) failure, got err %v", address, err)
		}
	}
}

func createCategories() {
	for _, c := range Seeds.Categories {
		category := products.Category{}
		category.Name = c.Name
		category.Code = strings.ToLower(c.Code)
		if err := DraftDB.Create(&category).Error; err != nil {
			log.Fatalf("create category (%v) failure, got err %v", category, err)
		}
	}
}

func createCollections() {
	for _, c := range Seeds.Collections {
		collection := products.Collection{}
		collection.Name = c.Name
		if err := DraftDB.Create(&collection).Error; err != nil {
			log.Fatalf("create collection (%v) failure, got err %v", collection, err)
		}
	}
}

func createColors() {
	for _, c := range Seeds.Colors {
		color := products.Color{}
		color.Name = c.Name
		color.Code = c.Code
		color.PublishReady = true
		if err := DraftDB.Create(&color).Error; err != nil {
			log.Fatalf("create color (%v) failure, got err %v", color, err)
		}
	}
}

func createSizes() {
	for _, s := range Seeds.Sizes {
		size := products.Size{}
		size.Name = s.Name
		size.Code = s.Code
		size.PublishReady = true
		if err := DraftDB.Create(&size).Error; err != nil {
			log.Fatalf("create size (%v) failure, got err %v", size, err)
		}
	}
}

func createMaterial() {
	for _, s := range Seeds.Materials {
		material := products.Material{}
		material.Name = s.Name
		material.Code = s.Code
		if err := DraftDB.Create(&material).Error; err != nil {
			log.Fatalf("create material (%v) failure, got err %v", material, err)
		}
	}
}

func createProducts() {
	for idx, p := range Seeds.Products {
		category := findCategoryByName(p.CategoryName)

		product := products.Product{}
		product.CategoryID = category.ID
		product.Name = p.Name
		product.NameWithSlug = slug.Slug{p.NameWithSlug}
		product.Code = p.Code
		product.Price = p.Price
		product.Description = p.Description
		product.MadeCountry = p.MadeCountry
		product.Gender = p.Gender
		product.PublishReady = true
		for _, c := range p.Collections {
			collection := findCollectionByName(c.Name)
			product.Collections = append(product.Collections, *collection)
		}

		if err := DraftDB.Create(&product).Error; err != nil {
			log.Fatalf("create product (%v) failure, got err %v", product, err)
		}

		for _, cv := range p.ColorVariations {
			color := findColorByName(cv.ColorName)

			colorVariation := products.ColorVariation{}
			colorVariation.ProductID = product.ID
			colorVariation.ColorID = color.ID
			colorVariation.ColorCode = cv.ColorCode

			for _, i := range cv.Images {
				image := products.ProductImage{Title: p.Name, SelectedType: "image"}
				if file, err := openFileByURL(i.URL); err != nil {
					fmt.Printf("open file (%q) failure, got err %v", i.URL, err)
				} else {
					defer file.Close()
					image.File.Scan(file)
				}
				if err := DraftDB.Create(&image).Error; err != nil {
					log.Fatalf("create color_variation_image (%v) failure, got err %v when %v", image, err, i.URL)
				} else {
					colorVariation.Images.Files = append(colorVariation.Images.Files, media_library.File{
						ID:  json.Number(fmt.Sprint(image.ID)),
						Url: image.File.URL(),
					})

					Admin := qoradmin.New(&qoradmin.AdminConfig{
						SiteName: "QOR DEMO",
						Auth:     auth.AdminAuth{},
						DB:       db.DB.Set(publish2.VisibleMode, publish2.ModeOff).Set(publish2.ScheduleMode, publish2.ModeOff),
					})

					colorVariation.Images.Crop(Admin.NewResource(&products.ProductImage{}), DraftDB, media_library.MediaOption{
						Sizes: map[string]*media.Size{
							"main":    {Width: 560, Height: 700},
							"icon":    {Width: 50, Height: 50},
							"preview": {Width: 300, Height: 300},
							"listing": {Width: 640, Height: 640},
						},
					})

					if len(product.MainImage.Files) == 0 {
						product.MainImage.Files = []media_library.File{{
							ID:  json.Number(fmt.Sprint(image.ID)),
							Url: image.File.URL(),
						}}
						DraftDB.Save(&product)
					}
				}
			}

			if err := DraftDB.Create(&colorVariation).Error; err != nil {
				log.Fatalf("create color_variation (%v) failure, got err %v", colorVariation, err)
			}
		}

		for _, cv := range p.SizeVariations {
			size := findSizeByName(cv.SizeName)

			sizeVariation := products.SizeVariation{}
			sizeVariation.ProductID = product.ID
			sizeVariation.SizeID = size.ID
			// sizeVariation.Name = size.Name

			if err := DraftDB.Create(&sizeVariation).Error; err != nil {
				log.Fatalf("create color_variation (%v) failure, got err %v", sizeVariation, err)
			}
		}

		product.Name = p.ZhName
		product.Description = p.ZhDescription
		product.MadeCountry = p.ZhMadeCountry
		product.Gender = p.ZhGender
		DraftDB.Set("l10n:locale", "zh-CN").Create(&product)

		if idx%3 == 0 {
			start := time.Now().AddDate(0, 0, idx-7)
			end := time.Now().AddDate(0, 0, idx-4)
			product.SetVersionName("v1")
			product.Name = p.Name + " - v1"
			product.Description = p.Description + " - v1"
			product.MadeCountry = p.MadeCountry
			product.Gender = p.Gender
			product.SetScheduledStartAt(&start)
			product.SetScheduledEndAt(&end)
			DraftDB.Save(&product)
		}

		if idx%2 == 0 {
			start := time.Now().AddDate(0, 0, idx-7)
			end := time.Now().AddDate(0, 0, idx-4)
			product.SetVersionName("v1")
			product.Name = p.ZhName + " - 版本 1"
			product.Description = p.ZhDescription + " - 版本 1"
			product.MadeCountry = p.ZhMadeCountry
			product.Gender = p.ZhGender
			product.SetScheduledStartAt(&start)
			product.SetScheduledEndAt(&end)
			DraftDB.Set("l10n:locale", "zh-CN").Save(&product)
		}
	}
}

func createStores() {
	for _, s := range Seeds.Stores {
		store := stores.Store{}
		store.StoreName = s.Name
		store.Phone = s.Phone
		store.Email = s.Email
		store.Country = s.Country
		store.City = s.City
		store.Region = s.Region
		store.Address = s.Address
		store.Zip = s.Zip
		store.Latitude = s.Latitude
		store.Longitude = s.Longitude
		if err := DraftDB.Create(&store).Error; err != nil {
			log.Fatalf("create store (%v) failure, got err %v", store, err)
		}
	}
}

func createOrders() {
	var Users []users.User
	if err := DraftDB.Preload("Addresses").Find(&Users).Error; err != nil {
		log.Fatalf("query users (%v) failure, got err %v", Users, err)
	}

	var sizeVariations []products.SizeVariation
	if err := DraftDB.Find(&sizeVariations).Error; err != nil {
		log.Fatalf("query sizeVariations (%v) failure, got err %v", sizeVariations, err)
	}

	var colorVariations []products.ColorVariation
	if err := DraftDB.Find(&colorVariations).Error; err != nil {
		log.Fatalf("query colorVariations (%v) failure, got err %v", colorVariations, err)
	}

	var sizeVariationsCount = len(sizeVariations)
	var colorVariationsCount = len(colorVariations)

	fmt.Println("sizeVariationsCount", sizeVariationsCount)
	fmt.Println("colorVariationsCount", colorVariationsCount)

	for _, user := range Users {
		count := 5
		if user.ID != 1 {
			count = rand.Intn(5)
		}

		for j := 0; j < count; j++ {
			order := orders.Order{}
			state := []string{"draft", "checkout", "cancelled", "paid", "paid_cancelled", "processing", "shipped", "returned"}[rand.Intn(10)%8]
			abandonedReasons := []string{
				"Unsatisfied with discount",
				"Dropped after check gift wrapping option",
				"Dropped after select expected delivery date",
				"Invalid credit card inputted",
				"Credit card balances insufficient",
				"Created a new order with more products",
				"Created a new order with fewer products",
			}
			abandonedReason := abandonedReasons[rand.Intn(len(abandonedReasons))]

			order.UserID = &user.ID
			order.ShippingAddressID = user.Addresses[0].ID
			order.BillingAddressID = user.Addresses[0].ID
			order.State = state
			if rand.Intn(15)%15 == 3 && state == "checkout" || state == "processing" || state == "paid_cancelled" {
				order.AbandonedReason = abandonedReason
			}
			if err := DraftDB.Create(&order).Error; err != nil {
				log.Fatalf("create order (%v) failure, got err %v", order, err)
			}

			sizeVariation := sizeVariations[rand.Intn(sizeVariationsCount)]
			colorVariation := colorVariations[rand.Intn(colorVariationsCount)]

			product := findProductByColorVariationID(colorVariation.ID)
			quantity := []uint{1, 2, 3, 4, 5}[rand.Intn(10)%5]
			discountRate := []uint{0, 5, 10, 15, 20, 25}[rand.Intn(10)%6]

			orderItem := orders.OrderItem{}
			orderItem.OrderID = order.ID
			orderItem.SizeVariationID = sizeVariation.ID
			orderItem.ColorVariationID = colorVariation.ID
			orderItem.Quantity = quantity
			orderItem.Price = product.Price
			orderItem.State = state
			orderItem.DiscountRate = discountRate
			if err := DraftDB.Create(&orderItem).Error; err != nil {
				log.Fatalf("create orderItem (%v) failure, got err %v", orderItem, err)
			}

			order.OrderItems = append(order.OrderItems, orderItem)
			order.CreatedAt = user.CreatedAt.Add(1 * time.Hour)
			order.PaymentAmount = order.Amount()
			order.PaymentMethod = orders.COD
			if err := DraftDB.Save(&order).Error; err != nil {
				log.Fatalf("Save order (%v) failure, got err %v", order, err)
			}

			var resolvedAt *time.Time
			if (rand.Intn(10) % 9) != 1 {
				now := time.Now()
				resolvedAt = &now
			}

			// Send welcome notification
			switch order.State {
			case "paid_cancelled":
				Notification.Send(&notification.Message{
					From:        user,
					To:          AdminUser,
					Title:       "Order Cancelled After Paid",
					Body:        fmt.Sprintf("Order #%v has been cancelled, its amount %.2f", order.ID, order.Amount()),
					MessageType: "order_paid_cancelled",
					ResolvedAt:  resolvedAt,
				}, &qor.Context{DB: DraftDB})
			case "processed":
				Notification.Send(&notification.Message{
					From:        user,
					To:          AdminUser,
					Title:       "Order Processed",
					Body:        fmt.Sprintf("Order #%v has been prepared to ship", order.ID),
					MessageType: "order_processed",
					ResolvedAt:  resolvedAt,
				}, &qor.Context{DB: DraftDB})
			case "returned":
				Notification.Send(&notification.Message{
					From:        user,
					To:          AdminUser,
					Title:       "Order Returned",
					Body:        fmt.Sprintf("Order #%v has been returned, its amount %.2f", order.ID, order.Amount()),
					MessageType: "order_returned",
					ResolvedAt:  resolvedAt,
				}, &qor.Context{DB: DraftDB})
			}
		}
	}
}

func createMediaLibraries() {
	for _, m := range Seeds.MediaLibraries {
		medialibrary := settings.MediaLibrary{}
		medialibrary.Title = m.Title

		if file, err := openFileByURL(m.Image); err != nil {
			fmt.Printf("open file (%q) failure, got err %v", m.Image, err)
		} else {
			defer file.Close()
			medialibrary.File.Scan(file)
		}

		if err := DraftDB.Create(&medialibrary).Error; err != nil {
			log.Fatalf("create medialibrary (%v) failure, got err %v", medialibrary, err)
		}
	}
}

func createWidgets() {
	// home page banner
	type ImageStorage struct{ oss.OSS }
	topBannerSetting := admin.QorWidgetSetting{}
	topBannerSetting.Name = "home page banner"
	topBannerSetting.Description = "This is a top banner"
	topBannerSetting.WidgetType = "NormalBanner"
	topBannerSetting.GroupName = "Banner"
	topBannerSetting.Scope = "from_google"
	topBannerSetting.Shared = true
	topBannerValue := &struct {
		Title           string
		ButtonTitle     string
		Link            string
		BackgroundImage ImageStorage `sql:"type:varchar(4096)"`
		Logo            ImageStorage `sql:"type:varchar(4096)"`
	}{
		Title:       "Welcome Googlistas!",
		ButtonTitle: "LEARN MORE",
		Link:        "http://getqor.com",
	}
	if file, err := openFileByURL("http://qor3.s3.amazonaws.com/slide01.jpg"); err == nil {
		defer file.Close()
		topBannerValue.BackgroundImage.Scan(file)
	} else {
		fmt.Printf("open file (%q) failure, got err %v", "banner", err)
	}

	if file, err := openFileByURL("http://qor3.s3.amazonaws.com/qor_logo.png"); err == nil {
		defer file.Close()
		topBannerValue.Logo.Scan(file)
	} else {
		fmt.Printf("open file (%q) failure, got err %v", "qor_logo", err)
	}

	topBannerSetting.SetSerializableArgumentValue(topBannerValue)
	if err := DraftDB.Create(&topBannerSetting).Error; err != nil {
		log.Fatalf("Save widget (%v) failure, got err %v", topBannerSetting, err)
	}

	// SlideShow banner
	type slideImage struct {
		Title    string
		SubTitle string
		Button   string
		Link     string
		Image    oss.OSS
	}
	slideshowSetting := admin.QorWidgetSetting{}
	slideshowSetting.Name = "home page banner"
	slideshowSetting.GroupName = "Banner"
	slideshowSetting.WidgetType = "SlideShow"
	slideshowSetting.Scope = "default"
	slideshowValue := &struct {
		SlideImages []slideImage
	}{}

	for _, s := range Seeds.Slides {
		slide := slideImage{Title: s.Title, SubTitle: s.SubTitle, Button: s.Button, Link: s.Link}
		if file, err := openFileByURL(s.Image); err == nil {
			defer file.Close()
			slide.Image.Scan(file)
		} else {
			fmt.Printf("open file (%q) failure, got err %v", "banner", err)
		}
		slideshowValue.SlideImages = append(slideshowValue.SlideImages, slide)
	}
	slideshowSetting.SetSerializableArgumentValue(slideshowValue)
	if err := DraftDB.Create(&slideshowSetting).Error; err != nil {
		fmt.Printf("Save widget (%v) failure, got err %v", slideshowSetting, err)
	}

	// Featured Products
	featureProducts := admin.QorWidgetSetting{}
	featureProducts.Name = "featured products"
	featureProducts.Description = "featured product list"
	featureProducts.WidgetType = "Products"
	featureProducts.SetSerializableArgumentValue(&struct {
		Products       []string
		ProductsSorter sorting.SortableCollection
	}{
		Products:       []string{"1", "2", "3", "4", "5", "6", "7", "8"},
		ProductsSorter: sorting.SortableCollection{PrimaryKeys: []string{"1", "2", "3", "4", "5", "6", "7", "8"}},
	})
	if err := DraftDB.Create(&featureProducts).Error; err != nil {
		log.Fatalf("Save widget (%v) failure, got err %v", featureProducts, err)
	}

	// Banner edit items
	for _, s := range Seeds.BannerEditorSettings {
		setting := banner_editor.QorBannerEditorSetting{}
		id, _ := strconv.Atoi(s.ID)
		setting.ID = uint(id)
		setting.Kind = s.Kind
		setting.Value.SerializedValue = s.Value
		if err := DraftDB.Create(&setting).Error; err != nil {
			log.Fatalf("Save QorBannerEditorSetting (%v) failure, got err %v", setting, err)
		}
	}

	// Men collection
	menCollectionWidget := admin.QorWidgetSetting{}
	menCollectionWidget.Name = "men collection"
	menCollectionWidget.Description = "Men collection baner"
	menCollectionWidget.WidgetType = "FullWidthBannerEditor"
	menCollectionWidget.Value.SerializedValue = `{"Value":"%3Cdiv+class%3D%22qor-bannereditor__html%22+style%3D%22position%3A+relative%3B+height%3A+100%25%3B%22+data-image-width%3D%221280%22+data-image-height%3D%22480%22%3E%3Cspan+class%3D%22qor-bannereditor-image%22%3E%3Cimg+src%3D%22http%3A%2F%2Fqor3.s3.amazonaws.com%2Fmedialibrary%2Fmen-collection.jpg%22%3E%3C%2Fspan%3E%3Cspan+class%3D%22qor-bannereditor__draggable%22+data-edit-id%3D%2212%22+style%3D%22position%3A+absolute%3B+left%3A+10.0781%25%3B+top%3A+18.125%25%3B+right%3A+auto%3B+bottom%3A+auto%3B%22+data-position-left%3D%22129%22+data-position-top%3D%2287%22%3E%3Ch1+class%3D%22banner-title%22+style%3D%22color%3A+%3B%22%3EMEN+COLLECTION%3C%2Fh1%3E%3C%2Fspan%3E%3Cspan+class%3D%22qor-bannereditor__draggable%22+data-edit-id%3D%2210%22+style%3D%22position%3A+absolute%3B+left%3A+9.92188%25%3B+top%3A+29.7917%25%3B+right%3A+auto%3B+bottom%3A+auto%3B%22+data-position-left%3D%22127%22+data-position-top%3D%22143%22%3E%3Ch2+class%3D%22banner-sub-title%22+style%3D%22color%3A+%3B%22%3ECheck+the+newcomming+collection%3C%2Fh2%3E%3C%2Fspan%3E%3Cspan+class%3D%22qor-bannereditor__draggable+qor-bannereditor__draggable-left%22+data-edit-id%3D%2211%22+style%3D%22position%3A+absolute%3B+left%3A+9.92188%25%3B+top%3A+47.0833%25%3B+right%3A+auto%3B+bottom%3A+auto%3B%22+data-position-left%3D%22127%22+data-position-top%3D%22226%22%3E%3Ca+class%3D%22button+button__primary+banner-button%22+href%3D%22%2Fmen%22%3Eview+more%3C%2Fa%3E%3C%2Fspan%3E%3C%2Fdiv%3E"}`
	if err := DraftDB.Create(&menCollectionWidget).Error; err != nil {
		log.Fatalf("Save widget (%v) failure, got err %v", menCollectionWidget, err)
	}

	// Women collection
	womenCollectionWidget := admin.QorWidgetSetting{}
	womenCollectionWidget.Name = "women collection"
	womenCollectionWidget.Description = "Women collection banner"
	womenCollectionWidget.WidgetType = "FullWidthBannerEditor"
	womenCollectionWidget.Value.SerializedValue = `{"Value":"%3Cdiv+class%3D%22qor-bannereditor__html%22+style%3D%22position%3A+relative%3B+height%3A+100%25%3B%22+data-image-width%3D%221280%22+data-image-height%3D%22480%22%3E%3Cspan+class%3D%22qor-bannereditor-image%22%3E%3Cimg+src%3D%22http%3A%2F%2Fqor3.s3.amazonaws.com%2Fmedialibrary%2Fwomen-collection.jpg%22%3E%3C%2Fspan%3E%3Cspan+class%3D%22qor-bannereditor__draggable%22+data-edit-id%3D%2223%22+style%3D%22position%3A+absolute%3B+left%3A+10.0781%25%3B+top%3A+18.125%25%3B+right%3A+auto%3B+bottom%3A+auto%3B%22+data-position-left%3D%22129%22+data-position-top%3D%2287%22%3E%3Ch1+class%3D%22banner-title%22+style%3D%22color%3A+%3B%22%3EWOMEN+COLLECTION%3C%2Fh1%3E%3C%2Fspan%3E%3Cspan+class%3D%22qor-bannereditor__draggable%22+data-edit-id%3D%2221%22+style%3D%22position%3A+absolute%3B+left%3A+9.92188%25%3B+top%3A+29.7917%25%3B+right%3A+auto%3B+bottom%3A+auto%3B%22+data-position-left%3D%22127%22+data-position-top%3D%22143%22%3E%3Ch2+class%3D%22banner-sub-title%22+style%3D%22color%3A+%3B%22%3ECheck+the+newcomming+collection%3C%2Fh2%3E%3C%2Fspan%3E%3Cspan+class%3D%22qor-bannereditor__draggable+qor-bannereditor__draggable-left%22+data-edit-id%3D%2222%22+style%3D%22position%3A+absolute%3B+left%3A+9.92188%25%3B+top%3A+47.0833%25%3B+right%3A+auto%3B+bottom%3A+auto%3B%22+data-position-left%3D%22127%22+data-position-top%3D%22226%22%3E%3Ca+class%3D%22button+button__primary+banner-button%22+href%3D%22%2Fwomen%22%3Eview+more%3C%2Fa%3E%3C%2Fspan%3E%3C%2Fdiv%3E"}`
	if err := DraftDB.Create(&womenCollectionWidget).Error; err != nil {
		log.Fatalf("Save widget (%v) failure, got err %v", womenCollectionWidget, err)
	}

	// New arrivals promotio
	newArrivalsCollectionWidget := admin.QorWidgetSetting{}
	newArrivalsCollectionWidget.Name = "new arrivals promotion"
	newArrivalsCollectionWidget.Description = "New arrivals promotion banner"
	newArrivalsCollectionWidget.WidgetType = "FullWidthBannerEditor"
	newArrivalsCollectionWidget.Value.SerializedValue = `{"Value":"%3Cdiv%20class%3D%22qor-bannereditor__html%22%20style%3D%22position%3A%20relative%3B%20height%3A%20100%25%3B%22%20data-image-width%3D%221172%22%20data-image-height%3D%22300%22%3E%3Cspan%20class%3D%22qor-bannereditor-image%22%3E%3Cimg%20src%3D%22http%3A%2F%2Fqor3.s3.amazonaws.com%2Fmedialibrary%2Fnew-arrivals-bg.jpg%22%3E%3C%2Fspan%3E%3Cspan%20class%3D%22qor-bannereditor__draggable%20qor-bannereditor__draggable-left%22%20data-edit-id%3D%2233%22%20style%3D%22position%3A%20absolute%3B%20left%3A%209.47099%25%3B%20top%3A%2030%25%3B%20right%3A%20auto%3B%20bottom%3A%20auto%3B%22%20data-position-left%3D%22111%22%20data-position-top%3D%2290%22%3E%3Ch1%20class%3D%22banner-title%22%20style%3D%22color%3A%20%3B%22%3ENew%20Arrivals%3C%2Fh1%3E%3C%2Fspan%3E%3Cspan%20class%3D%22qor-bannereditor__draggable%20qor-bannereditor__draggable-left%22%20data-edit-id%3D%2231%22%20style%3D%22position%3A%20absolute%3B%20left%3A%208.61775%25%3B%20top%3A%20auto%3B%20right%3A%20auto%3B%20bottom%3A%2030.6667%25%3B%22%20data-position-left%3D%22101%22%20data-position-top%3D%22173%22%3E%3Ca%20class%3D%22button%20button__primary%20banner-button%22%20href%3D%22%23%22%3ESHOP%20COLLECTION%3C%2Fa%3E%3C%2Fspan%3E%3Cspan%20class%3D%22qor-bannereditor__draggable%22%20data-edit-id%3D%2232%22%20style%3D%22position%3A%20absolute%3B%20left%3A%209.55631%25%3B%20top%3A%2016%25%3B%20right%3A%20auto%3B%20bottom%3A%20auto%3B%22%20data-position-left%3D%22112%22%20data-position-top%3D%2248%22%3E%3Cp%20class%3D%22banner-text%22%20style%3D%22color%3A%20%3B%22%3ETHE%20STYLE%20THAT%20FITS%20EVERYTHING%3C%2Fp%3E%3C%2Fspan%3E%3C%2Fdiv%3E"}`
	if err := DraftDB.Create(&newArrivalsCollectionWidget).Error; err != nil {
		log.Fatalf("Save widget (%v) failure, got err %v", newArrivalsCollectionWidget, err)
	}

	// Model products
	modelCollectionWidget := admin.QorWidgetSetting{}
	modelCollectionWidget.Name = "model products"
	modelCollectionWidget.Description = "Model products banner"
	modelCollectionWidget.WidgetType = "FullWidthBannerEditor"
	modelCollectionWidget.Value.SerializedValue = `{"Value":"%3Cdiv%20class%3D%22qor-bannereditor__html%22%20style%3D%22position%3A%20relative%3B%20height%3A%20100%25%3B%22%20data-image-width%3D%221100%22%20data-image-height%3D%221200%22%3E%3Cspan%20class%3D%22qor-bannereditor-image%22%3E%3Cimg%20src%3D%22http%3A%2F%2Fqor3.s3.amazonaws.com%2Fmedialibrary%2Fmodel-products.jpg%22%3E%3C%2Fspan%3E%3Cspan%20class%3D%22qor-bannereditor__draggable%22%20data-edit-id%3D%2249%22%20style%3D%22position%3A%20absolute%3B%20left%3A%2026.4545%25%3B%20top%3A%204.41667%25%3B%20right%3A%20auto%3B%20bottom%3A%20auto%3B%22%20data-position-left%3D%22291%22%20data-position-top%3D%2253%22%3E%3Ch1%20class%3D%22banner-title%22%20style%3D%22color%3A%20%3B%22%3EENJOY%20THE%20NEW%20FASHION%20EXPERIENCE%3C%2Fh1%3E%3C%2Fspan%3E%3Cspan%20class%3D%22qor-bannereditor__draggable%22%20data-edit-id%3D%2242%22%20style%3D%22position%3A%20absolute%3B%20left%3A%2043.2727%25%3B%20top%3A%208.41667%25%3B%20right%3A%20auto%3B%20bottom%3A%20auto%3B%22%20data-position-left%3D%22476%22%20data-position-top%3D%22101%22%3E%3Cp%20class%3D%22banner-text%22%20style%3D%22color%3A%20%3B%22%3ENew%20look%20of%202017%3C%2Fp%3E%3C%2Fspan%3E%3Cspan%20class%3D%22qor-bannereditor__draggable%20qor-bannereditor__draggable-left%22%20data-edit-id%3D%2243%22%20style%3D%22position%3A%20absolute%3B%20left%3A%205.45455%25%3B%20top%3A%2044.25%25%3B%20right%3A%20auto%3B%20bottom%3A%20auto%3B%22%20data-position-left%3D%2260%22%20data-position-top%3D%22531%22%3E%3Cdiv%20class%3D%22model-buy-block%22%3E%3Ch2%20class%3D%22banner-sub-title%22%3ETOP%3C%2Fh2%3E%3Cp%20class%3D%22banner-text%22%3E%2429.99%3C%2Fp%3E%3Ca%20class%3D%22button%20button__primary%20banner-button%22%20href%3D%22%23%22%3EVIEW%20DETAILS%3C%2Fa%3E%3C%2Fdiv%3E%3C%2Fspan%3E%3Cspan%20class%3D%22qor-bannereditor__draggable%22%20data-edit-id%3D%2244%22%20style%3D%22position%3A%20absolute%3B%20left%3A%20auto%3B%20top%3A%2050.8333%25%3B%20right%3A%209.58527%25%3B%20bottom%3A%20auto%3B%22%20data-position-left%3D%22841%22%20data-position-top%3D%22610%22%3E%3Cdiv%20class%3D%22model-buy-block%22%3E%3Ch2%20class%3D%22banner-sub-title%22%3EPINK%20JACKET%3C%2Fh2%3E%3Cp%20class%3D%22banner-text%22%3E%2469.99%3C%2Fp%3E%3Ca%20class%3D%22button%20button__primary%20banner-button%22%20href%3D%22%23%22%3EVIEW%20DETAILS%3C%2Fa%3E%3C%2Fdiv%3E%3C%2Fspan%3E%3Cspan%20class%3D%22qor-bannereditor__draggable%20qor-bannereditor__draggable-left%22%20data-edit-id%3D%2247%22%20style%3D%22position%3A%20absolute%3B%20left%3A%2012.3636%25%3B%20top%3A%20auto%3B%20right%3A%20auto%3B%20bottom%3A%2014.2032%25%3B%22%20data-position-left%3D%22136%22%20data-position-top%3D%22903%22%3E%3Cdiv%20class%3D%22model-buy-block%22%3E%3Ch2%20class%3D%22banner-sub-title%22%3EBOTTOM%3C%2Fh2%3E%3Cp%20class%3D%22banner-text%22%3E%2432.99%3C%2Fp%3E%3Ca%20class%3D%22button%20button__primary%20banner-button%22%20href%3D%22%23%22%3EVIEW%20DETAILS%3C%2Fa%3E%3C%2Fdiv%3E%3C%2Fspan%3E%3Cspan%20class%3D%22qor-bannereditor__draggable%22%20data-edit-id%3D%2245%22%20style%3D%22position%3A%20absolute%3B%20left%3A%2053.2727%25%3B%20top%3A%2048.5%25%3B%20right%3A%20auto%3B%20bottom%3A%20auto%3B%22%20data-position-left%3D%22586%22%20data-position-top%3D%22582%22%3E%3Cimg%20src%3D%22%2F%2Fqor3.s3.amazonaws.com%2Fmedialibrary%2Farrow-left.png%22%20class%3D%22banner-image%22%3E%3C%2Fspan%3E%3Cspan%20class%3D%22qor-bannereditor__draggable%22%20data-edit-id%3D%2246%22%20style%3D%22position%3A%20absolute%3B%20left%3A%2015.5455%25%3B%20top%3A%2043.0833%25%3B%20right%3A%20auto%3B%20bottom%3A%20auto%3B%22%20data-position-left%3D%22171%22%20data-position-top%3D%22517%22%3E%3Cimg%20src%3D%22%2F%2Fqor3.s3.amazonaws.com%2Fmedialibrary%2Farrow-right.png%22%20class%3D%22banner-image%22%3E%3C%2Fspan%3E%3Cspan%20class%3D%22qor-bannereditor__draggable%22%20data-edit-id%3D%2248%22%20style%3D%22position%3A%20absolute%3B%20left%3A%2019.2727%25%3B%20top%3A%20auto%3B%20right%3A%20auto%3B%20bottom%3A%2024.8333%25%3B%22%20data-position-left%3D%22212%22%20data-position-top%3D%22879%22%3E%3Cimg%20src%3D%22%2F%2Fqor3.s3.amazonaws.com%2Fmedialibrary%2Farrow-right.png%22%20class%3D%22banner-image%22%3E%3C%2Fspan%3E%3C%2Fdiv%3E"}`
	if err := DraftDB.Create(&modelCollectionWidget).Error; err != nil {
		log.Fatalf("Save widget (%v) failure, got err %v", modelCollectionWidget, err)
	}
}

func createHelps() {
	helps := map[string][]string{
		"How to setup a microsite":           []string{"micro_sites"},
		"How to create a user":               []string{"users"},
		"How to create an admin user":        []string{"users"},
		"How to handle abandoned order":      []string{"abandoned_orders", "orders"},
		"How to cancel a order":              []string{"orders"},
		"How to create a order":              []string{"orders"},
		"How to upload product images":       []string{"products", "product_images"},
		"How to create a product":            []string{"products"},
		"How to create a discounted product": []string{"products"},
		"How to create a store":              []string{"stores"},
		"How shop setting works":             []string{"shop_settings"},
		"How to setup seo settings":          []string{"seo_settings"},
		"How to setup seo for blog":          []string{"seo_settings"},
		"How to setup seo for product":       []string{"seo_settings"},
		"How to setup seo for microsites":    []string{"micro_sites", "seo_settings"},
		"How to setup promotions":            []string{"promotions"},
		"How to publish a promotion":         []string{"schedules", "promotions"},
		"How to create a publish event":      []string{"schedules", "scheduled_events"},
		"How to publish a product":           []string{"schedules", "products"},
		"How to publish a microsite":         []string{"schedules", "micro_sites"},
		"How to create a scheduled data":     []string{"schedules"},
		"How to take something offline":      []string{"schedules"},
	}

	for key, value := range helps {
		helpEntry := help.QorHelpEntry{
			Title: key,
			Body:  "Content of " + key,
			Categories: help.Categories{
				Categories: value,
			},
		}
		DraftDB.Create(&helpEntry)
	}
}

func createArticles() {
	for idx := 1; idx <= 10; idx++ {
		title := fmt.Sprintf("Article %v", idx)
		article := blogs.Article{Title: title}
		article.PublishReady = true
		DraftDB.Create(&article)

		for i := 1; i <= idx-5; i++ {
			article.SetVersionName(fmt.Sprintf("v%v", i))
			start := time.Now().AddDate(0, 0, i*2-3)
			end := time.Now().AddDate(0, 0, i*2-1)
			article.SetScheduledStartAt(&start)
			article.SetScheduledEndAt(&end)
			DraftDB.Save(&article)
		}
	}
}

func findCategoryByName(name string) *products.Category {
	category := &products.Category{}
	if err := DraftDB.Where(&products.Category{Name: name}).First(category).Error; err != nil {
		log.Fatalf("can't find category with name = %q, got err %v", name, err)
	}
	return category
}

func findCollectionByName(name string) *products.Collection {
	collection := &products.Collection{}
	if err := DraftDB.Where(&products.Collection{Name: name}).First(collection).Error; err != nil {
		log.Fatalf("can't find collection with name = %q, got err %v", name, err)
	}
	return collection
}

func findColorByName(name string) *products.Color {
	color := &products.Color{}
	if err := DraftDB.Where(&products.Color{Name: name}).First(color).Error; err != nil {
		log.Fatalf("can't find color with name = %q, got err %v", name, err)
	}
	return color
}

func findSizeByName(name string) *products.Size {
	size := &products.Size{}
	if err := DraftDB.Where(&products.Size{Name: name}).First(size).Error; err != nil {
		log.Fatalf("can't find size with name = %q, got err %v", name, err)
	}
	return size
}

func findProductByColorVariationID(colorVariationID uint) *products.Product {
	colorVariation := products.ColorVariation{}
	product := products.Product{}

	if err := DraftDB.Find(&colorVariation, colorVariationID).Error; err != nil {
		log.Fatalf("query colorVariation (%v) failure, got err %v", colorVariation, err)
		return &product
	}
	if err := DraftDB.Find(&product, colorVariation.ProductID).Error; err != nil {
		log.Fatalf("query product (%v) failure, got err %v", product, err)
		return &product
	}
	return &product
}

func randTime() time.Time {
	num := rand.Intn(10)
	return time.Now().Add(-time.Duration(num*24) * time.Hour)
}

func openFileByURL(rawURL string) (*os.File, error) {
	if fileURL, err := url.Parse(rawURL); err != nil {
		return nil, err
	} else {
		path := fileURL.Path
		segments := strings.Split(path, "/")
		fileName := segments[len(segments)-1]

		filePath := filepath.Join(os.TempDir(), fileName)
		fmt.Println("filepath for openFileByURL is ", filePath)

		if _, err := os.Stat(filePath); err == nil {
			return os.Open(filePath)
		}

		file, err := os.Create(filePath)
		if err != nil {
			return file, err
		}

		check := http.Client{
			CheckRedirect: func(r *http.Request, via []*http.Request) error {
				r.URL.Opaque = r.URL.Path
				return nil
			},
		}
		resp, err := check.Get(rawURL) // add a filter to check redirect
		if err != nil {
			return file, err
		}
		defer resp.Body.Close()
		fmt.Printf("----> Downloaded %v\n", rawURL)

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return file, err
		}
		return file, nil
	}
}
