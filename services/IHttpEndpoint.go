package services

import (
	"net/http"

	crun "github.com/pip-services3-go/pip-services3-commons-go/run"
	cvalid "github.com/pip-services3-go/pip-services3-commons-go/validate"
)

// IHttpEndpoint interface for custom implementation of endpoints
// Optional configurable and refererences interfaces can be implemented
type IHttpEndpoint interface {
	ITlsConfigurator
	crun.IOpenable

	Register(registration IRegisterable)
	Unregister(registration IRegisterable)
	GetCorrelationId(req *http.Request) string
	RegisterRoute(method string, route string, schema *cvalid.Schema, action http.HandlerFunc)
	RegisterRouteWithAuth(method string, route string, schema *cvalid.Schema,
		authorize func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc),
		action http.HandlerFunc)
	RegisterInterceptor(route string, action func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc))
	AddCorsHeader(header string, origin string)
}
