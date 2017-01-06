package validator

import (
	"reflect"
	"log"
)

func Validate(c interface{}) {

	isZero := func(x reflect.Value) bool {
		return x.Interface() == reflect.Zero(x.Type()).Interface()
	}

	v := reflect.ValueOf(c).Elem()
	t := reflect.TypeOf(c).Elem()
	for i := 0; i < v.NumField(); i++ {
		vField := v.Field(i)
		tField := t.Field(i)
		if isZero(vField) {
			if val, ok := tField.Tag.Lookup("default"); ok {
				vField.Set(reflect.ValueOf(val))
			} else if val, ok := tField.Tag.Lookup("required"); ok {
				if val == "true" {
					log.Fatalf("Field %s is required", tField.Name)
				}
			}

		}
	}
}
