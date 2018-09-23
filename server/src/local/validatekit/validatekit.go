package validatekit

import (
	"regexp"

	validator "gopkg.in/go-playground/validator.v9"
)

func NewRegexValidation(regex string) func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		reg := regexp.MustCompile(regex)
		return reg.MatchString(fl.Field().String())
	}
}

type Validate struct {
	validator.Validate
}

func NewValidate() *Validate {
	vv := validator.New()
	return &Validate{*vv}
}

func (v *Validate) RegisterRegexValidation(tag string, regex string) {
	v.RegisterValidation(tag, NewRegexValidation(regex))
}
