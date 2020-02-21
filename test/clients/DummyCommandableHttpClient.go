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

	calValue, calErr := c.CallCommand("get_dummies", correlationId, params, nil)
	if calErr != nil {
		return nil, calErr
	}

	convRes, err := clients.ConvertComandResult(calValue, reflect.TypeOf(&testrpc.DummyDataPage{}))
	result, _ = convRes.(*testrpc.DummyDataPage)
	return result, err
}

func (c *DummyCommandableHttpClient) GetDummyById(correlationId string, dummyId string) (result *testrpc.Dummy, err error) {

	params := cdata.NewEmptyStringValueMap()
	params.Put("dummy_id", dummyId)

	calValue, calErr := c.CallCommand("get_dummy_by_id", correlationId, params, nil)
	if calErr != nil {
		return nil, calErr
	}
	convRes, err := clients.ConvertComandResult(calValue, reflect.TypeOf(&testrpc.Dummy{}))
	result, _ = convRes.(*testrpc.Dummy)
	return result, err
}

func (c *DummyCommandableHttpClient) CreateDummy(correlationId string, dummy testrpc.Dummy) (result *testrpc.Dummy, err error) {

	bodyMap := make(map[string]interface{})
	bodyMap["dummy"] = dummy
	calValue, calErr := c.CallCommand("create_dummy", correlationId, nil, bodyMap)
	if calErr != nil {
		return nil, calErr
	}

	convRes, err := clients.ConvertComandResult(calValue, reflect.TypeOf(&testrpc.Dummy{}))
	result, _ = convRes.(*testrpc.Dummy)
	return result, err
}

func (c *DummyCommandableHttpClient) UpdateDummy(correlationId string, dummy testrpc.Dummy) (result *testrpc.Dummy, err error) {

	bodyMap := make(map[string]interface{})
	bodyMap["dummy"] = dummy
	calValue, calErr := c.CallCommand("update_dummy", correlationId, nil, bodyMap)
	if calErr != nil {
		return nil, calErr
	}
	convRes, err := clients.ConvertComandResult(calValue, reflect.TypeOf(&testrpc.Dummy{}))
	result, _ = convRes.(*testrpc.Dummy)
	return result, err
}

func (c *DummyCommandableHttpClient) DeleteDummy(correlationId string, dummyId string) (result *testrpc.Dummy, err error) {

	params := cdata.NewEmptyStringValueMap()
	params.Put("dummy_id", dummyId)

	calValue, calErr := c.CallCommand("delete_dummy", correlationId, params, nil)
	if calErr != nil {
		return nil, calErr
	}
	convRes, err := clients.ConvertComandResult(calValue, reflect.TypeOf(&testrpc.Dummy{}))
	result, _ = convRes.(*testrpc.Dummy)
	return result, err
}
