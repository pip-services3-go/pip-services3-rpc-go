package test_rpc_clients

import (
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	"github.com/pip-services3-go/pip-services3-rpc-go/clients"
	testrpc "github.com/pip-services3-go/pip-services3-rpc-go/test"
)

type DummyDirectClient struct {
	clients.DirectClient
	concreateController testrpc.IDummyController
}

func NewDummyDirectClient() *DummyDirectClient {
	ddc := DummyDirectClient{}
	ddc.DirectClient = *clients.NewDirectClient()
	ddc.DependencyResolver.Put("controller", cref.NewDescriptor("pip-services-dummies", "controller", "*", "*", "*"))
	return &ddc
}

func (c *DummyDirectClient) SetReferences(references cref.IReferences) {
	c.DirectClient.SetReferences(references)

	concreateController, ok := c.Controller.(testrpc.IDummyController)
	if !ok {
		panic("DummyDirectClient: Cant't resolv dependency 'controller' to IDummyController")
	}
	c.concreateController = concreateController

}

func (c *DummyDirectClient) GetDummies(correlationId string, filter *cdata.FilterParams, paging *cdata.PagingParams) (result *testrpc.DummyDataPage, err error) {

	timing := c.Instrument(correlationId, "dummy.get_page_by_filter")
	result, err = c.concreateController.GetPageByFilter(correlationId, filter, paging)
	timing.EndTiming()
	return result, err

}

func (c *DummyDirectClient) GetDummyById(correlationId string, dummyId string) (result *testrpc.Dummy, err error) {

	timing := c.Instrument(correlationId, "dummy.get_one_by_id")
	result, err = c.concreateController.GetOneById(correlationId, dummyId)
	timing.EndTiming()
	return result, err
}

func (c *DummyDirectClient) CreateDummy(correlationId string, dummy testrpc.Dummy) (result *testrpc.Dummy, err error) {

	timing := c.Instrument(correlationId, "dummy.create")
	result, err = c.concreateController.Create(correlationId, dummy)
	timing.EndTiming()
	return result, err
}

func (c *DummyDirectClient) UpdateDummy(correlationId string, dummy testrpc.Dummy) (result *testrpc.Dummy, err error) {

	timing := c.Instrument(correlationId, "dummy.update")
	result, err = c.concreateController.Update(correlationId, dummy)
	timing.EndTiming()
	return result, err
}

func (c *DummyDirectClient) DeleteDummy(correlationId string, dummyId string) (result *testrpc.Dummy, err error) {

	timing := c.Instrument(correlationId, "dummy.delete_by_id")
	result, err = c.concreateController.DeleteById(correlationId, dummyId)
	timing.EndTiming()
	return result, err
}
