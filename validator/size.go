package validator

import (
	"fmt"
)

// Requires an array or string to be at least a given length.
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
	return fmt.Sprintln("Minimum size is", m.Min)
}

// Requires an array or string to be at most a given length.
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
	return fmt.Sprintln("Maximum size is", m.Max)
}

// Requires an array or string to be exactly a given length.
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
	return fmt.Sprintln("Required length is", s.N)
}
