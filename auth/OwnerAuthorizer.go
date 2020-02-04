package auth

import (
	"encoding/base64"
	"net/http"
	"strings"

	cerr "github.com/pip-services3-go/pip-services3-commons-go/errors"
	services "github.com/pip-services3-go/pip-services3-rpc-go/services"
)

type OwnerAuthorizer struct {
}

func (c *OwnerAuthorizer) Owner(idParam string) func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if idParam == "" {
		idParam = "user_id"
	}
	return func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		auth := strings.SplitN(req.Header.Get("Authorization"), " ", 2)
		if len(auth) != 2 || auth[0] != "Basic" {
			http.Error(res, "authorization failed", http.StatusUnauthorized)
			return
		}
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if pair[0] == "" {
			services.HttpResponseSender.SendError(
				res, req,
				cerr.NewUnauthorizedError("", "NOT_SIGNED", "User must be signed in to perform this operation").WithStatus(401))
		} else {
			// params := req.URL.Query()
			// userId := params[idParam]

			// if req.user_id != userId {
			// 	services.HttpResponseSender.SendError(
			// 		res, req,
			// 		cerr.NewUnauthorizedError("", "FORBIDDEN",
			// 			"Only data owner can perform this operation").WithStatus(403))
			// } else {
			// 	next.ServeHTTP(res, req)
			// }
		}
	}
}

func (c *OwnerAuthorizer) ownerOrAdmin(idParam string) func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if idParam == "" {
		idParam = "user_id"
	}
	return func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

		auth := strings.SplitN(req.Header.Get("Authorization"), " ", 2)
		if len(auth) != 2 || auth[0] != "Basic" {
			http.Error(res, "authorization failed", http.StatusUnauthorized)
			return
		}
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if pair[0] == "" {
			services.HttpResponseSender.SendError(res, req,
				cerr.NewUnauthorizedError("", "NOT_SIGNED",
					"User must be signed in to perform this operation").WithStatus(401))
		} else {
			// /users/123/zones/456?user_id=123

			//  userId := req.params[idParam] || req.param(idParam);

			//  var roles interface{}
			//  if req.user != nil {
			// 	roles = req.user.roles
			//  }

			//  admin := _.includes(roles, "admin");
			// if req.user_id != userId && !admin {
			//     services.HttpResponseSender.SendError(
			//         req, res,
			//         cerr.NewUnauthorizedError("", "FORBIDDEN",
			//             "Only data owner can perform this operation").WithStatus(403));
			// } else {
			//     next();
			// }
		}
	}
}
