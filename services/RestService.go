package services

import (
	"net/http"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/v3/config"
	cerr "github.com/pip-services3-go/pip-services3-commons-go/v3/errors"
	crefer "github.com/pip-services3-go/pip-services3-commons-go/v3/refer"
	ccount "github.com/pip-services3-go/pip-services3-components-go/v3/count"
	clog "github.com/pip-services3-go/pip-services3-components-go/v3/log"
)

/*
Abstract service that receives remove calls via HTTP/REST protocol.

Configuration parameters:

- base_route:              base route for remote URI
- dependencies:
  - endpoint:              override for HTTP Endpoint dependency
  - controller:            override for Controller dependency
- connection(s):
  - discovery_key:         (optional) a key to retrieve the connection from connect.idiscovery.html IDiscovery
  - protocol:              connection protocol: http or https
  - host:                  host name or IP address
  - port:                  port number
  - uri:                   resource URI or connection string with all parameters in it
- credential - the HTTPS credentials:
  - ssl_key_file:         the SSL private key in PEM
  - ssl_crt_file:         the SSL certificate in PEM
  - ssl_ca_file:          the certificate authorities (root cerfiticates) in PEM

References:

- *:logger:*:*:1.0               (optional) ILogger components to pass log messages
- *:counters:*:*:1.0             (optional) ICounters components to pass collected measurements
- *:discovery:*:*:1.0            (optional) IDiscovery services to resolve connection
- *:endpoint:http:*:1.0          (optional) HttpEndpoint reference

See RestClient

 Example:

    class MyRestService extends RestService {
       private _controller: IMyController;
       ...
       func (c * RestService) constructor() {
          base();
          c.DependencyResolver.put(
              "controller",
              new Descriptor("mygroup","controller","*","*","1.0")
          );
       }

       func (c * RestService) setReferences(references: IReferences): void {
          base.setReferences(references);
          c._controller = c.DependencyResolver.getRequired<IMyController>("controller");
       }

       func (c * RestService) register(): void {
           registerRoute("get", "get_mydata", nil, (req, res) => {
               let correlationId = req.param("correlation_id");
               let id = req.param("id");
               c._controller.getMyData(correlationId, id, c.sendResult(req, res));
           });
           ...
       }
    }

    let service = new MyRestService();
    service.configure(ConfigParams.fromTuples(
        "connection.protocol", "http",
        "connection.host", "localhost",
        "connection.port", 8080
    ));
    service.setReferences(References.fromTuples(
       new Descriptor("mygroup","controller","default","default","1.0"), controller
    ));

    service.open("123", (err) => {
       console.log("The REST service is running on port 8080");
    });
*/
// implements IOpenable, IConfigurable, IReferenceable, IUnreferenceable, IRegisterable
type RestService struct {
	defaultConfig cconf.ConfigParams
	config        cconf.ConfigParams
	references    crefer.IReferences
	localEndpoint bool
	opened        bool
	//The base route.
	BaseRoute string
	//The HTTP endpoint that exposes this service.
	Endpoint HttpEndpoint
	//The dependency resolver.
	DependencyResolver DependencyResolver
	//The logger.
	Logger CompositeLogger
	//The performance counters.
	Counters CompositeCounters
}

// NewRestService is create new RestService
func NewRestService() *RestService {
	rs := RestService{}
	rs.defaultConfig = NewConfigParamsFromTuples(
		"base_route", "",
		"dependencies.endpoint", "*:endpoint:http:*:1.0",
	)
	rs.DependencyResolver = crefer.NewDependencyResolver(rs.defaultConfig)
	rs.Logger = clog.NewCompositeLogger()
	rs.Counters = ccount.NewCompositeCounters()
	return &rs
}

//Configures component by passing configuration parameters.
//- config    configuration parameters to be set.
func (c *RestService) Configure(config cconf.ConfigParams) {
	config = config.SetDefaults(RestService.defaultConfig)

	c.config = config
	c.DependencyResolver.configure(config)

	c.BaseRoute = config.getAsStringWithDefault("base_route", c.BaseRoute)
}

