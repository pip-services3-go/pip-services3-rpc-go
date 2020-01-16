package clients

import (
	"strings"
	"time"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/v3/config"
	cdata "github.com/pip-services3-go/pip-services3-commons-go/v3/data"
	cerr "github.com/pip-services3-go/pip-services3-commons-go/v3/errors"
	crefer "github.com/pip-services3-go/pip-services3-commons-go/v3/refer"
	ccount "github.com/pip-services3-go/pip-services3-components-go/v3/count"
	clog "github.com/pip-services3-go/pip-services3-components-go/v3/log"
	rpccon "github.com/pip-services3-go/pip-services3-rpc-go/v3/connect"
)

// import { HttpConnectionResolver } from "../connect/HttpConnectionResolver";

/*
Abstract client that calls remove endpoints using HTTP/REST protocol.

Configuration parameters:

- base_route:              base route for remote URI
- connection(s):
  - discovery_key:         (optional) a key to retrieve the connection from [[https://rawgit.com/pip-services-node/pip-services3-components-node/master/doc/api/interfaces/connect.idiscovery.html IDiscovery]]
  - protocol:              connection protocol: http or https
  - host:                  host name or IP address
  - port:                  port number
  - uri:                   resource URI or connection string with all parameters in it
- options:
  - retries:               number of retries (default: 3)
  - connectTimeout:       connection timeout in milliseconds (default: 10 sec)
  - timeout:               invocation timeout in milliseconds (default: 10 sec)

References:

- <code>\*:logger:\*:\*:1.0</code>         (optional) [[https://rawgit.com/pip-services-node/pip-services3-components-node/master/doc/api/interfaces/log.ilogger.html ILogger]] components to pass log messages
- <code>\*:counters:\*:\*:1.0</code>         (optional) [[https://rawgit.com/pip-services-node/pip-services3-components-node/master/doc/api/interfaces/count.icounters.html ICounters]] components to pass collected measurements
- <code>\*:discovery:\*:\*:1.0</code>        (optional) [[https://rawgit.com/pip-services-node/pip-services3-components-node/master/doc/api/interfaces/connect.idiscovery.html IDiscovery]] services to resolve connection
 *
See [[RestService]]
See [[CommandableHttpService]]
 *
 Example
 *
    class MyRestClient extends RestClient implements IMyClient {
       ...
 *
       func (c* RestClient) getData(correlationId: string, id: string,
           callback: (err: interface{}, result: MyData) => void): void {
 *
           let timing = c.instrument(correlationId, "myclient.get_data");
           c.call("get", "/get_data" correlationId, { id: id }, nil, (err, result) => {
               timing.endTiming();
               callback(err, result);
           });
       }
       ...
    }
 *
    let client = new MyRestClient();
    client.configure(ConfigParams.fromTuples(
        "connection.protocol", "http",
        "connection.host", "localhost",
        "connection.port", 8080
    ));
 *
    client.getData("123", "1", (err, result) => {
      ...
    });
*/
// implements IOpenable, IConfigurable, IReferenceable
type RestClient struct {
	defaultConfig cconf.ConfigParams
	//The HTTP client.
	Client *JsonHttpClient
	//The connection resolver.
	ConnectionResolver rpccon.HttpConnectionResolver
	//The logger.
	Logger clog.CompositeLogger
	//The performance counters.
	Counters ccount.CompositeCounters
	//The configuration options.
	Options cconf.ConfigParams
	//The base route.
	BaseRoute string
	//The number of retries.
	Retries int
	//The default headers to be added to every request.
	Headers cdata.StringValueMap
	//The connection timeout in milliseconds.
	ConnectTimeout int
	//The invocation timeout in milliseconds.
	Timeout int
	//The remote service uri which is calculated on open.
	Uri string
}

func NewRestClient() *RestClient {
	rc := RestClient{}
	rc.defaultConfig = *cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "0.0.0.0",
		"connection.port", 3000,

		"options.request_max_size", 1024*1024,
		"options.connectTimeout", 10000,
		"options.timeout", 10000,
		"options.retries", 3,
		"options.debug", true,
	)
	rc.ConnectionResolver = *rpccon.NewHttpConnectionResolver()
	rc.Logger = *clog.NewCompositeLogger()
	rc.Counters = *ccount.NewCompositeCounters()
	rc.Options = *cconf.NewEmptyConfigParams()
	rc.Retries = 1
	rc.ConnectTimeout = 10000
	return &rc
}

