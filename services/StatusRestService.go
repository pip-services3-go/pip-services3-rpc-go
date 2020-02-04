package services

import (
	"net/http"
	"time"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cconv "github.com/pip-services3-go/pip-services3-commons-go/convert"
	crefer "github.com/pip-services3-go/pip-services3-commons-go/refer"
	cinfo "github.com/pip-services3-go/pip-services3-components-go/info"
)

/*
Service that returns microservice status information via HTTP/REST protocol.

The service responds on /status route (can be changed) with a JSON object:
{
    - "id":            unique container id (usually hostname)
    - "name":          container name (from ContextInfo)
    - "description":   container description (from ContextInfo)
    - "start_time":    time when container was started
    - "current_time":  current time in UTC
    - "uptime":        duration since container start time in milliseconds
    - "properties":    additional container properties (from ContextInfo)
    - "components":    descriptors of components registered in the container
}

### Configuration parameters ###

- baseroute:              base route for remote URI
- route:                   status route (default: "status")
- dependencies:
  - endpoint:              override for HTTP Endpoint dependency
  - controller:            override for Controller dependency
- connection(s):
  - discovery_key:         (optional) a key to retrieve the connection from [[https://rawgit.com/pip-services-node/pip-services3-components-node/master/doc/api/interfaces/connect.idiscovery.html IDiscovery]]
  - protocol:              connection protocol: http or https
  - host:                  host name or IP address
  - port:                  port number
  - uri:                   resource URI or connection string with all parameters in it

### References ###

- <code>\*:logger:\*:\*:1.0</code>               (optional) [[https://rawgit.com/pip-services-node/pip-services3-components-node/master/doc/api/interfaces/log.ilogger.html ILogger]] components to pass log messages
- <code>\*:counters:\*:\*:1.0</code>             (optional) [[https://rawgit.com/pip-services-node/pip-services3-components-node/master/doc/api/interfaces/count.icounters.html ICounters]] components to pass collected measurements
- <code>\*:discovery:\*:\*:1.0</code>            (optional) [[https://rawgit.com/pip-services-node/pip-services3-components-node/master/doc/api/interfaces/connect.idiscovery.html IDiscovery]] services to resolve connection
- <code>\*:endpoint:http:\*:1.0</code>          (optional) [[HttpEndpoint]] reference

@see [[RestService]]
@see [[RestClient]]

### Example ###

    let service = new StatusService();
    service.configure(ConfigParams.fromTuples(
        "connection.protocol", "http",
        "connection.host", "localhost",
        "connection.port", 8080
    ));

    service.open("123", (err) => {
       console.log("The Status service is accessible at http://+:8080/status");
    });
*/
type StatusRestService struct {
	RestService
	startTime   time.Time
	references2 crefer.IReferences
	contextInfo *cinfo.ContextInfo
	route       string
}

/*
   Creates a new instance of c service.
*/
func NewStatusRestService() *StatusRestService {
	//super();
	srs := StatusRestService{}
	srs.startTime = time.Now()
	srs.route = "status"
	srs.DependencyResolver.Put("context-info", crefer.NewDescriptor("pip-services", "context-info", "default", "*", "1.0"))
	return &srs
}

/*
   Configures component by passing configuration parameters.

   @param config    configuration parameters to be set.
*/
func (c *StatusRestService) configure(config *cconf.ConfigParams) {
	c.RestService.Configure(config)
	c.route = config.GetAsStringWithDefault("route", c.route)
}

/*
	Sets references to dependent components.

	@param references 	references to locate the component dependencies.
*/
func (c *StatusRestService) SetReferences(references crefer.IReferences) {
	c.references2 = references
	c.RestService.SetReferences(references)

	depRes := c.DependencyResolver.GetOneOptional("context-info")
	if depRes != nil {
		c.contextInfo = depRes.(*cinfo.ContextInfo)
	}

}

/*
   Registers all service routes in HTTP endpoint.
*/
func (c *StatusRestService) Register() {
	c.RegisterRoute("get", c.route, nil, c.status)
}

/*
   Handles status requests

   @param req   an HTTP request
   @param res   an HTTP response
*/
func (c *StatusRestService) status(res http.ResponseWriter, req *http.Request) {

	id := ""
	if c.contextInfo != nil {
		id = c.contextInfo.ContextId
	}

	name := "Unknown"
	if c.contextInfo != nil {
		name = c.contextInfo.Name
	}

	description := ""
	if c.contextInfo != nil {
		description = c.contextInfo.Description
	}

	uptime := time.Now().Sub(c.startTime)

	properties := make(map[string]string, 0)
	if c.contextInfo != nil {
		properties = c.contextInfo.Properties
	}

	var components []string
	if c.references2 != nil {
		for _, locator := range c.references2.GetAllLocators() {
			components = append(components, cconv.StringConverter.ToString(locator))
		}
	}

	status := make(map[string]interface{})

	status["id"] = id
	status["name"] = name
	status["description"] = description
	status["start_time"] = cconv.StringConverter.ToString(c.startTime)
	status["current_time"] = cconv.StringConverter.ToString(time.Now())
	status["uptime"] = uptime
	status["properties"] = properties
	status["components"] = components

	c.SendResult(res, req, status, nil)
}
