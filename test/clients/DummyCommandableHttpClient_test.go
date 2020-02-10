package test_rpc_clients

// import (
// 	"testing"
// )

// func TestDummyCommandableHttpClient(t *testing.T) {

// 	 restConfig := cconf.NewConfigParamsFromTuples(
// 		"connection.protocol", "http",
// 		"connection.host", "localhost",
// 		"connection.port", "3000"
// 	);

// 		var service* DummyCommandableHttpService;
// 		var client* DummyCommandableHttpClient;

// 	 	var fixture *DummyClientFixture;

// 			ctrl := NewDummyController();

// 			service = NewDummyCommandableHttpService();
// 			service.Configure(restConfig);

// 			references = cref.NewReferencesFromTuples(
// 				cref.NewDescriptor("pip-services-dummies", "controller", "default", "default", "1.0"), ctrl,
// 				cref.NewDescriptor("pip-services-dummies", "service", "http", "default", "1.0"), service,
// 			);
// 			service.SetReferences(references);

// 			service.Open("");
// 			defer service.Close("")

// 			client = NewDummyCommandableHttpClient();
// 			fixture = NewDummyClientFixture(client);

// 			client.Configure(restConfig);
// 			client.SetReferences(cref.NewReferences());
// 			client.Open("");
// 	t.Run( "CRUD Operations",		fixture.testCrudOperations)
// }
