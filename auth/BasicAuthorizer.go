package auth

import (
	"encoding/base64"
	"net/http"
	"strings"

	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	cerr "github.com/pip-services3-go/pip-services3-commons-go/errors"
	services "github.com/pip-services3-go/pip-services3-rpc-go/services"
)

type BasicAuthorizer struct {
}

func (c *BasicAuthorizer) Anybody() func(res http.ResponseWriter, req *http.Request, user *cdata.AnyValueMap, next http.HandlerFunc) {
	return func(res http.ResponseWriter, req *http.Request, user *cdata.AnyValueMap, next http.HandlerFunc) {
		next.ServeHTTP(res, req)
	}
}

func (c *BasicAuthorizer) Signed() func(res http.ResponseWriter, req *http.Request, user *cdata.AnyValueMap, next http.HandlerFunc) {
	return func(res http.ResponseWriter, req *http.Request, user *cdata.AnyValueMap, next http.HandlerFunc) {

		auth := strings.SplitN(req.Header.Get("Authorization"), " ", 2)
		if len(auth) != 2 || auth[0] != "Basic" {
			http.Error(res, "authorization failed", http.StatusUnauthorized)
			return
		}
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if pair[0] == "" { // username
			services.HttpResponseSender.SendError(
				res, req,
				cerr.NewUnauthorizedError("", "NOT_SIGNED", "User must be signed in to perform this operation").WithStatus(401))
		} else {
			next.ServeHTTP(res, req)
		}
	}
}
