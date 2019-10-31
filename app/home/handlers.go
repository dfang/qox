package home

import (
	"net/http"

	"github.com/qor/qor"
	"github.com/qor/qor/utils"
	"github.com/qor/render"
	"github.com/rs/zerolog/log"
)

// Controller home controller
type Controller struct {
	View *render.Render
}

// Index home index page
func (ctrl Controller) Index(w http.ResponseWriter, req *http.Request) {
	log.Debug().Msg("viewpaths for home/index")
	for _, v := range ctrl.View.ViewPaths {
		log.Debug().Msg(v)
	}
	ctrl.View.Execute("index", map[string]interface{}{}, req, w)
}

// SwitchLocale switch locale
func (ctrl Controller) SwitchLocale(w http.ResponseWriter, req *http.Request) {
	utils.SetCookie(http.Cookie{Name: "locale", Value: req.URL.Query().Get("locale")}, &qor.Context{Request: req, Writer: w})
	http.Redirect(w, req, req.Referer(), http.StatusSeeOther)
}

// RedirectToAdmin 重定向到/admin, 暂时屏蔽前端界面
func (ctrl Controller) RedirectToAdmin(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "/admin", 302)
}
