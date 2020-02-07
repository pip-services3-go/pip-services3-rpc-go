package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	crefer "github.com/pip-services3-go/pip-services3-commons-go/refer"
	cvalid "github.com/pip-services3-go/pip-services3-commons-go/validate"
	ccount "github.com/pip-services3-go/pip-services3-components-go/count"
	clog "github.com/pip-services3-go/pip-services3-components-go/log"
	"github.com/pip-services3-go/pip-services3-rpc-go/connect"
)

/*
Used for creating HTTP endpoints. An endpoint is a URL, at which a given service can be accessed by a client.

 Configuration parameters

Parameters to pass to the configure method for component configuration:

- connection(s) - the connection resolver"s connections:
    - "connection.discovery_key" - the key to use for connection resolving in a discovery service;
    - "connection.protocol" - the connection"s protocol;
    - "connection.host" - the target host;
    - "connection.port" - the target port;
    - "connection.uri" - the target URI.
- credential - the HTTPS credentials:
    - "credential.ssl_key_file" - the SSL func (c *HttpEndpoint )key in PEM
    - "credential.ssl_crt_file" - the SSL certificate in PEM
    - "credential.ssl_ca_file" - the certificate authorities (root cerfiticates) in PEM

References:

A logger, counters, and a connection resolver can be referenced by passing the
following references to the object"s setReferences method:

- logger: "\*:logger:\*:\*:1.0";
- counters: "\*:counters:\*:\*:1.0";
- discovery: "\*:discovery:\*:\*:1.0" (for the connection resolver).

Examples:

    func (c *HttpEndpoint ) MyMethod(_config: ConfigParams, _references: IReferences) {
        let endpoint = new HttpEndpoint();
        if (c._config)
            endpoint.configure(c._config);
        if (c._references)
            endpoint.setReferences(c._references);
        ...

        c._endpoint.open(correlationId, (err) => {
                c._opened = err == nil;
                callback(err);
            });
        ...
    }
*/

//implements IOpenable, IConfigurable, IReferenceable

type HttpEndpoint struct {
	defaultConfig          *cconf.ConfigParams
	server                 *http.Server
	router                 *mux.Router
	connectionResolver     *connect.HttpConnectionResolver
	logger                 *clog.CompositeLogger
	counters               *ccount.CompositeCounters
	maintenanceEnabled     bool
	fileMaxSize            int64
	protocolUpgradeEnabled bool
	uri                    string
	registrations          []IRegisterable
}

// NewHttpEndpoint creates new HttpEndpoint
func NewHttpEndpoint() *HttpEndpoint {
	he := HttpEndpoint{}
	he.defaultConfig = cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "0.0.0.0",
		"connection.port", "3000",

		"credential.ssl_key_file", nil,
		"credential.ssl_crt_file", nil,
		"credential.ssl_ca_file", nil,

		"options.maintenance_enabled", false,
		"options.request_max_size", 1024*1024,
		"options.file_max_size", 200*1024*1024,
		"options.connect_timeout", "60000",
		"options.debug", "true",
	)
	he.connectionResolver = connect.NewHttpConnectionResolver()
	he.logger = clog.NewCompositeLogger()
	he.counters = ccount.NewCompositeCounters()
	he.maintenanceEnabled = false
	he.fileMaxSize = 200 * 1024 * 1024
	he.protocolUpgradeEnabled = false
	he.registrations = make([]IRegisterable, 0, 0)
	return &he
}

/*
   Configures this HttpEndpoint using the given configuration parameters.

   __Configuration parameters:__
   - __connection(s)__ - the connection resolver"s connections;
       - "connection.discovery_key" - the key to use for connection resolving in a discovery service;
       - "connection.protocol" - the connection"s protocol;
       - "connection.host" - the target host;
       - "connection.port" - the target port;
       - "connection.uri" - the target URI.
       - "credential.ssl_key_file" - SSL func (c *HttpEndpoint )key in PEM
       - "credential.ssl_crt_file" - SSL certificate in PEM
       - "credential.ssl_ca_file" - Certificate authority (root certificate) in PEM

   - config    configuration parameters, containing a "connection(s)" section.

   @see https://rawgit.com/pip-services-node/pip-services3-commons-node/master/doc/api/classes/config.configparams.html ConfigParams (in the PipServices "Commons" package)
*/
func (c *HttpEndpoint) Configure(config *cconf.ConfigParams) {
	config = config.SetDefaults(c.defaultConfig)
	c.connectionResolver.Configure(config)

	c.maintenanceEnabled = config.GetAsBooleanWithDefault("options.maintenance_enabled", c.maintenanceEnabled)
	c.fileMaxSize = config.GetAsLongWithDefault("options.file_max_size", c.fileMaxSize)
	c.protocolUpgradeEnabled = config.GetAsBooleanWithDefault("options.protocol_upgrade_enabled", c.protocolUpgradeEnabled)
}

