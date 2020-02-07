package build

import (
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	cbuild "github.com/pip-services3-go/pip-services3-components-go/build"
	"github.com/pip-services3-go/pip-services3-rpc-go/services"
)

/*
Creates RPC components by their descriptors.

See [[https://rawgit.com/pip-services-node/pip-services3-components-node/master/doc/api/classes/build.factory.html Factory]]
See [[HttpEndpoint]]
See [[HeartbeatRestService]]
See [[StatusRestService]]
*/
type DefaultRpcFactory struct {
	cbuild.Factory
	Descriptor                 *cref.Descriptor
	HttpEndpointDescriptor     *cref.Descriptor
	StatusServiceDescriptor    *cref.Descriptor
	HeartbeatServiceDescriptor *cref.Descriptor
}

/*
	Create a new instance of the factory.
*/
func NewDefaultRpcFactory() *DefaultRpcFactory {
	drf := DefaultRpcFactory{}
	drf.Factory = *cbuild.NewFactory()
	drf.Descriptor = cref.NewDescriptor("pip-services", "factory", "rpc", "default", "1.0")
	drf.HttpEndpointDescriptor = cref.NewDescriptor("pip-services", "endpoint", "http", "*", "1.0")
	drf.StatusServiceDescriptor = cref.NewDescriptor("pip-services", "status-service", "http", "*", "1.0")
	drf.HeartbeatServiceDescriptor = cref.NewDescriptor("pip-services", "heartbeat-service", "http", "*", "1.0")

	drf.RegisterType(drf.HttpEndpointDescriptor, services.NewHttpEndpoint)
	drf.RegisterType(drf.HeartbeatServiceDescriptor, services.NewHeartbeatRestService)
	drf.RegisterType(drf.StatusServiceDescriptor, services.NewStatusRestService)
	return &drf
}
