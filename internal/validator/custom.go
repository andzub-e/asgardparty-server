package validator

import (
	"github.com/go-playground/validator/v10"
)

func stringOneOfTheList(list []string) func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		str, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}

		for _, item := range list {
			if str == item {
				return true
			}
		}

		return false
	}
}
