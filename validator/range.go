package validator

import (
	"fmt"
)

// Requires an integer to be within Min, Max inclusive.
type Range struct {
	Min
	Max
}

func (r Range) IsSatisfied(obj interface{}) bool {
	return r.Min.IsSatisfied(obj) && r.Max.IsSatisfied(obj)
}

func (r Range) DefaultMessage() string {
	return fmt.Sprintln("Range is", r.Min.Min, "to", r.Max.Max)
}
