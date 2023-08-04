package main

import (
	"fmt"

	eendpoint "github.com/pip-services3-go/pip-services3-rpc-go/example/custom_endpoint/endpoint"
	elogic "github.com/pip-services3-go/pip-services3-rpc-go/example/custom_endpoint/logic"
	eservices "github.com/pip-services3-go/pip-services3-rpc-go/example/custom_endpoint/services"

	eclients "github.com/pip-services3-go/pip-services3-rpc-go/example/custom_endpoint/clients"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
)

func main() {
	service := BuildRestService()
	client := BuildRestClient()

	err := service.Open("")
	if err != nil {
		panic(err)
	}
	defer service.Close("")
	defer service.Endpoint.Close("")

	err = client.Open("")
	if err != nil {
		panic(err)
	}
	defer client.Close("")

	res, err := client.SayHello("example", "Jhon")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Response from Server: %s", res)
}

const (
	Port     = 3000
	Protocol = "https"
	Host     = "localhost"
)

func BuildRestService() *eservices.MyRestService {
	// Create custom endpoint
	endpointConfig := cconf.NewConfigParamsFromTuples(
		"root_path", "",
		"cors_headers", "x-session-id",

		"connection.protocol", Protocol,
		"connection.host", Host,
		"connection.port", Port,

		"endpoint_config.ssl_key_file", "../certs/server.key",
		"endpoint_config.ssl_crt_file", "../certs/server.crt",
		"endpoint_config.ssl_ca_file", "../certs/ca.crt",

		"endpoint_config.client_auth_type", "require_and_verify_client_cert",
	)

	eendpoint := eendpoint.NewMyHttpEndpoint()
	eendpoint.Configure(endpointConfig)

	// Create Service
	serviceConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", Protocol,
		"connection.host", Host,
		"connection.port", Port,
		"openapi_content", "swagger yaml or json content",
		"swagger.enable", "true",
	)

	ctrl := elogic.NewMyController()

	service := eservices.NewMyRestService()
	service.Configure(serviceConfig)

	var references *cref.References = cref.NewReferencesFromTuples(
		cref.NewDescriptor("pip-services-dummies", "controller", "default", "default", "1.0"), ctrl,
		cref.NewDescriptor("pip-services-dummies", "service", "rest", "default", "1.0"), service,
		cref.NewDescriptor("pip-services-dummies", "endpoint", "http", "custom", "1.0"), eendpoint,
	)

	eendpoint.SetReferences(references)
	service.SetReferences(references)

	eendpoint.Open("")
	return service
}

func BuildRestClient() *eclients.MyRestClient {
	clientConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", Protocol,
		"connection.host", Host,
		"connection.port", Port,

		"rest_client_config.correlation_id_place", "headers",

		"rest_client_config.ssl_key_file", "../certs/client.key",
		"rest_client_config.ssl_crt_file", "../certs/client.crt",
		"rest_client_config.ssl_ca_file", "../certs/ca.crt",
	)

	client := eclients.NewMyRestClient()

	client.Configure(clientConfig)
	client.SetReferences(cref.NewEmptyReferences())

	return client
}
