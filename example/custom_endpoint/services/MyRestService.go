package services

import (
	"encoding/json"
	"io"
	"net/http"

	elogic "github.com/pip-services3-go/pip-services3-rpc-go/example/custom_endpoint/logic"

	edata "github.com/pip-services3-go/pip-services3-rpc-go/example/custom_endpoint/data/version1"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cerr "github.com/pip-services3-go/pip-services3-commons-go/errors"
	crefer "github.com/pip-services3-go/pip-services3-commons-go/refer"
	cvalid "github.com/pip-services3-go/pip-services3-commons-go/validate"
	"github.com/pip-services3-go/pip-services3-rpc-go/services"
)

type MyRestService struct {
	*services.RestService
	controller elogic.IMyController
}

func NewMyRestService() *MyRestService {
	c := &MyRestService{}
	c.RestService = services.InheritRestService(c)
	c.DependencyResolver.Put("controller", crefer.NewDescriptor("pip-services-dummies", "controller", "default", "*", "*"))
	return c
}

func (c *MyRestService) Configure(config *cconf.ConfigParams) {
	c.RestService.Configure(config)
}

func (c *MyRestService) SetReferences(references crefer.IReferences) {
	c.RestService.SetReferences(references)
	depRes, depErr := c.DependencyResolver.GetOneRequired("controller")
	if depErr == nil && depRes != nil {
		c.controller = depRes.(elogic.IMyController)
	}
}

func (c *MyRestService) hello(res http.ResponseWriter, req *http.Request) {
	correlationId := c.GetCorrelationId(req)
	var hello edata.HelloV1

	body, bodyErr := io.ReadAll(req.Body)
	if bodyErr != nil {
		err := cerr.NewInternalError(correlationId, "JSON_CNV_ERR", "Cant convert from JSON to Dummy").WithCause(bodyErr)
		c.SendError(res, req, err)
		return
	}
	defer req.Body.Close()
	jsonErr := json.Unmarshal(body, &hello)

	if jsonErr != nil {
		err := cerr.NewInternalError(correlationId, "JSON_CNV_ERR", "Cant convert from JSON to Dummy").WithCause(jsonErr)
		c.SendError(res, req, err)
		return
	}

	result, err := c.controller.SayHello(
		correlationId,
		hello.Name,
	)
	c.SendResult(res, req, result, err)
}

func (c *MyRestService) Register() {
	c.RegisterRoute(
		"post", "/hello",
		&cvalid.NewObjectSchema().
			WithRequiredProperty("body", edata.NewHelloV1Schema()).Schema,
		c.hello,
	)
}
