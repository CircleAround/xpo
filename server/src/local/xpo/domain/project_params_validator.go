package domain

import (
	"local/validatekit"
	"local/xpo/entities"

	"github.com/pkg/errors"
)

func ValidateProject(params *entities.Project) (*validatekit.Validate, error) {
	v, err := newProjectValidator()
	if err != nil {
		return nil, errors.Wrap(err, "newProjectValidator failed")
	}

	err = v.Struct(params)
	if err != nil {
		return nil, errors.Wrap(err, "Validation failed")
	}
	return v, nil
}

func newProjectValidator() (*validatekit.Validate, error) {
	v := validatekit.NewValidate()
	v.RegisterRegexValidation("identity_name_format", `^[a-z][0-9a-z_]+$`)

	vali, err := NewReservedIdentityNameValidation()
	if err != nil {
		return nil, errors.Wrap(err, "NewReservedIdentityNameValidation failed")
	}
	v.RegisterValidation("reserved_identity_name", vali)

	return v, nil
}
