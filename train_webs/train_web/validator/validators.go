package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	return isMobile(mobile)
}

func isMobile(mobile string) bool {
	matched, _ := regexp.MatchString("1[3-9]\\d{9}", mobile)
	return matched
}
