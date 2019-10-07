package settings

import (
	"github.com/jinzhu/gorm"
)

type ServiceType struct {
	gorm.Model
	Name string
}
