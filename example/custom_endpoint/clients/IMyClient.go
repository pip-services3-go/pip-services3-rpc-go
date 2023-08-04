package clients

type IMyClient interface {
	SayHello(correlationId string, name string) (result string, err error)
}
