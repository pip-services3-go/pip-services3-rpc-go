package services

import (
	cerr "github.com/pip-services3-go/pip-services3-commons-go/v3/errors"
	"net/http"
)

/*
Helper class that handles HTTP-based responses.
*/
type THttpResponseSender struct {
}

var HttpResponseSender THttpResponseSender

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

	// res.Status(result.status)
	// res.JSON(result)
}

/**
 * Creates a callback function that sends result as JSON object.
 * That callack function call be called directly or passed
 * as a parameter to business logic components.
 *
 * If object is not nil it returns 200 status code.
 * For nil results it returns 204 status code.
 * If error occur it sends ErrorDescription with approproate status code.
 *
 * - req       a HTTP request object.
 * - res       a HTTP response object.
 * - callback function that receives execution result or error.
 */
func (c *THttpResponseSender) SendResult(req *http.Request, res http.ResponseWriter) func(result interface{}, err error) {
	return func(result interface{}, err error) {
		if err != nil {
			HttpResponseSender.SendError(req, res, err)
			return
		}
		if result == nil {
			//res.Send(204)
		} else {
			//res.Json(result)
		}
	}
}

/**
 * Creates a callback function that sends an empty result with 204 status code.
 * If error occur it sends ErrorDescription with approproate status code.
 *
 * - req       a HTTP request object.
 * - res       a HTTP response object.
 * - callback function that receives error or nil for success.
 */
func (c *THttpResponseSender) SendEmptyResult(req *http.Request, res http.ResponseWriter) func(err error) {
	return func(err error) {
		if err != nil {
			HttpResponseSender.SendError(req, res, err)
			return
		}
		//res.Send(204)
	}
}

/**
 * Creates a callback function that sends newly created object as JSON.
 * That callack function call be called directly or passed
 * as a parameter to business logic components.
 *
 * If object is not nil it returns 201 status code.
 * For nil results it returns 204 status code.
 * If error occur it sends ErrorDescription with approproate status code.
 *
 * - req       a HTTP request object.
 * - res       a HTTP response object.
 * - callback function that receives execution result or error.
 */
func (c *THttpResponseSender) SendCreatedResult(req *http.Request, res http.ResponseWriter) func(result interface{}, err error) {
	return func(result interface{}, err error) {
		if err != nil {
			HttpResponseSender.SendError(req, res, err)
			return
		}
		if result == nil {
			//res.Status(204)
		} else {
			//res.Status(201)
			//res.Json(result)
		}
	}
}

/**
 * Creates a callback function that sends deleted object as JSON.
 * That callack function call be called directly or passed
 * as a parameter to business logic components.
 *
 * If object is not nil it returns 200 status code.
 * For nil results it returns 204 status code.
 * If error occur it sends ErrorDescription with approproate status code.
 *
 * - req       a HTTP request object.
 * - res       a HTTP response object.
 * - callback function that receives execution result or error.
 */
func (c *THttpResponseSender) SendDeletedResult(req *http.Request, res http.ResponseWriter) func(result interface{}, err error) {
	return func(result interface{}, err error) {
		if err != nil {
			HttpResponseSender.SendError(req, res, err)
			return
		}
		if result == nil {
			//res.Status(204)
		} else {
			//res.Status(200)
			//res.Json(result)
		}
	}
}
