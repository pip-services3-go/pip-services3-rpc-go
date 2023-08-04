package logic

type MyController struct {
}

func NewMyController() *MyController {
	dc := MyController{}
	return &dc
}

func (c *MyController) SayHello(correlationId string, name string) (result string, err error) {
	if name == "" {
		name = "user"
	}
	return "Hello " + name + "!", nil
}
