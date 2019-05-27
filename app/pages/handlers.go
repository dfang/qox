package pages

import (
	"net/http"

	"github.com/dfang/qor-demo/config/db"
	"github.com/dfang/qor-demo/models/blogs"
	"github.com/dfang/qor-demo/utils"
	"github.com/qor/render"
)

// Controller home controller
type Controller struct {
	View *render.Render
}

// Index home index page
func (ctrl Controller) Index(w http.ResponseWriter, req *http.Request) {
	ctrl.View.Layout("blog").Execute("index", map[string]interface{}{}, req, w)
}

func (ctrl Controller) Show(w http.ResponseWriter, req *http.Request) {
	var page blogs.Page
	db.DB.Where("title_with_slug = ?", utils.URLParam("article", req)).Find(&page)
	if page.ID == 0 {
		ctrl.View.Layout("blog").Execute("404", map[string]interface{}{"Page": page}, req, w)
	} else {
		ctrl.View.Layout("blog").Execute("show", map[string]interface{}{"Page": page}, req, w)
	}
}
