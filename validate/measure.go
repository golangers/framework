package validate

import (
	"fmt"
)

type Min struct {
	Min int
}

func (m Min) IsSatisfied(obj interface{}) bool {
	num, ok := obj.(int)
	if ok {
		return num >= m.Min
	}
	return false
}

func (m Min) DefaultMessage() string {
	return fmt.Sprint("Minimum is", m.Min)
}

type Max struct {
	Max int
}

func (m Max) IsSatisfied(obj interface{}) bool {
	num, ok := obj.(int)
	if ok {
		return num <= m.Max
	}
	return false
}

func (m Max) DefaultMessage() string {
	return fmt.Sprint("Maximum is", m.Max)
}

type Range struct {
	Min
	Max
}

func (r Range) IsSatisfied(obj interface{}) bool {
	return r.Min.IsSatisfied(obj) && r.Max.IsSatisfied(obj)
}

func (r Range) DefaultMessage() string {
	return fmt.Sprint("Range is", r.Min.Min, "to", r.Max.Max)
}

type MinSize struct {
	Min int
}

func (m MinSize) IsSatisfied(obj interface{}) bool {
	if arr, ok := obj.([]interface{}); ok {
		return len(arr) >= m.Min
	}
	if str, ok := obj.(string); ok {
		return len(str) >= m.Min
	}
	return false
}

func (m MinSize) DefaultMessage() string {
	return fmt.Sprint("Minimum size is", m.Min)
}

type MaxSize struct {
	Max int
}

func (m MaxSize) IsSatisfied(obj interface{}) bool {
	if arr, ok := obj.([]interface{}); ok {
		return len(arr) <= m.Max
	}
	if str, ok := obj.(string); ok {
		return len(str) <= m.Max
	}
	return false
}

func (m MaxSize) DefaultMessage() string {
	return fmt.Sprint("Maximum size is", m.Max)
}

type Length struct {
	N int
}

func (s Length) IsSatisfied(obj interface{}) bool {
	if arr, ok := obj.([]interface{}); ok {
		return len(arr) == s.N
	}
	if str, ok := obj.(string); ok {
		return len(str) == s.N
	}
	return false
}

func (s Length) DefaultMessage() string {
	return fmt.Sprint("Required length is", s.N)
}
