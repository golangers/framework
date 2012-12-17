package validator

type Validator interface {
	IsSatisfied(interface{}) bool
	DefaultMessage() string
}
