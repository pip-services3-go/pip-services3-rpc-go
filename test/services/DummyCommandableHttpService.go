package test_rpc_services

import (
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	"github.com/pip-services3-go/pip-services3-rpc-go/services"
)

type DummyCommandableHttpService struct {
	*services.CommandableHttpService
}

func NewDummyCommandableHttpService() *DummyCommandableHttpService {
	dchs := DummyCommandableHttpService{
		CommandableHttpService: services.NewCommandableHttpService("dummies"),
	}
	dchs.DependencyResolver.Put("controller", cref.NewDescriptor("pip-services-dummies", "controller", "default", "*", "*"))
	return &dchs
}
