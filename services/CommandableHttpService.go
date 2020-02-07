package services

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	ccomands "github.com/pip-services3-go/pip-services3-commons-go/commands"
	crun "github.com/pip-services3-go/pip-services3-commons-go/run"
)

// /* @module services */
// import { ICommandable } from "pip-services3-commons-node";
// import { CommandSet } from "pip-services3-commons-node";
// import { Parameters } from "pip-services3-commons-node";

// import { RestService } from "./RestService";

/*
Abstract service that receives remove calls via HTTP/REST protocol
to operations automatically generated for commands defined in [[https://rawgit.com/pip-services-node/pip-services3-commons-node/master/doc/api/interfaces/commands.icommandable.html ICommandable components]].
Each command is exposed as POST operation that receives all parameters in body object.

Commandable services require only 3 lines of code to implement a robust external
HTTP-based remote interface.

### Configuration parameters ###

- base_route:              base route for remote URI
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

@see [[CommandableHttpClient]]
@see [[RestService]]

### Example ###

    class MyCommandableHttpService extends CommandableHttpService {
       public constructor() {
          base();
          c._dependencyResolver.put(
              "controller",
              new Descriptor("mygroup","controller","*","*","1.0")
          );
       }
    }

    let service = new MyCommandableHttpService();
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
// extends RestService

type CommandableHttpService struct {
	RestService
	commandSet *ccomands.CommandSet
}

/*
   Creates a new instance of the service.

   @param baseRoute a service base route.
*/
func NewCommandableHttpService(baseRoute string) *CommandableHttpService {
	chs := CommandableHttpService{}
	chs.RestService = *NewRestService()
	chs.RestService.IRegisterable = &chs
	chs.BaseRoute = baseRoute
	chs.DependencyResolver.Put("controller", "none")
	return &chs
}

/*
   Registers all service routes in HTTP endpoint.
*/
func (c *CommandableHttpService) Register() {
	resCtrl, depErr := c.DependencyResolver.GetOneRequired("controller")
	if depErr != nil {
		return
	}
	controller, ok := resCtrl.(ccomands.ICommandable)
	if !ok {
		c.Logger.Error("CommandableHttpService", nil, "Can't cast Controller to ICommandable")
		return
	}
	c.commandSet = controller.GetCommandSet()

	commands := c.commandSet.Commands()
	for index := 0; index < len(commands); index++ {
		command := commands[index]

		route := command.Name()
		if route[0] != "/"[0] {
			route = "/" + route
		}

		c.RegisterRoute("post", route, nil, func(res http.ResponseWriter, req *http.Request) {

			// Make copy of request
			bodyBuf, bodyErr := ioutil.ReadAll(req.Body)
			if bodyErr != nil {
				HttpResponseSender.SendError(res, req, bodyErr)
				return
			}
			req.Body.Close()
			req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBuf))
			//-------------------------

			var params interface{}
			json.Unmarshal(bodyBuf, &params)
			urlParams := req.URL.Query()
			correlationId := urlParams.Get("correlation_id")
			args := crun.NewParametersFromValue(params)
			timing := c.Instrument(correlationId, c.BaseRoute+"."+command.Name())

			execRes, execErr := command.Execute(correlationId, args)
			timing.EndTiming()
			instrRes, instrErr := c.InstrumentError(correlationId,
				c.BaseRoute+"."+command.Name(),
				execErr, execRes)
			c.SendResult(res, req, instrRes, instrErr)

		})
	}
}
