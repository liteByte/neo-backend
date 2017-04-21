package structval

import (
	"reflect"
	"fmt"
)

func Validate(c interface{}) (err error) {

	v := reflect.ValueOf(c).Elem()
	t := reflect.TypeOf(c).Elem()
	for i := 0; i < v.NumField(); i++ {
		vField := v.Field(i)
		tField := t.Field(i)
		if v.Field(i).CanSet() {
			if equalsZero(vField) {
				if val, ok := tField.Tag.Lookup("default"); ok {
					vField.Set(reflect.ValueOf(val))
				} else if val, ok := tField.Tag.Lookup("required"); ok {
					if val == "true" {
						return fmt.Errorf("Field %s is required", tField.Name)
					}
				}
			}
		}

	}

	return nil
}

func equalsZero(x reflect.Value) bool {
	return x.Interface() == reflect.Zero(x.Type()).Interface()
}
