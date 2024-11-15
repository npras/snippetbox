package validator

import (
	"slices"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

func (v *Validator) IsValid() bool {
	return len(v.FieldErrors) == 0
}

func (v *Validator) AddFieldError(key, msg string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = msg
	}
}

func (v *Validator) CheckAndAddFieldError(ok bool, key, msg string) {
	if !ok {
		v.AddFieldError(key, msg)
	}
}

//

func NotBlank(val string) bool {
	return strings.TrimSpace(val) != ""
}

func LessThanMaxChars(val string, n int) bool {
	return utf8.RuneCountInString(val) <= n
}

func PermittedValue[T comparable](val T, permittedVals ...T) bool {
	return slices.Contains(permittedVals, val)
}
