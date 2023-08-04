package main

import (
	"fmt"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	eclients "github.com/pip-services3-go/pip-services3-rpc-go/example/basic_tls/clients"
	edata "github.com/pip-services3-go/pip-services3-rpc-go/example/basic_tls/data"
	elogic "github.com/pip-services3-go/pip-services3-rpc-go/example/basic_tls/logic"
	eservices "github.com/pip-services3-go/pip-services3-rpc-go/example/basic_tls/services"
)

func main() {
	service := BuildRestService()
	client := BuildRestClient()

	err := service.Open("")
	if err != nil {
		panic(err)
	}
	defer service.Close("")

	err = client.Open("")
	if err != nil {
		panic(err)
	}
	defer client.Close("")

	res, err := client.CreateDummy("example", edata.Dummy(*edata.NewDummy("1", "my_key", "some content")))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Created Dummy: id=%s, key=%s, content=%s", res.Id, res.Key, res.Content)
}

const (
	Port     = 3000
	Protocol = "https"
	Host     = "localhost"
)

func BuildRestService() *eservices.DummyRestService {

	serviceConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", Protocol,
		"connection.host", Host,
		"connection.port", Port,
		"openapi_content", "swagger yaml or json content",
		"swagger.enable", "true",

		"options.client_auth_type", "require_and_verify_client_cert",
		"options.certificate_server_name", "localhost",

		"credential.ssl_key_file", "../certs/server.key",
		"credential.ssl_crt_file", "../certs/server.crt",
		"credential.ssl_ca_file", "../certs/ca.crt",
	)

	var service *eservices.DummyRestService
	ctrl := elogic.NewDummyController()

	service = eservices.NewDummyRestService()
	service.Configure(serviceConfig)

	var references *cref.References = cref.NewReferencesFromTuples(
		cref.NewDescriptor("pip-services-dummies", "controller", "default", "default", "1.0"), ctrl,
		cref.NewDescriptor("pip-services-dummies", "service", "rest", "default", "1.0"), service,
	)
	service.SetReferences(references)

	return service
}

func BuildRestClient() *eclients.DummyRestClient {
	clientConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", Protocol,
		"connection.host", Host,
		"connection.port", Port,

		"options.correlation_id_place", "headers",
		"options.certificate_server_name", "localhost",

		"credential.ssl_key_file", "../certs/client.key",
		"credential.ssl_crt_file", "../certs/client.crt",
		"credential.ssl_ca_file", "../certs/ca.crt",
	)

	client := eclients.NewDummyRestClient()

	client.Configure(clientConfig)
	client.SetReferences(cref.NewEmptyReferences())

	return client
}
