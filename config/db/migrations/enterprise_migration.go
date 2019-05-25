// +build enterprise

package migrations

import "github.com/dfang/qor-demo/app/enterprise"

func init() {
	AutoMigrate(&enterprise.QorMicroSite{})
}