/*
   Configures component by passing configuration parameters.
    *
   @param config    configuration parameters to be set.
*/
func (c *RestClient) Configure(config *cconf.ConfigParams) {
	config = config.SetDefaults(&c.defaultConfig)
	c.ConnectionResolver.Configure(config)
	c.Options = *c.Options.Override(config.GetSection("options"))
	c.Retries = config.GetAsIntegerWithDefault("options.retries", c.Retries)
	c.ConnectTimeout = config.GetAsIntegerWithDefault("options.connectTimeout", c.ConnectTimeout)
	c.Timeout = config.GetAsIntegerWithDefault("options.timeout", c.Timeout)

	c.BaseRoute = config.GetAsStringWithDefault("base_route", c.BaseRoute)
}

/*
	Sets references to dependent components.
	 *
	@param references 	references to locate the component dependencies.
*/
func (c *RestClient) SetReferences(references crefer.IReferences) {
	c.Logger.SetReferences(references)
	c.Counters.SetReferences(references)
	c.ConnectionResolver.SetReferences(references)
}

/*
   Adds instrumentation to log calls and measure call time.
   It returns a Timing object that is used to end the time measurement.
    *
   @param correlationId     (optional) transaction id to trace execution through call chain.
   @param name              a method name.
   @returns Timing object to end the time measurement.
*/
func (c *RestClient) Instrument(correlationId string, name string) *ccount.Timing {
	c.Logger.Trace(correlationId, "Calling %s method", name)
	c.Counters.IncrementOne(name + ".call_count")
	return c.Counters.BeginTiming(name + ".call_time")
}

/*
   Adds instrumentation to error handling.
    *
   @param correlationId     (optional) transaction id to trace execution through call chain.
   @param name              a method name.
   @param err               an occured error
   @param result            (optional) an execution result
   @param callback          (optional) an execution callback
*/
func (c *RestClient) instrumentError(correlationId string, name string, inErr error, inRes interface{}) (result interface{}, err error) {
	if inErr != nil {
		c.Logger.Error(correlationId, inErr, "Failed to call %s method", name)
		c.Counters.IncrementOne(name + ".call_errors")
	}

	return inRes, inErr
}

/*
	Checks if the component is opened.
	 *
	@returns true if the component has been opened and false otherwise.
*/
func (c *RestClient) IsOpen() bool {
	return c.Client != nil
}

/*
	Opens the component.
	 *
	@param correlationId 	(optional) transaction id to trace execution through call chain.
    @param callback 			callback function that receives error or nil no errors occured.
*/
func (c *RestClient) Open(correlationId string) error {
	if c.IsOpen() {
		return nil
	}

	connection, _, conErr := c.ConnectionResolver.Resolve(correlationId)
	if conErr != nil {
		return conErr
	}

	//try {
	c.Uri = connection.Uri()
	c.Client = NewJsonHttpClient()
	if c.Client == nil {
		//c.Client = nil
		//ex := cerr.NewConnectionError(correlationId, "CANNOT_CONNECT", "Connection to REST service failed").Wrap(err).WithDetails("url", c.Uri)
		ex := cerr.NewConnectionError(correlationId, "CANNOT_CONNECT", "Connection to REST service failed").WithDetails("url", c.Uri)
		return ex
	}
	c.Client.SetUrl(c.Uri)
	c.Client.RetryWaitMin = (time.Duration)(c.Timeout) * time.Millisecond
	c.Client.RetryMax = c.Retries
	c.Client.Client.HTTPClient.Timeout = (time.Duration)(c.ConnectTimeout) * time.Millisecond
	c.Client.SetHeaders(&c.Headers)
	// let restify = require("restify-clients");
	// c.Client = restify.createJsonClient({
	//     url: c.Uri,
	//     connectTimeout: c.ConnectTimeout,
	//     requestTimeout: c.Timeout,
	//     headers: c.Headers,
	//     retry: {
	//         minTimeout: c.Timeout,
	//         maxTimeout: Infinity,
	//         retries: c.Retries
	//     },
	//     version: "*"
	// });

	return nil
	// } catch (err) {

	//}
}

/*
	Closes component and frees used resources.
	 *
	@param correlationId 	(optional) transaction id to trace execution through call chain.
    @param callback 			callback function that receives error or nil no errors occured.
*/
func (c *RestClient) Close(correlationId string) error {
	if c.Client != nil {
		// Eat exceptions
		//try {
		c.Logger.Debug(correlationId, "Closed REST service at %s", c.Uri)
		//} catch (ex) {
		// c.Logger.Warn(correlationId, "Failed while closing REST service: %s", ex);
		//}
		c.Client = nil
		c.Uri = ""
	}
	return nil
}

