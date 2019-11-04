package slice

import (
	"database/sql/driver"
	"reflect"
)

// Flatten do flatten the input slice.
func Flatten(args []interface{}) []interface{} {
	ret := make([]interface{}, 0, len(args))
	for _, arg := range args {
		if driver.IsValue(arg) {
			ret = append(ret, arg)
			continue
		}

		v := reflect.ValueOf(arg)
		kind := v.Kind()
		if kind == reflect.Array || kind == reflect.Slice {
			ret = appendList(ret, v)
		} else {
			ret = append(ret, arg)
		}
	}
	return ret
}

func appendList(args []interface{}, v reflect.Value) []interface{} {
	vlen := v.Len()
	for i := 0; i < vlen; i++ {
		vv := v.Index(i)
		val := vv.Interface()

		if driver.IsValue(val) {
			args = append(args, val)
			continue
		}

		if vv.Kind() == reflect.Interface {
			vv = vv.Elem()
		}

		kind := vv.Kind()
		if kind == reflect.Array || kind == reflect.Slice {
			args = appendList(args, vv)
		} else {
			args = append(args, val)
		}
	}
	return args
}
