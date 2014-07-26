package validation

import (
	"fmt"

	"gopkg.in/validator.v1"
)

type Error map[string][]error

func (err Error) Error() string {
	for k, errs := range err {
		return fmt.Sprintf("%s has %s", k, Errors(errs))
	}
	return ""
}

type Errors []error

func (errs Errors) Error() string {
	if len(errs) > 0 {
		return errs[0].Error()
	}
	return ""
}

var (
	DefaultValidator = validator.NewValidator()
)

func init() {
	DefaultValidator.SetValidationFunc("min", Minimum)
	DefaultValidator.SetValidationFunc("email", Email)
	DefaultValidator.SetValidationFunc("nested", Nested)
	DefaultValidator.SetValidationFunc("required", Required)
}

func Valid(v interface{}, s string) error {
	if ok, err := DefaultValidator.Valid(v, s); !ok {
		return Errors(err)
	}
	return nil
}

func Validate(v interface{}) error {
	if ok, err := DefaultValidator.Validate(v); !ok {
		return Error(err)
	}
	return nil
}
