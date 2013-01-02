package middleware

import (
	"fmt"
	"sync"
)

var mu sync.RWMutex

var Middleware = middleware{}

type middleware map[string]interface{}

func (m middleware) Add(key string, value interface{}) {
	mu.RLock()
	_, ok := m[key]
	mu.RUnlock()
	if !ok {
		mu.Lock()
		m[key] = value
		mu.Unlock()
	} else {
		fmt.Println("key:" + key + " already exists")
	}
}

func (m middleware) Set(key string, value interface{}) {
	mu.RLock()
	_, ok := m[key]
	mu.RUnlock()
	if ok {
		mu.Lock()
		m[key] = value
		mu.Unlock()
	} else {
		fmt.Println("key:" + key + " does not exists")
	}
}

func (m middleware) Get(key string) interface{} {
	mu.RLock()
	defer mu.RUnlock()
	return m[key]
}

func (m middleware) Del(key string) {
	mu.Lock()
	delete(m, key)
	mu.Unlock()
}
