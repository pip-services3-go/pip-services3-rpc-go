package logic

type IMyController interface {
	SayHello(correlationId string, name string) (result string, err error)
}
