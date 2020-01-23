package auth

import (
	"encoding/base64"
	cerr "github.com/pip-services3-go/pip-services3-commons-go/errors"
	services "github.com/pip-services3-go/pip-services3-rpc-go/services"
	"net/http"
	"strings"
)

// import { HttpResponseSender } from "../services/HttpResponseSender";

type BasicAuthorizer struct {
}

func (c *BasicAuthorizer) Anybody() func(req *http.Request, res http.ResponseWriter, next func()) {
	return func(req *http.Request, res http.ResponseWriter, next func()) {
		next()
	}
}

func (c *BasicAuthorizer) Signed() func(req *http.Request, res http.ResponseWriter, next func()) {
	return func(req *http.Request, res http.ResponseWriter, next func()) {

		auth := strings.SplitN(req.Header.Get("Authorization"), " ", 2)
		if len(auth) != 2 || auth[0] != "Basic" {
			http.Error(res, "authorization failed", http.StatusUnauthorized)
			return
		}
		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if pair[0] == "" { // username
			services.HttpResponseSender.SendError(
				req, res,
				cerr.NewUnauthorizedError("", "NOT_SIGNED", "User must be signed in to perform this operation").WithStatus(401))
		} else {
			next()
		}
	}
}
