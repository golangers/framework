package validate

import (
	"fmt"
	"golanger.com/framework/log"
	"regexp"
	"runtime"
	"strconv"
)

type ValidationError struct {
	Message, Key string
}

func (e *ValidationError) String() string {
	if e == nil {
		return ""
	}
	return e.Message
}

type Validation struct {
	Errors []*ValidationError
	keep   bool
}

func (v *Validation) Keep() {
	v.keep = true
}

func (v *Validation) Clear() {
	v.Errors = []*ValidationError{}
}

func (v *Validation) HasErrors() bool {
	return len(v.Errors) > 0
}

func (v *Validation) ErrorMap() map[string]*ValidationError {
	m := map[string]*ValidationError{}
	for _, e := range v.Errors {
		if _, ok := m[e.Key]; !ok {
			m[e.Key] = e
		}
	}
	return m
}

func (v *Validation) Error(message string, args ...interface{}) *ValidationResult {
	return (&ValidationResult{
		Ok:    false,
		Error: &ValidationError{},
	}).Message(message, args)
}

type ValidationResult struct {
	Error *ValidationError
	Ok    bool
}

func (r *ValidationResult) Key(key string) *ValidationResult {
	if r.Error != nil {
		r.Error.Key = key
	}
	return r
}

func (r *ValidationResult) Message(message string, args ...interface{}) *ValidationResult {
	if r.Error != nil {
		if len(args) == 0 {
			r.Error.Message = message
		} else {
			r.Error.Message = fmt.Sprintf(message, args)
		}
	}
	return r
}

func (v *Validation) apply(chk Validator, obj interface{}) *ValidationResult {
	if chk.IsSatisfied(obj) {
		return &ValidationResult{Ok: true}
	}

	//e.g key: validate.(*Validation).apply#84
	var key string
	if pc, _, line, ok := runtime.Caller(2); ok {
		f := runtime.FuncForPC(pc)
		key = f.Name() + "#" + strconv.Itoa(line)
	} else {
		log.Info("Failed to get Caller information to look up Validation key")
	}

	err := &ValidationError{
		Message: chk.DefaultMessage(),
		Key:     key,
	}
	v.Errors = append(v.Errors, err)

	return &ValidationResult{
		Ok:    false,
		Error: err,
	}
}

func (v *Validation) Check(obj interface{}, checks ...Validator) *ValidationResult {
	var result *ValidationResult
	for _, check := range checks {
		result = v.apply(check, obj)
		if !result.Ok {
			return result
		}
	}
	return result
}

func (v *Validation) Required(obj interface{}) *ValidationResult {
	return v.apply(Required{}, obj)
}

func (v *Validation) Min(n int, min int) *ValidationResult {
	return v.apply(Min{min}, n)
}

func (v *Validation) Max(n int, max int) *ValidationResult {
	return v.apply(Max{max}, n)
}

func (v *Validation) Range(n, min, max int) *ValidationResult {
	return v.apply(Range{Min{min}, Max{max}}, n)
}

func (v *Validation) MinSize(obj interface{}, min int) *ValidationResult {
	return v.apply(MinSize{min}, obj)
}

func (v *Validation) MaxSize(obj interface{}, max int) *ValidationResult {
	return v.apply(MaxSize{max}, obj)
}

func (v *Validation) Length(obj interface{}, n int) *ValidationResult {
	return v.apply(Length{n}, obj)
}

func (v *Validation) Match(str string, regex *regexp.Regexp) *ValidationResult {
	return v.apply(Match{regex}, str)
}

func (v *Validation) Email(str string) *ValidationResult {
	return v.apply(Email{Match{emailPattern}}, str)
}
