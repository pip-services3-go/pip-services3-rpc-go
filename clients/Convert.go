package clients

import (
	"encoding/json"
	"reflect"
)

func ConvertComandResult(comRes interface{}, prototype reflect.Type) (convRes interface{}, err error) {

	str, ok := comRes.([]byte)
	if !ok || string(str) == "null" {
		return nil, nil
	}

	if prototype.Kind() == reflect.Ptr {
		prototype = prototype.Elem()
	}
	convRes = reflect.New(prototype).Interface()

	convErr := json.Unmarshal(comRes.([]byte), &convRes)
	if convErr != nil {
		return nil, convErr
	}

	return convRes, nil
}
