package example_clients

import (
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	edata "github.com/pip-services3-go/pip-services3-rpc-go/example/basic_tls/data"
)

type IDummyClient interface {
	GetDummies(correlationId string, filter *cdata.FilterParams, paging *cdata.PagingParams) (result *edata.DummyDataPage, err error)
	GetDummyById(correlationId string, dummyId string) (result *edata.Dummy, err error)
	CreateDummy(correlationId string, dummy edata.Dummy) (result *edata.Dummy, err error)
	UpdateDummy(correlationId string, dummy edata.Dummy) (result *edata.Dummy, err error)
	DeleteDummy(correlationId string, dummyId string) (result *edata.Dummy, err error)

	CheckCorrelationId(correlationId string) (result map[string]string, err error)

	CheckErrorPropagation(correlationId string) error
}
