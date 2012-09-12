package utils

import (
	"reflect"
	"strconv"
)

type M map[string]interface{}

// convert map to struct
func (m M) MapToStruct(s interface{}) {
	v := reflect.Indirect(reflect.ValueOf(s))

	for i := 0; i < v.NumField(); i++ {
		f := v.Type().Field(i)
		key := f.Name
		tag := f.Tag
		fieldName := tag.Get("field")
		vf := v.Field(i)
		if _, ok := m[key]; ok {
			vf.Set(reflect.ValueOf(m[key]))
		} else if fieldName != "" {
			vv := reflect.ValueOf(m[fieldName])
			if vf.Type().Kind().String() != vv.Type().Kind().String() {
				if vf.Type().Kind().String() == "bool" {
					if vv.Type().Kind().String() == "int64" && vv.Int() > 0 {
						vf.SetBool(true)
					}

					if vv.Type().Kind().String() == "string" && vv.String() != "" {
						ii, _ := strconv.ParseInt(vv.String(), 10, 64)
						if ii > 0 {
							vf.SetBool(true)
						}
					}
				}
			} else {
				vf.Set(vv)
			}
		}
	}
}
