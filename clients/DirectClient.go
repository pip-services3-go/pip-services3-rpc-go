package clients

import (
	cconf "github.com/pip-services3-go/pip-services3-commons-go/v3/config"
	cerr "github.com/pip-services3-go/pip-services3-commons-go/v3/errors"
	crefer "github.com/pip-services3-go/pip-services3-commons-go/v3/refer"
	ccount "github.com/pip-services3-go/pip-services3-components-go/v3/count"
	clog "github.com/pip-services3-go/pip-services3-components-go/v3/log"
)

/*
Abstract client that calls controller directly in the same memory space.

It is used when multiple microservices are deployed in a single container (monolyth)
and communication between them can be done by direct calls rather then through
the network.

Configuration parameters:

- dependencies:
  - controller:            override controller descriptor

References:

- \*:logger:\*:\*:1.0         (optional) [[https://rawgit.com/pip-services-node/pip-services3-components-node/master/doc/api/interfaces/log.ilogger.html ILogger]] components to pass log messages
- \*:counters:\*:\*:1.0       (optional) [[https://rawgit.com/pip-services-node/pip-services3-components-node/master/doc/api/interfaces/count.icounters.html ICounters]] components to pass collected measurements
- \*:controller:\*:\*:1.0     controller to call business methods

Example:

    class MyDirectClient extends DirectClient<IMyController> implements IMyClient {

        func (c * DirectClient) constructor() {
          super();
          c.DependencyResolver.put("controller", new Descriptor(
              "mygroup", "controller", "*", "*", "*"));
        }
        ...

        func (c * DirectClient) getData(correlationId: string, id: string,
          callback: (err: any, result: MyData) => void): void {

          let timing = c.instrument(correlationId, "myclient.get_data");
          c.Controller.getData(correlationId, id, (err, result) => {
             timing.endTiming();
             c.instrumentError(correlationId, "myclient.get_data", err, result, callback);
          });
        }
        ...
    }

    let client = new MyDirectClient();
    client.setReferences(References.fromTuples(
        new Descriptor("mygroup","controller","default","default","1.0"), controller
    ));

    client.getData("123", "1", (err, result) => {
      ...
    });
*/

//implements IConfigurable, IReferenceable, IOpenable

type DirectClient struct {
	//The controller reference.
	Controller interface{}
	//The open flag.
	Opened bool
	//The logger.
	Logger *clog.CompositeLogger
	//The performance counters
	Counters *ccount.CompositeCounters
	//The dependency resolver to get controller reference.
	DependencyResolver crefer.DependencyResolver
}

/*
   Creates a new instance of the client.
*/
func NewDirectClient() *DirectClient {
	dc := DirectClient{
		Opened:             true,
		Logger:             clog.NewCompositeLogger(),
		Counters:           ccount.NewCompositeCounters(),
		DependencyResolver: *crefer.NewDependencyResolver(),
	}
	dc.DependencyResolver.Put("controller", "none")
	return &dc
}

/*
   Configures component by passing configuration parameters.
    *
   - config    configuration parameters to be set.
*/
func (c *DirectClient) Configure(config *cconf.ConfigParams) {
	c.DependencyResolver.Configure(config)
}

/*
	Sets references to dependent components.
	 *
	- references 	references to locate the component dependencies.
*/
func (c *DirectClient) SetReferences(references crefer.IReferences) {
	c.Logger.SetReferences(references)
	c.Counters.SetReferences(references)
	c.DependencyResolver.SetReferences(references)
	res, cErr := c.DependencyResolver.GetOneRequired("controller")
	if cErr != nil {
		panic("DirectClient: Cant't resolv dependency 'controller'")
	}
	c.Controller = res
}

/*
   Adds instrumentation to log calls and measure call time.
   It returns a Timing object that is used to end the time measurement.
    *
   - correlationId     (optional) transaction id to trace execution through call chain.
   - name              a method name.
   // Returns Timing object to end the time measurement.
*/
func (c *DirectClient) Instrument(correlationId string, name string) *ccount.Timing {
	c.Logger.Trace(correlationId, "Calling %s method", name)
	c.Counters.IncrementOne(name + ".call_count")
	return c.Counters.BeginTiming(name + ".call_time")
}

/*
   Adds instrumentation to error handling.
    *
   - correlationId     (optional) transaction id to trace execution through call chain.
   - name              a method name.
   - err               an occured error
   - result            (optional) an execution result
   - callback          (optional) an execution callback
*/
func (c *DirectClient) InstrumentError(correlationId string, name string, inErr error, inRes interface{}) (result interface{}, err error) {
	if inErr != nil {
		c.Logger.Error(correlationId, inErr, "Failed to call %s method", name)
		c.Counters.IncrementOne(name + ".call_errors")
	}
	return inRes, inErr
}

/*
	Checks if the component is opened.
	 *
	// Returns true if the component has been opened and false otherwise.
*/
func (c *DirectClient) IsOpen() bool {
	return c.Opened
}

/*
	Opens the component.
	 *
	- correlationId 	(optional) transaction id to trace execution through call chain.
    - callback 			callback function that receives error or nil no errors occured.
*/
func (c *DirectClient) Open(correlationId string) error {
	if c.Opened {
		return nil
	}

	if c.Controller == nil {
		err := cerr.NewConnectionError(correlationId, "NO_CONTROLLER", "Controller reference is missing")
		return err
	}

	c.Opened = true

	c.Logger.Info(correlationId, "Opened direct client")
	return nil
}

/*
	Closes component and frees used resources.
	 *
	- correlationId 	(optional) transaction id to trace execution through call chain.
    - callback 			callback function that receives error or nil no errors occured.
*/
func (c *DirectClient) Close(correlationId string) error {
	if c.Opened {
		c.Logger.Info(correlationId, "Closed direct client")
	}
	c.Opened = false
	return nil
}
