package test_rpc_clients

import (
	"reflect"

	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	"github.com/pip-services3-go/pip-services3-rpc-go/clients"
	testrpc "github.com/pip-services3-go/pip-services3-rpc-go/test"
)

type DummyCommandableHttpClient struct {
	clients.CommandableHttpClient
}

func NewDummyCommandableHttpClient() *DummyCommandableHttpClient {
	dchc := DummyCommandableHttpClient{}
	dchc.CommandableHttpClient = *clients.NewCommandableHttpClient("dummies")
	return &dchc
}

func (c *DummyCommandableHttpClient) GetDummies(correlationId string, filter *cdata.FilterParams, paging *cdata.PagingParams) (result *testrpc.DummyDataPage, err error) {

	params := cdata.NewEmptyStringValueMap()
	c.AddFilterParams(params, filter)
	c.AddPagingParams(params, paging)

	calValue, calErr := c.CallCommand(dummyDataPageType, "get_dummies", correlationId, cdata.NewAnyValueMapFromValue(params.Value()))
	if calErr != nil {
		return nil, calErr
	}
	result, _ = calValue.(*testrpc.DummyDataPage)
	return result, err
}

func (c *DummyCommandableHttpClient) GetDummyById(correlationId string, dummyId string) (result *testrpc.Dummy, err error) {

	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy_id", dummyId)

	calValue, calErr := c.CallCommand(dummyType, "get_dummy_by_id", correlationId, params)
	if calErr != nil {
		return nil, calErr
	}
	result, _ = calValue.(*testrpc.Dummy)
	return result, err
}

func (c *DummyCommandableHttpClient) CreateDummy(correlationId string, dummy testrpc.Dummy) (result *testrpc.Dummy, err error) {

	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy", dummy)

	calValue, calErr := c.CallCommand(dummyType, "create_dummy", correlationId, params)
	if calErr != nil {
		return nil, calErr
	}
	result, _ = calValue.(*testrpc.Dummy)
	return result, err
}

func (c *DummyCommandableHttpClient) UpdateDummy(correlationId string, dummy testrpc.Dummy) (result *testrpc.Dummy, err error) {

	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy", dummy)

	calValue, calErr := c.CallCommand(dummyType, "update_dummy", correlationId, params)
	if calErr != nil {
		return nil, calErr
	}
	result, _ = calValue.(*testrpc.Dummy)
	return result, err
}

func (c *DummyCommandableHttpClient) DeleteDummy(correlationId string, dummyId string) (result *testrpc.Dummy, err error) {

	params := cdata.NewEmptyAnyValueMap()
	params.Put("dummy_id", dummyId)

	calValue, calErr := c.CallCommand(dummyType, "delete_dummy", correlationId, params)
	if calErr != nil {
		return nil, calErr
	}
	result, _ = calValue.(*testrpc.Dummy)
	return result, err
}

func (c *DummyCommandableHttpClient) CheckCorrelationId(correlationId string) (result map[string]string, err error) {

	params := cdata.NewEmptyAnyValueMap()

	calValue, calErr := c.CallCommand(reflect.TypeOf(make(map[string]string)), "check_correlation_id", correlationId, params)
	if calErr != nil {
		return nil, calErr
	}
	val, _ := calValue.(*(map[string]string))
	return *val, err
}
