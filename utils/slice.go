package utils

import (
	"reflect"
)

type Slice struct {
	I interface{}
}

// convert slice to map
func (s Slice) SliceToMap(is ...interface{}) M {
	m := M{}
	var inter []interface{}
	var interNum int
	switch si := s.I.(type) {
	case []string:
		if len(is) > 0 {
			var ok bool
			inter, ok = is[0].([]interface{})
			if !ok {
				return m
			}
		}

		interNum = len(inter)
		for key, value := range si {
			if key <= interNum {
				v := reflect.ValueOf(inter[key])
				if v.Kind() == reflect.Ptr {
					v = v.Elem()
				}

				m[value] = v.Interface()
			} else {
				m[value] = nil
			}
		}
	}

	return m
}
