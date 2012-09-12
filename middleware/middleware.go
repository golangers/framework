package middleware

import (
	"fmt"
)

var Middleware = middleware{}

type middleware map[string]interface{}

func (m middleware) Add(key string, value interface{}) {
	if _, ok := m[key]; !ok {
		m[key] = value
	} else {
		fmt.Println("key:" + key + " already exists")
	}
}

func (m middleware) Set(key string, value interface{}) {
	if _, ok := m[key]; ok {
		m[key] = value
	} else {
		fmt.Println("key:" + key + " does not exists")
	}
}

func (m middleware) Get(key string) interface{} {
	return m[key]
}

func (m middleware) Del(key string) {
	delete(m, key)
}
