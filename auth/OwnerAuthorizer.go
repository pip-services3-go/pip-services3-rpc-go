package auth

import (
	"encoding/base64"
	cerr "github.com/pip-services3-go/pip-services3-commons-go/errors"
	services "github.com/pip-services3-go/pip-services3-rpc-go/services"
	"net/http"
	"strings"
)

type OwnerAuthorizer struct {
}

func (c *OwnerAuthorizer) Owner(idParam string) func(req *http.Request, res http.ResponseWriter, next func()) {
	if idParam == "" {
		idParam = "user_id"
	}
	return func(req *http.Request, res http.ResponseWriter, next func()) {
		auth := strings.SplitN(req.Header.Get("Authorization"), " ", 2)
		if len(auth) != 2 || auth[0] != "Basic" {
			http.Error(res, "authorization failed", http.StatusUnauthorized)
			return
		}
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if pair[0] == "" {
			services.HttpResponseSender.SendError(
				req, res,
				cerr.NewUnauthorizedError("", "NOT_SIGNED", "User must be signed in to perform this operation").WithStatus(401))
		} else {
			// userId := req.params[idParam] || req.param(idParam);
			// if req.user_id != userId {
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

func (c *OwnerAuthorizer) ownerOrAdmin(idParam string) func(req *http.Request, res http.ResponseWriter, next func()) {
	if idParam == "" {
		idParam = "user_id"
	}
	return func(req *http.Request, res http.ResponseWriter, next func()) {

		auth := strings.SplitN(req.Header.Get("Authorization"), " ", 2)
		if len(auth) != 2 || auth[0] != "Basic" {
			http.Error(res, "authorization failed", http.StatusUnauthorized)
			return
		}
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if pair[0] == "" {
			services.HttpResponseSender.SendError(req, res,
				cerr.NewUnauthorizedError("", "NOT_SIGNED",
					"User must be signed in to perform this operation").WithStatus(401))
		} else {
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
