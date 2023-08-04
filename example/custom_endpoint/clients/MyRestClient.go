package clients

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"reflect"
	"strings"

	edata "github.com/pip-services3-go/pip-services3-rpc-go/example/custom_endpoint/data/version1"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cauth "github.com/pip-services3-go/pip-services3-components-go/auth"
	"github.com/pip-services3-go/pip-services3-rpc-go/clients"
)

var (
	helloType = reflect.TypeOf(edata.HelloV1{}.Name)
)

const (
	REST_CLIENT_CONFIG = "rest_client_config"
	CLIENT_AUTH_TYPE   = "client_auth_type"

	CERTS_CA  = "ssl_ca_file"
	CERTS_CRT = "ssl_crt_file"
	CERTS_KEY = "ssl_key_file"
)

type MyRestClient struct {
	*clients.RestClient
	credentialStore *cauth.MemoryCredentialStore
}

func NewMyRestClient() *MyRestClient {
	c := &MyRestClient{}
	c.credentialStore = cauth.NewEmptyMemoryCredentialStore()
	c.RestClient = clients.InheritTlsRestClient(c)
	return c
}

func (c *MyRestClient) SayHello(correlationId string, name string) (result string, err error) {
	calValue, calErr := c.Call(helloType, "post", "/hello", correlationId, nil, edata.HelloV1{Name: name})
	if calErr != nil {
		return "", calErr
	}

	val, _ := calValue.(*string)
	result = *val
	c.Instrument(correlationId, "myservice.hello")
	return result, nil
}

func (c *MyRestClient) Configure(config *cconf.ConfigParams) {
	c.RestClient.Configure(config)
	c.credentialStore.Configure(config)
}

func (c *MyRestClient) GetClientAuthType() tls.ClientAuthType {
	config, err := c.credentialStore.Lookup("", REST_CLIENT_CONFIG)
	if err != nil {
		panic("Credentials config is empty: " + err.Error())
	}

	clientAuthType := config.GetAsString(CLIENT_AUTH_TYPE)

	switch strings.ToLower(clientAuthType) {
	default:
		return tls.NoClientCert
	case "request_client_cert":
		return tls.RequestClientCert
	case "require_any_client_cert":
		return tls.RequireAnyClientCert
	case "verify_client_cert_if_given":
		return tls.VerifyClientCertIfGiven
	case "require_and_verify_client_cert":
		return tls.RequireAndVerifyClientCert
	}
}
func (c *MyRestClient) GetCertificates() ([]tls.Certificate, error) {
	config, err := c.credentialStore.Lookup("", REST_CLIENT_CONFIG)
	if err != nil {
		panic("Credentials config is empty: " + err.Error())
	}

	sslKeyFile := config.GetAsString(CERTS_KEY)
	sslCrtFile := config.GetAsString(CERTS_CRT)

	var certificates []tls.Certificate

	if sslCrtFile != "" && sslKeyFile != "" {
		certificate, err := tls.LoadX509KeyPair(sslCrtFile, sslKeyFile)
		if err != nil {
			return nil, err
		}
		certificates = append(certificates, certificate)
	}
	return certificates, nil
}
func (c *MyRestClient) GetCaCert() (*x509.CertPool, error) {
	config, err := c.credentialStore.Lookup("", REST_CLIENT_CONFIG)
	if err != nil {
		panic("Credentials config is empty: " + err.Error())
	}

	sslCaFile := config.GetAsString(CERTS_CA)

	caCertPool := x509.NewCertPool()
	if sslCaFile != "" {
		bytes, err := os.ReadFile(sslCaFile)
		if err != nil {
			return nil, err
		}
		caCertPool.AppendCertsFromPEM(bytes)

		return caCertPool, nil
	} else {
		return nil, nil
	}
}