/*
	Sets references to dependent components.

	- references 	references to locate the component dependencies.
*/
func (c *RestService) SetReferences(references cref.IReferences) {
	c.references = references

	c.Logger.setReferences(references)
	c.Counters.setReferences(references)
	c.DependencyResolver.setReferences(references)

	// Get endpoint
	c.Endpoint = c.DependencyResolver.getOneOptional("endpoint")
	// Or create a local one
	if c.Endpoint == nil {
		c.Endpoint = c.createEndpoint()
		c.localEndpoint = true
	} else {
		c.localEndpoint = false
	}
	// Add registration callback to the endpoint
	c.Endpoint.register(c)
}

/*
	Unsets (clears) previously set references to dependent components.
*/
func (c *RestService) UnsetReferences() {
	// Remove registration callback from endpoint
	if c.Endpoint != nil {
		c.Endpoint.unregister(c)
		c.Endpoint = nil
	}
}

func (c *RestService) createEndpoint() HttpEndpoint {
	endpoint := NewHttpEndpoint()

	if c.config {
		endpoint.Configure(c.config)
	}
	if c.references {
		endpoint.SetReferences(c.references)
	}

	return endpoint
}

/*
   Adds instrumentation to log calls and measure call time.
   It returns a Timing object that is used to end the time measurement.

   - correlationId     (optional) transaction id to trace execution through call chain.
   - name              a method name.
   @returns Timing object to end the time measurement.
*/
func (c *RestService) Instrument(correlationId string, name string) ccount.Timing {
	c.Logger.Trace(correlationId, "Executing %s method", name)
	c.Counters.IncrementOne(name + ".exec_count")
	return c.Counters.BeginTiming(name + ".exec_time")
}

/*
   Adds instrumentation to error handling.

   - correlationId     (optional) transaction id to trace execution through call chain.
   - name              a method name.
   - err               an occured error
   - result            (optional) an execution result
   - callback          (optional) an execution callback
*/
func (c *RestService) InstrumentError(correlationId string, name string, err any,
	result interface{}) (result interface{}, err error) {
	if err != nil {
		c.Logger.error(correlationId, err, "Failed to execute %s method", name)
		c.Counters.incrementOne(name + ".exec_errors")
	}

	return result, err
}

/*
	Checks if the component is opened.

	@returns true if the component has been opened and false otherwise.
*/
func (c *RestService) IsOpen() bool {
	return c.opened
}

/*
	Opens the component.

	- correlationId 	(optional) transaction id to trace execution through call chain.
    - callback 			callback function that receives error or nil no errors occured.
*/
func (c *RestService) Open(correlationId string) error {
	if c.opened {
		return nil
	}

	if c.Endpoint == nil {
		c.Endpoint = c.createEndpoint()
		c.Endpoint.register(c)
		c.localEndpoint = true
	}

	if c.localEndpoint {
		oErr := c.Endpoint.Open(correlationId)
		if oErr != nil {
			c.opened = false
			return oErr
		}
	} else {
		c.opened = true
		return nil
	}
}

/*
	Closes component and frees used resources.

	- correlationId 	(optional) transaction id to trace execution through call chain.
    - callback 			callback function that receives error or nil no errors occured.
*/
func (c *RestService) Close(correlationId string) error {
	if !c.opened {
		return nil
	}

	if c.Endpoint == nil {
		return cerr.NewInvalidStateException(correlationId, "NO_ENDPOINT", "HTTP endpoint is missing")
	}

	if c.localEndpoint {
		cErr := c.Endpoint.Close(correlationId)
		if cErr != nil {
			c.opened = false
			return cErr
		}
	} else {
		c.opened = false
		return nil
	}
}

/*
   Creates a callback function that sends result as JSON object.
   That callack function call be called directly or passed
   as a parameter to business logic components.

   If object is not nil it returns 200 status code.
   For nil results it returns 204 status code.
   If error occur it sends ErrorDescription with approproate status code.

   - req       a HTTP request object.
   - res       a HTTP response object.
   - callback function that receives execution result or error.
*/
func (c *RestService) SendResult(req *http.Request, res http.ResponseWriter) (result interface{}, err error) {
	return HttpResponseSender.SendResult(req, res)
}

