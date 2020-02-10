package test_rpc_clients

// import (
// 	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
// 	"github.com/pip-services3-go/pip-services3-rpc-go/clients"
// )

// type DummyCommandableHttpClient struct {
// 	clients.CommandableHttpClient
// }

// func NewDummyCommandableHttpClient() *DummyCommandableHttpClient {
// 	dchc := DummyCommandableHttpClient{}
// 	dchc.CommandableHttpClient = *clients.NewCommandableHttpClient("dummy")
// 	return &dchc
// }

// func (c *DummyCommandableHttpClient) GetDummies(correlationId string, filter *cdata.FilterParams, paging *cdata.PagingParams) (result *DummyDataPage, err error) {

// 	params := make(map[string]interface{})
// 	params["filter"] = *filter
// 	params["paging"] = *paging
// 	c.CallCommand("get_dummies", correlationId, params)

// }

// func (c *DummyCommandableHttpClient) GetDummyById(correlationId string, dummyId string) (result *tesstrpc.Dummy, err error) {

// 	params := make(map[string]interface{})
// 	params["dummy_id"] = dummyId
// 	c.CallCommand(
// 		"get_dummy_by_id",
// 		correlationId,
// 		params)
// }

// func (c *DummyCommandableHttpClient) CreateDummy(correlationId string, dummy tesstrpc.Dummy) (result *tesstrpc.Dummy, err error) {

// 	params := make(map[string]interface{})
// 	params["dummy"] = dummy
// 	c.CallCommand(
// 		"create_dummy",
// 		correlationId,
// 		params)
// }

// func (c *DummyCommandableHttpClient) UpdateDummy(correlationId string, dummy tesstrpc.Dummy) (result *tesstrpc.Dummy, err error) {
// 	params := make(map[string]interface{})
// 	params["dummy"] = dummy
// 	c.CallCommand(
// 		"update_dummy",
// 		correlationId,
// 		params)
// }

// func (c *DummyCommandableHttpClient) DeleteDummy(correlationId string, dummyId string) (result *tesstrpc.Dummy) {

// 	params := make(map[string]interface{})
// 	params["dummy_id"] = dummyId
// 	c.CallCommand(
// 		"delete_dummy",
// 		correlationId,
// 		params)
// }
