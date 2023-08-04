package data

import (
	cconv "github.com/pip-services3-go/pip-services3-commons-go/convert"
	cvalid "github.com/pip-services3-go/pip-services3-commons-go/validate"
)

type HelloV1Schema struct {
	cvalid.ObjectSchema
}

func NewHelloV1Schema() *HelloV1Schema {
	ds := HelloV1Schema{}
	ds.ObjectSchema = *cvalid.NewObjectSchema()

	ds.WithOptionalProperty("name", cconv.String)
	return &ds
}
