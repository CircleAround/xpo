package domain

import (
	"local/validatekit"
)

func NewReportValidate() *validatekit.Validate {
	v := validatekit.NewValidate()
	v.RegisterValidation("languages", NewLanguagesValidation())
	return v
}