/*
   Sets references to this endpoint"s logger, counters, and connection resolver.

   __References:__
   - logger: "\*:logger:\*:\*:1.0"
   - counters: "\*:counters:\*:\*:1.0"
   - discovery: "\*:discovery:\*:\*:1.0" (for the connection resolver)

   - references    an IReferences object, containing references to a logger, counters,
                        and a connection resolver.

   @see https://rawgit.com/pip-services-node/pip-services3-commons-node/master/doc/api/interfaces/refer.ireferences.html IReferences (in the PipServices "Commons" package)
*/
func (c *HttpEndpoint) SetReferences(references crefer.IReferences) {
	c.logger.SetReferences(references)
	c.counters.SetReferences(references)
	c.connectionResolver.SetReferences(references)
}

/*
   @returns whether or not this endpoint is open with an actively listening REST server.
*/
func (c *HttpEndpoint) IsOpen() bool {
	return c.server != nil
}

//TODO: check for correct understanding.
/*
   Opens a connection using the parameters resolved by the referenced connection
   resolver and creates a REST server (service) using the set options and parameters.

   - correlationId     (optional) transaction id to trace execution through call chain.
   - callback          (optional) the function to call once the opening process is complete.
                            Will be called with an error if one is raised.
*/
func (c *HttpEndpoint) Open(correlationId string) error {
	if c.IsOpen() {
		return nil
	}

	connection, credential, err := c.connectionResolver.Resolve(correlationId)

	if err != nil {
		return err
	}

	c.uri = connection.Uri()
	url := connection.Host() + ":" + strconv.Itoa(connection.Port())
	c.server = &http.Server{Addr: url}
	c.router = mux.NewRouter()

	c.router.Use(c.addCors)
	c.router.Use(c.addCompatibility)
	c.router.Use(c.noCache)
	c.router.Use(c.doMaintenance)

	c.server.Handler = c.router

	c.performRegistrations()

	if connection.Protocol() == "https" { //"http"
		sslKeyFile := credential.GetAsString("ssl_key_file")
		sslCrtFile := credential.GetAsString("ssl_crt_file")

		go func() {
			servErr := c.server.ListenAndServeTLS(sslKeyFile, sslCrtFile)
			if servErr != nil {
				fmt.Println("Server error %s", servErr.Error())
			}
		}()

	} else {
		go func() {
			servErr := c.server.ListenAndServe()
			if servErr != nil {
				fmt.Println("Server error %s", servErr.Error())
			}
		}()
	}

	regErr := c.connectionResolver.Register(correlationId)
	if regErr != nil {
		c.logger.Error(correlationId, regErr, "ERROR_REG_SRV", "Can't register REST service at %s", c.uri)
	}
	c.logger.Debug(correlationId, "Opened REST service at %s", c.uri)
	return regErr
}

func (c *HttpEndpoint) addCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		next.ServeHTTP(w, r)
	})
}

func (c *HttpEndpoint) addCompatibility(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//TODO: Write code

		// req.param = (name) => {
		//     if (req.query) {
		//         let param = req.query[name];
		//         if (param) return param;
		//     }
		//     if (req.body) {
		//         let param = req.body[name];
		//         if (param) return param;
		//     }
		//     if (req.params) {
		//         let param = req.params[name];
		//         if (param) return param;
		//     }
		//     return nil;

		// }
		// req.route.params = req.params;

		next.ServeHTTP(w, r)
	})
}

// Prevents IE from caching REST requests
func (c *HttpEndpoint) noCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Add("Pragma", "no-cache")
		w.Header().Add("Expires", "0")
		next.ServeHTTP(w, r)
	})
}

