package test_rpc_services

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	var _dummy2 testrpc.Dummy

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
	_dummy2 = testrpc.Dummy{Id: "", Key: "Key 2", Content: "Content 2"}

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

	dummy1 := dummy

	// Create another dummy
	bodyMap = make(map[string]interface{})
	bodyMap["dummy"] = _dummy2

	jsonBody, _ = json.Marshal(bodyMap)

	bodyReader = bytes.NewReader(jsonBody)
	postResponse, postErr = http.Post(url+"/dummies/create_dummy", "application/json", bodyReader)
	assert.Nil(t, postErr)
	resBody, bodyErr = ioutil.ReadAll(postResponse.Body)
	assert.Nil(t, bodyErr)

	jsonErr = json.Unmarshal(resBody, &dummy)

	assert.Nil(t, jsonErr)
	assert.NotNil(t, dummy)
	assert.Equal(t, dummy.Content, _dummy2.Content)
	assert.Equal(t, dummy.Key, _dummy2.Key)

	// Get all dummies

	postResponse, postErr = http.Post(url+"/dummies/get_dummies", "application/json", nil)
	assert.Nil(t, postErr)
	resBody, bodyErr = ioutil.ReadAll(postResponse.Body)
	assert.Nil(t, bodyErr)
	var dummies testrpc.DummyDataPage
	jsonErr = json.Unmarshal(resBody, &dummies)
	assert.Nil(t, jsonErr)
	assert.NotNil(t, dummies)
	assert.Len(t, dummies.Data, 2)

	// Update the dummy
	dummy1.Content = "Updated Content 1"
	bodyMap = make(map[string]interface{})
	bodyMap["dummy"] = dummy1

	jsonBody, _ = json.Marshal(bodyMap)

	bodyReader = bytes.NewReader(jsonBody)
	postResponse, postErr = http.Post(url+"/dummies/update_dummy", "application/json", bodyReader)
	assert.Nil(t, postErr)
	resBody, bodyErr = ioutil.ReadAll(postResponse.Body)
	assert.Nil(t, bodyErr)
	jsonErr = json.Unmarshal(resBody, &dummy)
	assert.Nil(t, jsonErr)
	assert.NotNil(t, dummy)
	assert.Equal(t, dummy.Content, "Updated Content 1")
	assert.Equal(t, dummy.Key, _dummy1.Key)

	// Delete dummy
	bodyMap = make(map[string]interface{})
	bodyMap["dummy_id"] = dummy1.Id
	jsonBody, _ = json.Marshal(bodyMap)
	bodyReader = bytes.NewReader(jsonBody)
	postResponse, postErr = http.Post(url+"/dummies/delete_dummy", "application/json", bodyReader)
	assert.Nil(t, postErr)
	resBody, bodyErr = ioutil.ReadAll(postResponse.Body)
	assert.Nil(t, bodyErr)

	// Try to get delete dummy
	bodyMap = make(map[string]interface{})
	bodyMap["dummy_id"] = dummy1.Id
	jsonBody, _ = json.Marshal(bodyMap)
	bodyReader = bytes.NewReader(jsonBody)
	postResponse, postErr = http.Post(url+"/dummies/get_dummy_by_id", "application/json", bodyReader)
	assert.Nil(t, postErr)
	resBody, bodyErr = ioutil.ReadAll(postResponse.Body)
	assert.Nil(t, bodyErr)

	fmt.Println((string)(resBody))

	jsonErr = json.Unmarshal(resBody, &dummy)
	assert.Nil(t, jsonErr)

}