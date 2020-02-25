package test_rpc_clients

import (
	"reflect"

	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	"github.com/pip-services3-go/pip-services3-rpc-go/clients"
	testrpc "github.com/pip-services3-go/pip-services3-rpc-go/test"
)

var (
	dummyDataPageType = reflect.TypeOf(&testrpc.DummyDataPage{})
	dummyType         = reflect.TypeOf(&testrpc.Dummy{})
)

type DummyRestClient struct {
	clients.RestClient
}

func NewDummyRestClient() *DummyRestClient {
	drc := DummyRestClient{}
	drc.RestClient = *clients.NewRestClient()
	return &drc
}

func (c *DummyRestClient) GetDummies(correlationId string, filter *cdata.FilterParams,
	paging *cdata.PagingParams) (result *testrpc.DummyDataPage, err error) {

	params := cdata.NewEmptyStringValueMap()
	c.AddFilterParams(params, filter)
	c.AddPagingParams(params, paging)

	calValue, calErr := c.Call("get", "/dummies", correlationId, params, nil, dummyDataPageType)
	if calErr != nil {
		return nil, calErr
	}

	result, _ = calValue.(*testrpc.DummyDataPage)
	c.Instrument(correlationId, "dummy.get_page_by_filter")
	return result, nil
}

func (c *DummyRestClient) GetDummyById(correlationId string, dummyId string) (result *testrpc.Dummy, err error) {
	calValue, calErr := c.Call("get", "/dummies/"+dummyId, correlationId, nil, nil, dummyType)

	if calErr != nil {
		return nil, calErr
	}

	result, _ = calValue.(*testrpc.Dummy)
	c.Instrument(correlationId, "dummy.get_one_by_id")
	return result, nil
}

func (c *DummyRestClient) CreateDummy(correlationId string, dummy testrpc.Dummy) (result *testrpc.Dummy, err error) {
	calValue, calErr := c.Call("post", "/dummies", correlationId, nil, dummy, dummyType)
	if calErr != nil {
		return nil, calErr
	}

	result, _ = calValue.(*testrpc.Dummy)
	c.Instrument(correlationId, "dummy.create")
	return result, nil
}

func (c *DummyRestClient) UpdateDummy(correlationId string, dummy testrpc.Dummy) (result *testrpc.Dummy, err error) {
	calValue, calErr := c.Call("put", "/dummies", correlationId, nil, dummy, dummyType)
	if calErr != nil {
		return nil, calErr
	}

	result, _ = calValue.(*testrpc.Dummy)
	c.Instrument(correlationId, "dummy.update")
	return result, nil
}

func (c *DummyRestClient) DeleteDummy(correlationId string, dummyId string) (result *testrpc.Dummy, err error) {
	calValue, calErr := c.Call("delete", "/dummies/"+dummyId, correlationId, nil, nil, dummyType)
	if calErr != nil {
		return nil, calErr
	}

	result, _ = calValue.(*testrpc.Dummy)
	c.Instrument(correlationId, "dummy.delete_by_id")
	return result, nil
}
