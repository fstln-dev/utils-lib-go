package http

import (
	"fmt"
	"net/url"
	"reflect"
)

func StructToUrl(v interface{}, data *url.Values) (err error) {

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Panic: %v", r)
		}
	}()

	rType, rVal := reflect.TypeOf(v), reflect.ValueOf(v)
	if rType.Kind() == reflect.Ptr {
		rType, rVal = rType.Elem(), rVal.Elem()
	} else {
		panic("v must be ptr to struct")
	}

	for i := 0; i < rType.NumField(); i++ {
		name := rType.Field(i).Tag.Get("json")
		if name == "" {
			name = rType.Field(i).Name
		}
		// 偷懒
		data.Set(name, fmt.Sprint(rVal.Field(i)))
	}
	return nil
}
