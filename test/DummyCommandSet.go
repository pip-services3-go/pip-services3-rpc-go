package test_rpc

import (
	ccomand "github.com/pip-services3-go/pip-services3-commons-go/commands"
	cconv "github.com/pip-services3-go/pip-services3-commons-go/convert"
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	crun "github.com/pip-services3-go/pip-services3-commons-go/run"
	cvalid "github.com/pip-services3-go/pip-services3-commons-go/validate"
)

type DummyCommandSet struct {
	ccomand.CommandSet
	controller IDummyController
}

func NewDummyCommandSet(controller IDummyController) *DummyCommandSet {
	//super();
	dcs := DummyCommandSet{}
	dcs.CommandSet = *ccomand.NewCommandSet()

	dcs.controller = controller

	dcs.AddCommand(dcs.makeGetPageByFilterCommand())
	dcs.AddCommand(dcs.makeGetOneByIdCommand())
	dcs.AddCommand(dcs.makeCreateCommand())
	dcs.AddCommand(dcs.makeUpdateCommand())
	dcs.AddCommand(dcs.makeDeleteByIdCommand())
	return &dcs
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
		cvalid.NewObjectSchema().WithRequiredProperty("dummy", NewDummySchema()),
		func(correlationId string, args *crun.Parameters) (result interface{}, err error) {
			var entity Dummy = args.Get("dummy").(Dummy)
			return c.controller.Create(correlationId, entity)
		},
	)
}

func (c *DummyCommandSet) makeUpdateCommand() ccomand.ICommand {
	return ccomand.NewCommand(
		"update_dummy",
		cvalid.NewObjectSchema().WithRequiredProperty("dummy", NewDummySchema()),
		func(correlationId string, args *crun.Parameters) (result interface{}, err error) {
			var entity Dummy = args.Get("dummy").(Dummy)
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
