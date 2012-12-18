package validate

type Validator interface {
	IsSatisfied(interface{}) bool
	DefaultMessage() string
}
