package test_rpc_services

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	testrpc "github.com/pip-services3-go/pip-services3-rpc-go/test"
	"github.com/stretchr/testify/assert"
)

func TestDummyCommandableHttpService(t *testing.T) {

	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", "3000",
	)
	var _dummy1 testrpc.Dummy
	//var _dummy2 testrpc.Dummy

	var service *DummyCommandableHttpService

	ctrl := testrpc.NewDummyController()

	service = NewDummyCommandableHttpService()
	service.Configure(restConfig)

	references := cref.NewReferencesFromTuples(
		cref.NewDescriptor("pip-services-dummies", "controller", "default", "default", "1.0"), ctrl,
		cref.NewDescriptor("pip-services-dummies", "service", "http", "default", "1.0"), service,
	)
	service.SetReferences(references)

	service.Open("")
	defer service.Close("")

	url := "http://localhost:3000"

	_dummy1 = testrpc.Dummy{Id: "", Key: "Key 1", Content: "Content 1"}
	//_dummy2 = testrpc.Dummy{Id: "", Key: "Key 2", Content: "Content 2"}

	// Create one dummy

	bodyMap := make(map[string]interface{})
	bodyMap["dummy"] = _dummy1

	jsonBody, _ := json.Marshal(bodyMap)

	bodyReader := bytes.NewReader(jsonBody)
	postResponse, postErr := http.Post(url+"/dummies/create_dummy", "application/json", bodyReader)
	assert.Nil(t, postErr)
	resBody, bodyErr := ioutil.ReadAll(postResponse.Body)
	assert.Nil(t, bodyErr)

	var dummy testrpc.Dummy
	jsonErr := json.Unmarshal(resBody, &dummy)

	assert.Nil(t, jsonErr)
	assert.NotNil(t, dummy)
	assert.Equal(t, dummy.Content, _dummy1.Content)
	assert.Equal(t, dummy.Key, _dummy1.Key)

	// 						dummy1 = dummy;

	// 		// Create another dummy
	// 			(callback) => {
	// 				rest.post("/dummy/create_dummy",
	// 					{
	// 						dummy: _dummy2
	// 					},
	// 					(err, req, res, dummy) => {
	// 						assert.isNull(err);

	// 						assert.isObject(dummy);
	// 						assert.equal(dummy.content, _dummy2.content);
	// 						assert.equal(dummy.key, _dummy2.key);

	// 						dummy2 = dummy;

	// 						callback();
	// 					}
	// 				);
	// 			},
	// 		// Get all dummies
	// 			(callback) => {
	// 				rest.post("/dummy/get_dummies",
	// 					null,
	// 					(err, req, res, dummies) => {
	// 						assert.isNull(err);

	// 						assert.isObject(dummies);
	// 						assert.lengthOf(dummies.data, 2);

	// 						callback();
	// 					}
	// 				);
	// 			},
	// 		// Update the dummy
	// 			(callback) => {
	// 				dummy1.content = "Updated Content 1";
	// 				rest.post("/dummy/update_dummy",
	// 					{
	// 						dummy: dummy1
	// 					},
	// 					(err, req, res, dummy) => {
	// 						assert.isNull(err);

	// 						assert.isObject(dummy);
	// 						assert.equal(dummy.content, "Updated Content 1");
	// 						assert.equal(dummy.key, _dummy1.key);

	// 						dummy1 = dummy;

	// 						callback();
	// 					}
	// 				);
	// 			},
	// 		// Delete dummy
	// 			(callback) => {
	// 				rest.post("/dummy/delete_dummy",
	// 					{
	// 						dummy_id: dummy1.id
	// 					},
	// 					(err, req, res) => {
	// 						assert.isNull(err);

	// 						callback();
	// 					}
	// 				);
	// 			},
	// 		// Try to get delete dummy
	// 			(callback) => {
	// 				rest.post("/dummy/get_dummy_by_id",
	// 					{
	// 						dummy_id: dummy1.id
	// 					},
	// 					(err, req, res, dummy) => {
	// 						assert.isNull(err);

	// 						// assert.isObject(dummy);

	// 						callback();
	// 					}
	// 				);
	// 			}

}
