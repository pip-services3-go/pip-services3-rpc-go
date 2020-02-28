package auth

import (
	"net/http"
	"strings"

	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	cerr "github.com/pip-services3-go/pip-services3-commons-go/errors"
	services "github.com/pip-services3-go/pip-services3-rpc-go/services"
)

type RoleAuthorizer struct {
}

func (c *RoleAuthorizer) UserInRoles(roles []string) func(res http.ResponseWriter, req *http.Request, user *cdata.AnyValueMap, next http.HandlerFunc) {
	return func(res http.ResponseWriter, req *http.Request, user *cdata.AnyValueMap, next http.HandlerFunc) {

		if user == nil {
			services.HttpResponseSender.SendError(
				res, req,
				cerr.NewUnauthorizedError("", "NOT_SIGNED",
					"User must be signed in to perform this operation").WithStatus(401))
		} else {
			authorized := false
			userRoles := user.GetAsNullableArray("roles")

			if userRoles == nil {
				services.HttpResponseSender.SendError(
					res, req,
					cerr.NewUnauthorizedError("", "NOT_SIGNED",
						"User must be signed in to perform this operation").WithStatus(401))
				return
			}

			for _, role := range roles {
				for _, userRole := range userRoles.Value() {
					if role == userRole.(string) {
						authorized = true
					}
				}
			}

			if !authorized {

				services.HttpResponseSender.SendError(
					res, req,
					cerr.NewUnauthorizedError(
						"", "NOT_IN_ROLE",
						"User must be "+strings.Join(roles, " or ")+" to perform this operation").WithDetails("roles", roles).WithStatus(403))
			} else {
				next.ServeHTTP(res, req)
			}
		}
	}
}

func (c *RoleAuthorizer) UserInRole(role string) func(res http.ResponseWriter, req *http.Request, user *cdata.AnyValueMap, next http.HandlerFunc) {
	roles := make([]string, 1)
	roles[0] = role
	return c.UserInRoles(roles)
}

func (c *RoleAuthorizer) Admin() func(res http.ResponseWriter, req *http.Request, user *cdata.AnyValueMap, next http.HandlerFunc) {
	return c.UserInRole("admin")
}
