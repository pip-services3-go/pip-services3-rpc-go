package test_rpc_clients

// import (
// 	"encoding/json"

// 	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
// 	"github.com/pip-services3-go/pip-services3-rpc-go/clients"
// 	testrpc "github.com/pip-services3-go/pip-services3-rpc-go/test"
// )

// type DummyRestClient struct {
// 	clients.RestClient
// }

// func NewDummyRestClient() *DummyRestClient {
// 	drc := DummyRestClient{}
// 	drc.RestClient = *clients.NewRestClient()
// 	return &drc
// }

// func (c *DummyRestClient) GetDummies(correlationId string, filter *cdata.FilterParams, paging *cdata.PagingParams) (result *testrpc.DummyDataPage, err error) {
// 	var params cdata.StringValueMap
// 	c.AddFilterParams(&params, filter)
// 	c.AddPagingParams(&params, paging)

// 	var data testrpc.DummyDataPage
// 	calValue, calErr := c.Call("get", "/dummies", correlationId, &params, nil)
// 	if calErr != nil {
// 		return nil, calErr
// 	}
// 	convErr := json.Unmarshal(calValue.([]byte), &data)
// 	if convErr != nil {
// 		return nil, convErr
// 	}

// 	c.Instrument(correlationId, "dummy.get_page_by_filter")
// 	return &data, nil
// }

// func (c *DummyRestClient) GetDummyById(correlationId string, dummyId string) (result *testrpc.Dummy, err error) {
// 	calValue, calErr := c.Call("get", "/dummies/"+dummyId, correlationId, nil, nil)

// 	if calErr != nil {
// 		return nil, calErr
// 	}
// 	var data testrpc.Dummy
// 	convErr := json.Unmarshal(calValue.([]byte), &data)
// 	if convErr != nil {
// 		return nil, convErr
// 	}

// 	c.Instrument(correlationId, "dummy.get_one_by_id")
// 	return &data, nil
// }

// func (c *DummyRestClient) CreateDummy(correlationId string, dummy testrpc.Dummy) (result *testrpc.Dummy, err error) {
// 	calValue, calErr := c.Call("post", "/dummies", correlationId, nil, dummy)
// 	if calErr != nil {
// 		return nil, calErr
// 	}
// 	var data testrpc.Dummy
// 	convErr := json.Unmarshal(calValue.([]byte), &data)
// 	if convErr != nil {
// 		return nil, convErr
// 	}
// 	c.Instrument(correlationId, "dummy.create")
// 	return &data, nil
// }

// func (c *DummyRestClient) UpdateDummy(correlationId string, dummy testrpc.Dummy) (result *testrpc.Dummy, err error) {
// 	calValue, calErr := c.Call("put", "/dummies", correlationId, nil, dummy)
// 	if calErr != nil {
// 		return nil, calErr
// 	}
// 	var data testrpc.Dummy
// 	convErr := json.Unmarshal(calValue.([]byte), &data)
// 	if convErr != nil {
// 		return nil, convErr
// 	}
// 	c.Instrument(correlationId, "dummy.update")
// 	return &data, nil
// }

// func (c *DummyRestClient) DeleteDummy(correlationId string, dummyId string) (result *testrpc.Dummy, err error) {
// 	calValue, calErr := c.Call("delete", "/dummies/"+dummyId, correlationId, nil, nil)
// 	if calErr != nil {
// 		return nil, calErr
// 	}
// 	var data testrpc.Dummy
// 	convErr := json.Unmarshal(calValue.([]byte), &data)
// 	if convErr != nil {
// 		return nil, convErr
// 	}
// 	c.Instrument(correlationId, "dummy.delete_by_id")
// 	return &data, nil

// }
