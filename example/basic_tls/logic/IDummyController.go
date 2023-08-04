package example_logic

import (
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	edata "github.com/pip-services3-go/pip-services3-rpc-go/example/basic_tls/data"
)

type IDummyController interface {
	GetPageByFilter(correlationId string, filter *cdata.FilterParams, paging *cdata.PagingParams) (result *edata.DummyDataPage, err error)
	GetOneById(correlationId string, id string) (result *edata.Dummy, err error)
	Create(correlationId string, entity edata.Dummy) (result *edata.Dummy, err error)
	Update(correlationId string, entity edata.Dummy) (result *edata.Dummy, err error)
	DeleteById(correlationId string, id string) (result *edata.Dummy, err error)

	CheckCorrelationId(correlationId string) (result map[string]string, err error)

	CheckErrorPropagation(correlationId string) error
}
