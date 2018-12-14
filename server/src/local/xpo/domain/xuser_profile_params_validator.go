package domain

import (
	"local/validatekit"
	"local/xpo/entities"

	"github.com/pkg/errors"
)

func ValidateXUserProfileParams(params entities.XUserProfileParams) (*validatekit.Validate, error) {
	v, err := newXUserValidator()
	if err != nil {
		return nil, errors.Wrap(err, "newXUserValidator failed")
	}

	err = v.Struct(params)
	if err != nil {
		return nil, errors.Wrap(err, "Validation failed")
	}
	return v, nil
}

func newXUserValidator() (*validatekit.Validate, error) {
	v := validatekit.NewValidate()
	v.RegisterRegexValidation("identity_name_format", `^[a-z][0-9a-z_]+$`)
	v.RegisterRegexValidation("usernickname_format", `^[^<>/:"'\s]+$`)

	vali, err := NewReservedIdentityNameValidation()
	if err != nil {
		return nil, errors.Wrap(err, "NewReservedIdentityNameValidation failed")
	}
	v.RegisterValidation("reserved_identity_name", vali)

	return v, nil
}
