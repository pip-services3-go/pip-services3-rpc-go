package example_logic

import (
	"encoding/json"

	ccomand "github.com/pip-services3-go/pip-services3-commons-go/commands"
	cconv "github.com/pip-services3-go/pip-services3-commons-go/convert"
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	crun "github.com/pip-services3-go/pip-services3-commons-go/run"
	cvalid "github.com/pip-services3-go/pip-services3-commons-go/validate"
	edata "github.com/pip-services3-go/pip-services3-rpc-go/example/basic_tls/data"
)

type DummyCommandSet struct {
	ccomand.CommandSet
	controller IDummyController
}

func NewDummyCommandSet(controller IDummyController) *DummyCommandSet {
	c := DummyCommandSet{}
	c.CommandSet = *ccomand.NewCommandSet()

	c.controller = controller

	c.AddCommand(c.makeGetPageByFilterCommand())
	c.AddCommand(c.makeGetOneByIdCommand())
	c.AddCommand(c.makeCreateCommand())
	c.AddCommand(c.makeUpdateCommand())
	c.AddCommand(c.makeDeleteByIdCommand())
	c.AddCommand(c.makeCheckCorrelationIdCommand())
	c.AddCommand(c.makeCheckErrorPropagationCommand())
	return &c
}

func (c *DummyCommandSet) makeGetPageByFilterCommand() ccomand.ICommand {
	return ccomand.NewCommand(
		"get_dummies",
		cvalid.NewObjectSchema().WithOptionalProperty("filter", cvalid.NewFilterParamsSchema()).WithOptionalProperty("paging", cvalid.NewPagingParamsSchema()),
		func(correlationId string, args *crun.Parameters) (result interface{}, err error) {
			filter := cdata.NewFilterParamsFromValue(args.Get("filter"))
			paging := cdata.NewPagingParamsFromValue(args.Get("paging"))
			return c.controller.GetPageByFilter(correlationId, filter, paging)
		},
	)
}

func (c *DummyCommandSet) makeGetOneByIdCommand() ccomand.ICommand {
	return ccomand.NewCommand(
		"get_dummy_by_id",
		cvalid.NewObjectSchema().WithRequiredProperty("dummy_id", cconv.String),
		func(correlationId string, args *crun.Parameters) (result interface{}, err error) {
			id := args.GetAsString("dummy_id")
			return c.controller.GetOneById(correlationId, id)
		},
	)
}

func (c *DummyCommandSet) makeCreateCommand() ccomand.ICommand {
	return ccomand.NewCommand(
		"create_dummy",
		cvalid.NewObjectSchema().WithRequiredProperty("dummy", edata.NewDummySchema()),
		func(correlationId string, args *crun.Parameters) (result interface{}, err error) {
			val, _ := json.Marshal(args.Get("dummy"))
			var entity edata.Dummy
			json.Unmarshal(val, &entity)

			return c.controller.Create(correlationId, entity)
		},
	)
}

func (c *DummyCommandSet) makeUpdateCommand() ccomand.ICommand {
	return ccomand.NewCommand(
		"update_dummy",
		cvalid.NewObjectSchema().WithRequiredProperty("dummy", edata.NewDummySchema()),
		func(correlationId string, args *crun.Parameters) (result interface{}, err error) {
			val, _ := json.Marshal(args.Get("dummy"))
			var entity edata.Dummy
			json.Unmarshal(val, &entity)
			return c.controller.Update(correlationId, entity)
		},
	)
}

func (c *DummyCommandSet) makeDeleteByIdCommand() ccomand.ICommand {
	return ccomand.NewCommand(
		"delete_dummy",
		cvalid.NewObjectSchema().WithRequiredProperty("dummy_id", cconv.String),
		func(correlationId string, args *crun.Parameters) (result interface{}, err error) {
			id := args.GetAsString("dummy_id")
			return c.controller.DeleteById(correlationId, id)
		},
	)
}

func (c *DummyCommandSet) makeCheckCorrelationIdCommand() ccomand.ICommand {
	return ccomand.NewCommand(
		"check_correlation_id",
		cvalid.NewObjectSchema(),
		func(correlationId string, args *crun.Parameters) (result interface{}, err error) {
			return c.controller.CheckCorrelationId(correlationId)
		},
	)
}

func (c *DummyCommandSet) makeCheckErrorPropagationCommand() ccomand.ICommand {
	return ccomand.NewCommand(
		"check_error_propagation",
		cvalid.NewObjectSchema(),
		func(correlationId string, args *crun.Parameters) (result interface{}, err error) {
			return nil, c.controller.CheckErrorPropagation(correlationId)
		},
	)
}
