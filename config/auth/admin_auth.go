package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dfang/qor-demo/models/users"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/qor/roles"
	"github.com/rs/zerolog/log"
)

func init() {
	roles.Register("admin", func(req *http.Request, currentUser interface{}) bool {
		return currentUser != nil && strings.ToLower(currentUser.(*users.User).Role) == "admin"
	})

	roles.Register("operator", func(req *http.Request, currentUser interface{}) bool {
		return currentUser != nil && strings.ToLower(currentUser.(*users.User).Role) == "operator"
	})

	roles.Register("workman", func(req *http.Request, currentUser interface{}) bool {
		return currentUser != nil && strings.ToLower(currentUser.(*users.User).Role) == "workman"
	})
}

type AdminAuth struct {
}

func (AdminAuth) LoginURL(c *admin.Context) string {
	return "/auth/login"
}

func (AdminAuth) LogoutURL(c *admin.Context) string {
	return "/auth/logout"
}

func (AdminAuth) GetCurrentUser(c *admin.Context) qor.CurrentUser {
	currentUser := Auth.GetCurrentUser(c.Request)
	if currentUser != nil {
		qorCurrentUser, ok := currentUser.(qor.CurrentUser)
		if !ok {
			fmt.Printf("User %#v haven't implement qor.CurrentUser interface\n", currentUser)
		}

		log.Debug().Msgf("role for currentUser is %s ", currentUser.(*users.User).Role)
		return qorCurrentUser
	}
	return nil
}
