package services

import (
	"net/http"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	cerr "github.com/pip-services3-go/pip-services3-commons-go/errors"
	crefer "github.com/pip-services3-go/pip-services3-commons-go/refer"
	ccount "github.com/pip-services3-go/pip-services3-components-go/count"
	clog "github.com/pip-services3-go/pip-services3-components-go/log"
)

type RestOperations struct {
	Logger             *clog.CompositeLogger
	Counters           *ccount.CompositeCounters
	DependencyResolver *crefer.DependencyResolver
}

func NewRestOperations() *RestOperations {
	ro := RestOperations{}
	ro.Logger = clog.NewCompositeLogger()
	ro.Counters = ccount.NewCompositeCounters()
	ro.DependencyResolver = crefer.NewDependencyResolver()
	return &ro
}

func (c *RestOperations) Configure(config *cconf.ConfigParams) {
	c.DependencyResolver.Configure(config)
}

func (c *RestOperations) SetReferences(references crefer.IReferences) {
	c.Logger.SetReferences(references)
	c.Counters.SetReferences(references)
	c.DependencyResolver.SetReferences(references)
}

func (c *RestOperations) GetCorrelationId(req *http.Request) string {
	params := req.URL.Query()
	return params.Get("correlation_id")
}

func (c *RestOperations) GetFilterParams(req *http.Request) *cdata.FilterParams {

	params := req.URL.Query()
	delete(params, "skip")
	delete(params, "take")
	delete(params, "total")
	filter := cdata.NewFilterParamsFromValue(
		params,
	)
	return filter
}

func (c *RestOperations) GetPagingParams(req *http.Request) *cdata.PagingParams {

	params := req.URL.Query()
	paginParams := make(map[string]string, 0)

	paginParams["skip"] = params.Get("skip")
	paginParams["take"] = params.Get("take")
	paginParams["total"] = params.Get("total")

	paging := cdata.NewPagingParamsFromValue(
		paginParams,
	)
	return paging
}

func (c *RestOperations) SendResult(res http.ResponseWriter, req *http.Request, result interface{}, err error) {
	HttpResponseSender.SendResult(res, req, result, err)
}

func (c *RestOperations) SendEmptyResult(res http.ResponseWriter, req *http.Request, err error) {
	HttpResponseSender.SendEmptyResult(res, req, err)
}

func (c *RestOperations) SendCreatedResult(res http.ResponseWriter, req *http.Request, result interface{}, err error) {
	HttpResponseSender.SendCreatedResult(res, req, result, err)
}

func (c *RestOperations) SendDeletedResult(res http.ResponseWriter, req *http.Request, result interface{}, err error) {
	HttpResponseSender.SendDeletedResult(res, req, result, err)
}

func (c *RestOperations) SendError(res http.ResponseWriter, req *http.Request, err error) {
	HttpResponseSender.SendError(res, req, err)
}

func (c *RestOperations) SendBadRequest(res http.ResponseWriter, req *http.Request, message string) {
	correlationId := c.GetCorrelationId(req)
	error := cerr.NewBadRequestError(correlationId, "BAD_REQUEST", message)
	c.SendError(res, req, error)
}

func (c *RestOperations) SendUnauthorized(res http.ResponseWriter, req *http.Request, message string) {
	correlationId := c.GetCorrelationId(req)
	error := cerr.NewUnauthorizedError(correlationId, "UNAUTHORIZED", message)
	c.SendError(res, req, error)
}

func (c *RestOperations) SendNotFound(res http.ResponseWriter, req *http.Request, message string) {
	correlationId := c.GetCorrelationId(req)
	error := cerr.NewNotFoundError(correlationId, "NOT_FOUND", message)
	c.SendError(res, req, error)
}

func (c *RestOperations) SendConflict(res http.ResponseWriter, req *http.Request, message string) {
	correlationId := c.GetCorrelationId(req)
	error := cerr.NewConflictError(correlationId, "CONFLICT", message)
	c.SendError(res, req, error)
}

func (c *RestOperations) SendSessionExpired(res http.ResponseWriter, req *http.Request, message string) {
	correlationId := c.GetCorrelationId(req)
	err := cerr.NewUnknownError(correlationId, "SESSION_EXPIRED", message)
	err.Status = 440
	c.SendError(res, req, err)
}

func (c *RestOperations) SendInternalError(res http.ResponseWriter, req *http.Request, message string) {
	correlationId := c.GetCorrelationId(req)
	error := cerr.NewUnknownError(correlationId, "INTERNAL", message)
	c.SendError(res, req, error)
}

func (c *RestOperations) SendServerUnavailable(res http.ResponseWriter, req *http.Request, message string) {
	correlationId := c.GetCorrelationId(req)
	err := cerr.NewConflictError(correlationId, "SERVER_UNAVAILABLE", message)
	err.Status = 503
	c.SendError(res, req, err)
}

// func (c *RestOperations) Invoke(operation string) func(res http.ResponseWriter, req *http.Request) {
// 	return func(res http.ResponseWriter, req *http.Request) {
// 		// TODO: what is it
// 		//c[operation](res http.ResponseWriter, req *http.Request);
// 	}
// }
