package utils

import (
	"encoding/gob"
	"reflect"
	"strconv"
)

func init() {
	gob.Register([]M{})
	gob.Register(M{})
}

type M map[string]interface{}

// convert map to struct
func (m M) MapToStruct(s interface{}) {
	v := reflect.Indirect(reflect.ValueOf(s))

	for i := 0; i < v.NumField(); i++ {
		f := v.Type().Field(i)
		key := f.Name
		scnKey := Strings(key).SnakeCasedName()
		tag := f.Tag
		fieldName := tag.Get("field")
		vf := v.Field(i)
		doRes := false
		if fieldName != "" {
			if val, ok := m[fieldName]; ok {
				vv := reflect.ValueOf(val)
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

				doRes = true
			}
		}

		if !doRes {
			if _, ok := m[key]; ok {
				vf.Set(reflect.ValueOf(m[key]))
			}
		}

		if !doRes {
			if _, ok := m[scnKey]; ok {
				vf.Set(reflect.ValueOf(m[scnKey]))
			}
		}
	}
}
