package services

import (
	"encoding/json"
	cerr "github.com/pip-services3-go/pip-services3-commons-go/errors"
	"io"
	"net/http"
)

/*
Helper class that handles HTTP-based responses.
*/
type THttpResponseSender struct {
}

var HttpResponseSender THttpResponseSender = THttpResponseSender{}

/*
   Sends error serialized as ErrorDescription object
   and appropriate HTTP status code.
   If status code is not defined, it uses 500 status code.

   - req       a HTTP request object.
   - res       a HTTP response object.
   - error     an error object to be sent.
*/
func (c *THttpResponseSender) SendError(req *http.Request, res http.ResponseWriter, err error) {
	// if err == nil {
	// 	err = Error{}
	// }
	appErr := cerr.ApplicationError{}
	err = appErr.Wrap(err)

	// result = _.pick(error, "code", "status", "name", "details", "component", "message", "stack", "cause");
	// result = _.defaults(result, { code: "Undefined", status: 500, message: "Unknown error" });

	result := make(map[string]string, 8)
	result["code"] = "Undefined"
	result["status"] = "500"
	result["name"] = ""
	result["details"] = ""
	result["component"] = ""
	result["message"] = "Unknown error"
	result["stack"] = ""
	result["cause"] = ""

	res.WriteHeader(500)
	// res.JSON(result)
	jsonObj, jsonErr := json.Marshal(result)
	if jsonErr == nil {
		io.WriteString(res, (string)(jsonObj))
	}
}

/*
Creates a callback function that sends result as JSON object.
That callack function call be called directly or passed
as a parameter to business logic components.
 *
If object is not nil it returns 200 status code.
For nil results it returns 204 status code.
If error occur it sends ErrorDescription with approproate status code.
 *
- req       a HTTP request object.
- res       a HTTP response object.
- callback function that receives execution result or error.
*/
func (c *THttpResponseSender) SendResult(req *http.Request, res http.ResponseWriter) func(result interface{}, err error) {
	return func(result interface{}, err error) {
		if err != nil {
			HttpResponseSender.SendError(req, res, err)
			return
		}
		if result == nil {
			res.WriteHeader(204)
		} else {
			jsonObj, jsonErr := json.Marshal(result)
			if jsonErr == nil {
				io.WriteString(res, (string)(jsonObj))
			}
		}
	}
}

/*
Creates a callback function that sends an empty result with 204 status code.
If error occur it sends ErrorDescription with approproate status code.
 *
- req       a HTTP request object.
- res       a HTTP response object.
- callback function that receives error or nil for success.
*/
func (c *THttpResponseSender) SendEmptyResult(req *http.Request, res http.ResponseWriter) func(err error) {
	return func(err error) {
		if err != nil {
			HttpResponseSender.SendError(req, res, err)
			return
		}
		res.WriteHeader(204)
	}
}

/*
Creates a callback function that sends newly created object as JSON.
That callack function call be called directly or passed
as a parameter to business logic components.
 *
If object is not nil it returns 201 status code.
For nil results it returns 204 status code.
If error occur it sends ErrorDescription with approproate status code.
 *
- req       a HTTP request object.
- res       a HTTP response object.
- callback function that receives execution result or error.
*/
func (c *THttpResponseSender) SendCreatedResult(req *http.Request, res http.ResponseWriter) func(result interface{}, err error) {
	return func(result interface{}, err error) {
		if err != nil {
			HttpResponseSender.SendError(req, res, err)
			return
		}
		if result == nil {
			res.WriteHeader(204)
		} else {
			res.WriteHeader(201)
			//res.Json(result)
			jsonObj, jsonErr := json.Marshal(result)
			if jsonErr == nil {
				io.WriteString(res, (string)(jsonObj))
			}
		}
	}
}

/*
Creates a callback function that sends deleted object as JSON.
That callack function call be called directly or passed
as a parameter to business logic components.
 *
If object is not nil it returns 200 status code.
For nil results it returns 204 status code.
If error occur it sends ErrorDescription with approproate status code.
 *
- req       a HTTP request object.
- res       a HTTP response object.
- callback function that receives execution result or error.
*/
func (c *THttpResponseSender) SendDeletedResult(req *http.Request, res http.ResponseWriter) func(result interface{}, err error) {
	return func(result interface{}, err error) {
		if err != nil {
			HttpResponseSender.SendError(req, res, err)
			return
		}
		if result == nil {
			res.WriteHeader(204)
		} else {
			res.WriteHeader(200)
			//res.Json(result)
			jsonObj, jsonErr := json.Marshal(result)
			if jsonErr == nil {
				io.WriteString(res, (string)(jsonObj))
			}
		}
	}
}