/*
   Adds a correlation id (correlation_id) to invocation parameter map.
    *
   @param params            invocation parameters.
   @param correlationId     (optional) a correlation id to be added.
   @returns invocation parameters with added correlation id.
*/
func (c *RestClient) AddCorrelationId(params *cdata.StringValueMap, correlationId string) *cdata.StringValueMap {
	// Automatically generate short ids for now
	if correlationId == "" {
		//correlationId = IdGenerator.NextShort()
		return params
	}

	if params == nil {
		params = cdata.NewEmptyStringValueMap()
	}
	params.Put("correlation_id", correlationId)
	return params
}

/*
   Adds filter parameters (with the same name as they defined)
   to invocation parameter map.
    *
   @param params        invocation parameters.
   @param filter        (optional) filter parameters
   @returns invocation parameters with added filter parameters.
*/
func (c *RestClient) AddFilterParams(params *cdata.StringValueMap, filter *cdata.FilterParams) *cdata.StringValueMap {

	if params == nil {
		params = cdata.NewEmptyStringValueMap()
	}
	if filter != nil {
		for k, v := range filter.Value() {
			params.Put(k, v)
		}
		// for (let prop in filter) {
		//     if (filter.HasOwnProperty(prop))
		//         params[prop] = filter[prop];
		// }
	}
	return params
}

/*
   Adds paging parameters (skip, take, total) to invocation parameter map.
    *
   @param params        invocation parameters.
   @param paging        (optional) paging parameters
   @returns invocation parameters with added paging parameters.
*/
func (c *RestClient) AddPagingParams(params *cdata.StringValueMap, paging *cdata.PagingParams) *cdata.StringValueMap {
	if params == nil {
		params = cdata.NewEmptyStringValueMap()
	}

	if paging != nil {
		params.Put("total", paging.Total)
		if paging.Skip != nil {
			params.Put("skip", *paging.Skip)
		}
		if paging.Take != nil {
			params.Put("take", *paging.Take)
		}
	}

	return params
}

func (c *RestClient) createRequestRoute(route string) string {
	builder := ""

	if c.BaseRoute != "" && len(c.BaseRoute) > 0 {
		if c.BaseRoute[0] != "/"[0] {
			builder += "/"
		}
		builder += c.BaseRoute
	}

	if route[0] != "/"[0] {
		builder += "/"
	}
	builder += route

	return builder
}

/*
   Calls a remote method via HTTP/REST protocol.
    *
   @param method            HTTP method: "get", "head", "post", "put", "delete"
   @param route             a command route. Base route will be added to c route
   @param correlationId     (optional) transaction id to trace execution through call chain.
   @param params            (optional) query parameters.
   @param data              (optional) body object.
   @param callback          (optional) callback function that receives result object or error.
*/
func (c *RestClient) Call(method string, route string, correlationId string, params *cdata.StringValueMap, data interface{}) (result interface{}, err error) {

	method = strings.ToLower(method)

	// if _.isFunction(data) {
	// 	callback = data
	// 	//data = {};
	// }

	// route = c.createRequestRoute(route)
	// params = c.AddCorrelationId(params, correlationId)
	// if !_.isEmpty(params) {
	// 	route += "?" + querystring.stringify(params)
	// }

	// self := c
	// action := nil
	// if callback == nil {
	// 	action = func(err error, req *http.Request, res http.Response, data interface{}) {
	// 		// Handling 204 codes
	// 		if res && res.StatusCode == 204 {
	// 			callback.Call(self, nil, nil)
	// 		} else if err == nil {
	// 			callback.Call(self, nil, data)
	// 		} else {
	// 			// Restore application exception
	// 			if data != nil {
	// 				err = cerr.ApplicationErrorFactory.Create(data).WithCause(err)
	// 			}
	// 			callback.Call(self, err, nil)
	// 		}
	// 	}
	// }

	// if method == "get" {
	// 	c.Client.Get(route, action)
	// } else if method == "head" {
	// 	c.Client.Head(route, action)
	// } else if method == "post" {
	// 	c.Client.Post(route, data, action)
	// } else if method == "put" {
	// 	c.Client.Put(route, data, action)
	// } else if method == "delete" {
	// 	c.Client.Del(route, action)
	// } else {
	// 	err = cerr.NewUnknownError(correlationId, "UNSUPPORTED_METHOD", "Method is not supported by REST client").WithDetails("verb", method)
	// 	return nil, err
	// }
	return
}
