package clients

import (
	"reflect"

	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
)

/* CommandableHttpClient is abstract client that calls commandable HTTP service.

Commandable services are generated automatically for ICommandable objects.
Each command is exposed as POST operation that receives all parameters
in body object.

 Configuration parameters

base_route:              base route for remote URI

- connection(s):
  - discovery_key:         (optional) a key to retrieve the connection from connect.idiscovery.html IDiscovery]]
  - protocol:              connection protocol: http or https
  - host:                  host name or IP address
  - port:                  port number
  - uri:                   resource URI or connection string with all parameters in it
- options:
  - retries:               number of retries (default: 3)
  - connect_timeout:       connection timeout in milliseconds (default: 10 sec)
  - timeout:               invocation timeout in milliseconds (default: 10 sec)

References:

- *:logger:*:*:1.0         (optional) ILogger components to pass log messages
- *:counters:*:*:1.0         (optional) ICounters components to pass collected measurements
- *:discovery:*:*:1.0        (optional) IDiscovery services to resolve connection

Example:

    class MyCommandableHttpClient extends CommandableHttpClient implements IMyClient {
       ...

        func (c * CommandableHttpClient) getData(correlationId: string, id: string,
           callback: (err: any, result: MyData) => void): void {

           c.callCommand(
               "get_data",
               correlationId,
               { id: id },
               (err, result) => {
                   callback(err, result);
               }
            );
        }
        ...
    }

    let client = new MyCommandableHttpClient();
    client.configure(ConfigParams.fromTuples(
        "connection.protocol", "http",
        "connection.host", "localhost",
        "connection.port", 8080
    ));

    client.getData("123", "1", (err, result) => {
    ...
    });
*/
type CommandableHttpClient struct {
	RestClient
}

// NewCommandableHttpClientr is creates a new instance of the client.
// - baseRoute     a base route for remote service.
func NewCommandableHttpClient(baseRoute string) *CommandableHttpClient {
	chc := CommandableHttpClient{}
	chc.RestClient = *NewRestClient()
	chc.BaseRoute = baseRoute
	return &chc
}

// CallCommand is calls a remote method via HTTP commadable protocol.
// The call is made via POST operation and all parameters are sent in body object.
// The complete route to remote method is defined as baseRoute + "/" + name.
//
// - name              a name of the command to call.
// - correlationId     (optional) transaction id to trace execution through call chain.
// - params            command parameters.
// - callback          callback function that receives result or error.

func (c *CommandableHttpClient) CallCommand(name string, correlationId string, params *cdata.StringValueMap,
	data interface{}, prototype reflect.Type) (result interface{}, err error) {
	timing := c.Instrument(correlationId, c.BaseRoute+"."+name)
	cRes, cErr := c.Call("post", name, correlationId, params, data, prototype)
	timing.EndTiming()
	return c.InstrumentError(correlationId, c.BaseRoute+"."+name, cErr, cRes)
}
