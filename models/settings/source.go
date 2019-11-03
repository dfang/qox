package settings

import (
	"github.com/jinzhu/gorm"
)

type Source struct {
	gorm.Model
	Name string
}
