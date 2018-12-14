package domain

import (
	"bufio"
	"local/stdkit"
	"local/xpo/assets"

	validator "gopkg.in/go-playground/validator.v9"
)

func NewReservedIdentityNameValidation() (validator.Func, error) {
	return func(fl validator.FieldLevel) bool {
		f, err := assets.Assets.Open("/assets/reserved_username_list")
		if err != nil {
			panic(err)
		}

		defer f.Close()

		reader := bufio.NewReaderSize(f, 128)
		target := fl.Field().String()
		hit, err := stdkit.FindLine(reader, func(word string) bool {
			return target == word
		})

		if err != nil {
			panic(err)
		}
		return !hit
	}, nil
}
