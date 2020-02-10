package test_rpc

//import { IStringIdentifiable } from 'pip-services3-commons-node';
// IStringIdentifiable
type Dummy struct {
	Id      string `json:"id"`
	Key     string `json:"key"`
	Content string `json:"content"`
}

func NewDummy(id string, key string, content string) *Dummy {
	return &Dummy{
		Id:      id,
		Key:     key,
		Content: content,
	}
}
