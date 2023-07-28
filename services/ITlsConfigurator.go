package services

import (
	"crypto/tls"
	"crypto/x509"
)

type ITlsConfigurator interface {
	GetClientAuthType() tls.ClientAuthType
	GetCertificates() ([]tls.Certificate, error)
	GetCaCert() (*x509.CertPool, error)
}
