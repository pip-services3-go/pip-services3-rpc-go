package test_rpc_clients

import (
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	"github.com/pip-services3-go/pip-services3-rpc-go/clients"
)

type DummyDirectClient struct {
	clients.DirectClient
}

func NewDummyDirectClient() *DummyDirectClient {
	ddc := DummyDirectClient{}
	ddc.DirectClient = *clients.NewDirectClient()
	ddc.DependencyResolver.Put("controller", cref.NewDescriptor("pip-services-dummies", "controller", "*", "*", "*"))
	return &ddc
}

// func (c *DummyDirectClient) getDummies(correlationId string, filter *cdata.FilterParams, paging *cdata.PagingParams) (result testrpc.DummyDataPage, err error) {

// 	timing := c.Instrument(correlationId, "dummy.get_page_by_filter")
// 	result, err = c.Controller.GetPageByFilter(correlationId, filter, paging)

// 	timing.EndTiming()
// 	return result, err

// }

// func (c * DummyDirectClient) getDummyById(correlationId: string, dummyId: string, callback: (err: any, result: Dummy) => void): void {
//     let timing = c.instrument(correlationId, "dummy.get_one_by_id");
//     c._controller.getOneById(
//         correlationId,
//         dummyId,
//         (err, result) => {
//             timing.endTiming();
//             callback(err, result);
//         }
//     );
// }

// func (c * DummyDirectClient) createDummy(correlationId: string, dummy: any,
//     callback: (err: any, result: Dummy) => void): void {

//     let timing = c.instrument(correlationId, "dummy.create");
//     c._controller.create(
//         correlationId,
//         dummy,
//         (err, result) => {
//             timing.endTiming();
//             callback(err, result);
//         }
//     );
// }

// func (c * DummyDirectClient) updateDummy(correlationId: string, dummy: any,
//     callback: (err: any, result: Dummy) => void): void {

//     let timing = c.instrument(correlationId, "dummy.update");
//     c._controller.update(
//         correlationId,
//         dummy,
//         (err, result) => {
//             timing.endTiming();
//             callback(err, result);
//         }
//     );
// }

// func (c * DummyDirectClient) deleteDummy(correlationId: string, dummyId: string,
//     callback: (err: any, result: Dummy) => void): void {

//     let timing = c.instrument(correlationId, "dummy.delete_by_id");
//     c._controller.deleteById(
//         correlationId,
//         dummyId,
//         (err, result) => {
//             timing.endTiming();
//             callback(err, result);
//         }
//     );
// }
