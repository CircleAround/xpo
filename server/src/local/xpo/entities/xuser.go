package entities

import "local/validatekit"

// XUser struct
type XUser struct {
	ID          string `datastore:"-" goon:"id" json:"id"`
	Name        string `json:"name" validate:"required,min=3,max=15,username_format"`
	Email       string `json:"email" validate:"required"`
	Nickname    string `json:"nickname" validate:"required,min=3,max=24,usernickname_format"`
	ReportCount int64  `json:"reportCount"`
}

// XUserProfileParams is parameter's basic
type XUserProfileParams struct {
	Name     string `json:"name" validate:"required,min=3,max=15,username_format"`
	Nickname string `json:"nickname" validate:"required,min=3,max=24,usernickname_format"`
}

func NewXUserValidator() *validatekit.Validate {
	v := validatekit.NewValidate()
	v.RegisterRegexValidation("username_format", `^[a-z][0-9a-z_]+$`)
	v.RegisterRegexValidation("usernickname_format", `^[^<>/:"'\s]+$`)
	return v
}
