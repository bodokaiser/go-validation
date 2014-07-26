// The validation package offers access to a validator and constraints which
// work with go-validator to extend its validation functionality.
package validation

import (
	"fmt"

	"gopkg.in/validator.v1"
)

// Error is returned as the validation result of a struct and contains all
// errors mapped by field name.
type Error map[string][]error

// Implements error interface by returning first error in map.
func (err Error) Error() string {
	for k, errs := range err {
		return fmt.Sprintf("%s has %s", k, Errors(errs))
	}
	return ""
}

// Error is internally used to implement the error interface on multiple errors.
type Errors []error

// Implements error interface by returning first error or empty string.
func (errs Errors) Error() string {
	if len(errs) > 0 {
		return errs[0].Error()
	}
	return ""
}

// Validator instance used to register ValidationFunc's.
var DefaultValidator = validator.NewValidator()

func init() {
	DefaultValidator.SetValidationFunc("min", Minimum)
	DefaultValidator.SetValidationFunc("email", Email)
	DefaultValidator.SetValidationFunc("nested", Nested)
	DefaultValidator.SetValidationFunc("required", Required)
}

// Validates a value with the given tag configuration.
// This is more used for testing internal validators you may prefer Validate().
func Valid(v interface{}, s string) error {
	if ok, err := DefaultValidator.Valid(v, s); !ok {
		return Errors(err)
	}
	return nil
}

// Validates a value with the configuration defined in tags and returns nil or
// an error map type of Error.
func Validate(v interface{}) error {
	if ok, err := DefaultValidator.Validate(v); !ok {
		return Error(err)
	}
	return nil
}
