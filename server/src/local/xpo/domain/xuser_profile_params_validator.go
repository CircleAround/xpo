package domain

import (
	"bufio"
	"local/apikit"
	"local/stdkit"
	"local/validatekit"
	"local/xpo/assets"
	"local/xpo/entities"
)

func ValidateXUserProfileParams(params entities.XUserProfileParams) (*validatekit.Validate, error) {
	v := newXUserValidator()
	err := v.Struct(params)
	if err != nil {
		return nil, err
	}

	var prop string
	ng, err := checkBlockedWord(func(word string) bool {
		if params.Name == word {
			prop = "name"
			return true
		}
		if params.Nickname == word {
			prop = "nickname"
			return true
		}
		return false
	})

	if err != nil {
		return nil, err
	}

	if ng {
		return nil, apikit.NewInvalidParameterError(prop)
	}

	return v, nil
}

func checkBlockedWord(callback func(line string) bool) (bool, error) {
	f, err := assets.Assets.Open("/assets/reserved_username_list")
	if err != nil {
		return false, err
	}

	defer f.Close()

	reader := bufio.NewReaderSize(f, 128)
	hit, err := stdkit.FindLine(reader, callback)
	return hit, err
}

func newXUserValidator() *validatekit.Validate {
	v := validatekit.NewValidate()
	v.RegisterRegexValidation("username_format", `^[a-z][0-9a-z_]+$`)
	v.RegisterRegexValidation("usernickname_format", `^[^<>/:"'\s]+$`)
	return v
}
