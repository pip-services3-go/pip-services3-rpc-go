package endpoint

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"strings"

	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cauth "github.com/pip-services3-go/pip-services3-components-go/auth"
	cserv "github.com/pip-services3-go/pip-services3-rpc-go/services"
)

const (
	ENDPOINT_CONFIG  = "endpoint_config"
	CLIENT_AUTH_TYPE = "client_auth_type"

	CERTS_CA  = "ssl_ca_file"
	CERTS_CRT = "ssl_crt_file"
	CERTS_KEY = "ssl_key_file"
)

type MyHttpEndpoint struct {
	*cserv.HttpEndpoint
	credentialStore *cauth.MemoryCredentialStore
}

func NewMyHttpEndpoint() *MyHttpEndpoint {
	c := &MyHttpEndpoint{}
	c.HttpEndpoint = cserv.InheritTlsHttpEndpoint(c)
	c.credentialStore = cauth.NewEmptyMemoryCredentialStore()

	return c
}

func (c *MyHttpEndpoint) Configure(config *cconf.ConfigParams) {
	c.HttpEndpoint.Configure(config)
	c.credentialStore.Configure(config)
}

func (c *MyHttpEndpoint) GetClientAuthType() tls.ClientAuthType {
	config, err := c.credentialStore.Lookup("", ENDPOINT_CONFIG)
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
func (c *MyHttpEndpoint) GetCertificates() ([]tls.Certificate, error) {
	config, err := c.credentialStore.Lookup("", ENDPOINT_CONFIG)
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
func (c *MyHttpEndpoint) GetCaCert() (*x509.CertPool, error) {
	config, err := c.credentialStore.Lookup("", ENDPOINT_CONFIG)
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