// Returns maintenance error code
func (c *HttpEndpoint) doMaintenance(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Make this more sophisticated
		if c.maintenanceEnabled {
			w.Header().Add("Retry-After", "3600")
			jsonStr, _ := json.Marshal(503)
			w.Write(jsonStr)
			next.ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

/*
Closes this endpoint and the REST server (service) that was opened earlier.

- correlationId     (optional) transaction id to trace execution through call chain.
- callback          (optional) the function to call once the closing process is complete.
                         Will be called with an error if one is raised.
*/
func (c *HttpEndpoint) Close(correlationId string) error {
	if c.server != nil {
		// Attempt a graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		clErr := c.server.Shutdown(ctx)
		if clErr != nil {
			c.logger.Warn(correlationId, "Failed while closing REST service: %s", clErr.Error())
			return clErr
		}
		c.logger.Debug(correlationId, "Closed REST service at %s", c.uri)
		c.server = nil
		c.uri = ""
	}

	return nil
}

/*
Registers a registerable object for dynamic endpoint discovery.

- registration      the registration to add.

@see IRegisterable
*/
func (c *HttpEndpoint) Register(registration IRegisterable) {
	c.registrations = append(c.registrations, registration)
}

/*
Unregisters a registerable object, so that it is no longer used in dynamic
endpoint discovery.

- registration      the registration to remove.

@see IRegisterable
*/
func (c *HttpEndpoint) Unregister(registration IRegisterable) {
	for i := 0; i < len(c.registrations); {
		if c.registrations[i] == registration {
			if i == len(c.registrations)-1 {
				c.registrations = c.registrations[:i]
			} else {
				c.registrations = append(c.registrations[:i], c.registrations[i+1:]...)
			}
		} else {
			i++
		}
	}
}

func (c *HttpEndpoint) performRegistrations() {
	for _, registration := range c.registrations {
		registration.Register()
	}
}

func (c *HttpEndpoint) fixRoute(route string) string {
	if len(route) > 0 && !strings.HasPrefix(route, "/") {
		route = "/" + route
	}
	return route
}

/*
Registers an action in this objects REST server (service) by the given method and route.

- method        the HTTP method of the route.
- route         the route to register in this object"s REST server (service).
- schema        the schema to use for parameter validation.
- action        the action to perform at the given route.
*/
func (c *HttpEndpoint) RegisterRoute(method string, route string, schema *cvalid.Schema,
	action http.HandlerFunc) {

	method = strings.ToLower(method)
	if method == "del" {
		method = "delete"
	}

	route = c.fixRoute(route)

	// Hack!!! Wrapping action to preserve prototyping context
	actionCurl := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//  Perform validation
		if schema != nil {
			//params = _.extend({}, req.params, { body: req.body })
			var params map[string]interface{} = make(map[string]interface{}, 0)
			for k, v := range r.URL.Query() {
				params[k] = v[0]
			}

			for k, v := range mux.Vars(r) {
				params[k] = v
			}

			// Make copy of request
			bodyBuf, bodyErr := ioutil.ReadAll(r.Body)
			if bodyErr != nil {
				HttpResponseSender.SendError(w, r, bodyErr)
				return
			}
			r.Body.Close()
			r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBuf))
			//-------------------------
			var body interface{}
			json.Unmarshal(bodyBuf, &body)
			params["body"] = body

			correlationId := r.URL.Query().Get("correlaton_id")
			err := schema.ValidateAndReturnError(correlationId, params, false)
			if err != nil {
				HttpResponseSender.SendError(w, r, err)
				return
			}
		}
		action(w, r)
	})

	// Wrapping to preserve "this"
	// let self = c;
	// c.server[method](route, actionCurl);
	c.router.Handle(route, actionCurl).Methods(strings.ToUpper(method))
}

/*
Registers an action with authorization in this objects REST server (service)
by the given method and route.

- method        the HTTP method of the route.
- route         the route to register in this object"s REST server (service).
- schema        the schema to use for parameter validation.
- authorize     the authorization interceptor
- action        the action to perform at the given route.
*/
func (c *HttpEndpoint) RegisterRouteWithAuth(method string, route string, schema *cvalid.Schema,
	authorize func(w http.ResponseWriter, r *http.Request, user *cdata.AnyValueMap, next http.HandlerFunc),
	action http.HandlerFunc) {

	if authorize != nil {
		nextAction := action
		action = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorize(w, r, nil, nextAction)
		})
	}

	c.RegisterRoute(method, route, schema, action)
}

/*
Registers a middleware action for the given route.

- route         the route to register in this object"s REST server (service).
- action        the middleware action to perform at the given route.
*/
func (c *HttpEndpoint) RegisterInterceptor(route string, action func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)) {

	route = c.fixRoute(route)
	interceptorFunc := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if route != "" && !strings.HasPrefix(r.URL.String(), route) {
				next.ServeHTTP(w, r)
			} else {
				action(w, r, next.ServeHTTP)
			}
		})
	}
	c.router.Use(interceptorFunc)
}