/*
   Creates a callback function that sends newly created object as JSON.
   That callack function call be called directly or passed
   as a parameter to business logic components.

   If object is not nil it returns 201 status code.
   For nil results it returns 204 status code.
   If error occur it sends ErrorDescription with approproate status code.

   - req       a HTTP request object.
   - res       a HTTP response object.
   - callback function that receives execution result or error.
*/
func (c *RestService) SendCreatedResult(req *http.Request, res http.ResponseWriter) (result interface{}, err errro) {
	return HttpResponseSender.sendCreatedResult(req, res)
}

/*
   Creates a callback function that sends deleted object as JSON.
   That callack function call be called directly or passed
   as a parameter to business logic components.

   If object is not nil it returns 200 status code.
   For nil results it returns 204 status code.
   If error occur it sends ErrorDescription with approproate status code.

   - req       a HTTP request object.
   - res       a HTTP response object.
   - callback function that receives execution result or error.
*/
func (c *RestService) SendDeletedResult(req *http.Request, res http.ResponseWriter) (result interface{}, err error) {
	return HttpResponseSender.sendDeletedResult(req, res)
}

/*
   Sends error serialized as ErrorDescription object
   and appropriate HTTP status code.
   If status code is not defined, it uses 500 status code.

   - req       a HTTP request object.
   - res       a HTTP response object.
   - error     an error object to be sent.
*/
func (c *RestService) SendError(req *http.Request, res http.ResponseWriter, err error) {
	HttpResponseSender.SendError(req, res, err)
}

func (c *RestService) appendBaseRoute(route string) string {
	route = route || ""

	if c.BaseRoute != nil && c.BaseRoute.length > 0 {
		baseRoute := c.BaseRoute
		if baseRoute[0] != "/"[0] {
			baseRoute = "/" + baseRoute
		}
		route = baseRoute + route
	}
	return route
}

/*
   Registers a route in HTTP endpoint.

   - method        HTTP method: "get", "head", "post", "put", "delete"
   - route         a command route. Base route will be added to this route
   - schema        a validation schema to validate received parameters.
   - action        an action function that is called when operation is invoked.
*/
func (c *RestService) RegisterRoute(method string, route string, schema Schema,
	action func(req *http.Request, res http.ResponseWriter)) {
	if c.Endpoint == nil {
		return
	}

	route = c.appendBaseRoute(route)

	c.Endpoint.registerRoute(
		method, route, schema, func(req *http.Request, res http.ResponseWriter) {
			action.call(c, req, res)
		})
}

/*
   Registers a route with authorization in HTTP endpoint.

   - method        HTTP method: "get", "head", "post", "put", "delete"
   - route         a command route. Base route will be added to this route
   - schema        a validation schema to validate received parameters.
   - authorize     an authorization interceptor
   - action        an action function that is called when operation is invoked.
*/
func (c *RestService) RegisterRouteWithAuth(method string, route string, schema Schema,
	authorize func(req *http.Request, res http.ResponseWriter, next func()),
	action func(req *http.Request, res http.ResponseWriter)) {
	if c.Endpoint == nil {
		return
	}

	route = c.appendBaseRoute(route)

	c.Endpoint.registerRouteWithAuth(
		method, route, schema,
		func(req *http.Request, res http.ResponseWriter, next func()) {
			if authorize {
				authorize.call(c, req, res, next)
			} else {
				next()
			}
		}, func(req *http.Request, res http.ResponseWriter) {
			action.call(c, req, res)
		})
}

/*
   Registers a middleware for a given route in HTTP endpoint.

   - route         a command route. Base route will be added to this route
   - action        an action function that is called when middleware is invoked.
*/
func (c *RestService) RegisterInterceptor(route string,
	action func(req *http.Request, res http.ResponseWriter, next func())) {
	if c.Endpoint == nil {
		return
	}

	route = c.appendBaseRoute(route)

	c.Endpoint.registerInterceptor(
		route, func(req *http.Request, res http.ResponseWriter, next func()) {
			action.call(c, req, res, next)
		})
}

/*
   Registers all service routes in HTTP endpoint.

   This method is called by the service and must be overriden
   in child classes.
*/
func (c *RestService) Register() {}
