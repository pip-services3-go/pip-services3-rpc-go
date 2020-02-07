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

func TestDummyRestService(t *testing.T) {
	restConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "http",
		"connection.host", "localhost",
		"connection.port", "3000",
	)

	var _dummy1 testrpc.Dummy
	var _dummy2 testrpc.Dummy
	var service *DummyRestService
	ctrl := testrpc.NewDummyController()

	service = NewDummyRestService()
	service.Configure(restConfig)

	var references *cref.References = cref.NewReferencesFromTuples(
		cref.NewDescriptor("pip-services-dummies", "controller", "default", "default", "1.0"), ctrl,
		cref.NewDescriptor("pip-services-dummies", "service", "rest", "default", "1.0"), service,
	)
	service.SetReferences(references)
	opnErr := service.Open("")
	assert.Nil(t, opnErr)
	defer service.Close("")

	url := "http://localhost:3000"

	_dummy1 = testrpc.Dummy{Id: "", Key: "Key 1", Content: "Content 1"}
	_dummy2 = testrpc.Dummy{Id: "", Key: "Key 2", Content: "Content 2"}

	//var dummy1, dummy2 testrpc.Dummy

	// Create one dummy
	jsonBody, _ := json.Marshal(_dummy1)

	bodyReader := bytes.NewReader(jsonBody)
	postResponse, postErr := http.Post(url+"/dummies", "application/json", bodyReader)
	assert.Nil(t, postErr)
	resBody, bodyErr := ioutil.ReadAll(postResponse.Body)
	assert.Nil(t, bodyErr)

	var dummy testrpc.Dummy
	jsonErr := json.Unmarshal(resBody, &dummy)

	assert.Nil(t, jsonErr)
	assert.NotNil(t, dummy)
	assert.Equal(t, dummy.Content, _dummy1.Content)
	assert.Equal(t, dummy.Key, _dummy1.Key)

	//dummy1 = dummy

	// Create another dummy
	jsonBody, _ = json.Marshal(_dummy2)

	bodyReader = bytes.NewReader(jsonBody)
	postResponse, postErr = http.Post(url+"/dummies", "application/json", bodyReader)
	assert.Nil(t, postErr)
	resBody, bodyErr = ioutil.ReadAll(postResponse.Body)
	assert.Nil(t, bodyErr)

	jsonErr = json.Unmarshal(resBody, &dummy)

	assert.Nil(t, jsonErr)
	assert.NotNil(t, dummy)
	assert.Equal(t, dummy.Content, _dummy2.Content)
	assert.Equal(t, dummy.Key, _dummy2.Key)
	//dummy2 = dummy

	// Get all dummies
	postResponse, postErr = http.Get(url + "/dummies")

	resBody, bodyErr = ioutil.ReadAll(postResponse.Body)
	assert.Nil(t, bodyErr)
	var dummies testrpc.DummyDataPage
	jsonErr = json.Unmarshal(resBody, &dummies)
	assert.Nil(t, postErr)
	assert.NotNil(t, dummies)
	assert.Len(t, dummies.Data, 2)

	//         // Update the dummy
	//             (callback) => {
	//                 dummy1.content = "Updated Content 1";
	//                 rest.put("/dummies",
	//                     dummy1,
	//                     (err, req, res, dummy) => {
	//                         assert.isNull(err);

	//                         assert.isObject(dummy);
	//                         assert.equal(dummy.content, "Updated Content 1");
	//                         assert.equal(dummy.key, _dummy1.key);

	//                         dummy1 = dummy;

	//                         callback();
	//                     }
	//                 );
	//             },
	//         // Delete dummy
	//             (callback) => {
	//                 rest.del("/dummies/" + dummy1.id,
	//                     (err, req, res) => {
	//                         assert.isNull(err);

	//                         callback();
	//                     }
	//                 );
	//             },
	//         // Try to get delete dummy
	//             (callback) => {
	//                 rest.get("/dummies/" + dummy1.id,
	//                     (err, req, res, dummy) => {
	//                         assert.isNull(err);

	//                         // assert.isObject(dummy);

	//                         callback();
	//                     }
	//                 );
	//             }
	//         ], done);
	//     });

	// });
}
