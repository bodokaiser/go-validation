// The validation package offers access to a validator and constraints which
// work with go-validator to extend its validation functionality.
package validation

import (
	"fmt"
	"net/mail"
	"reflect"
	"strconv"
	"strings"

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
//
// This is more used for testing internal validators you may prefer Validate().
func Valid(value interface{}, params string) error {
	if ok, err := DefaultValidator.Valid(value, params); !ok {
		return Errors(err)
	}

	return nil
}

// Validates a value with the configuration defined in tags and returns nil or
// an error map type of Error.
func Validate(value interface{}) error {
	if ok, err := DefaultValidator.Validate(value); !ok {
		return Error(err)
	}

	return nil
}

// Returns error if given value is not an email.
func Email(value interface{}, _ string) error {
	if err := Required(value, ""); err != nil {
		return nil
	}

	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.String:
		s := v.String()

		if _, err := mail.ParseAddress(s); err != nil {
			return validator.ErrInvalid
		}
		if i := strings.LastIndex(s, "."); i == -1 || i == len(s)-1 {
			return validator.ErrInvalid
		}
	default:
		return validator.ErrUnsupported
	}

	return nil
}

// Returns error if the given value does not validate.
func Nested(value interface{}, _ string) error {
	if value == nil {
		return nil
	}

	if v := reflect.ValueOf(value); v.Kind() == reflect.Slice {
		errs := make(Errors, 0)

		for i := 0; i < v.Len(); i++ {
			err := Validate(v.Index(i).Interface())

			if err != nil {
				errs = append(errs, err)
			}
		}

		if len(errs) > 0 {
			return errs
		}

		return nil
	}

	return Validate(value)
}

// Returns error if given value exceeds minimum.
func Minimum(value interface{}, param string) error {
	if err := Required(value, ""); err != nil {
		return nil
	}

	n, err := parseInt(param)
	if err != nil {
		return err
	}

	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.String, reflect.Slice, reflect.Map, reflect.Array:
		if int64(v.Len()) < n {
			return validator.ErrMin
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Int() < n {
			return validator.ErrMin
		}
	default:
		return validator.ErrUnsupported
	}

	return nil
}

// Returns error if given value is nil or zero value.
func Required(value interface{}, _ string) error {
	if value == nil {
		return validator.ErrZeroValue
	}

	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.Slice, reflect.Map, reflect.Array:
		if v.Len() == 0 {
			return validator.ErrZeroValue
		}
	default:
		if value == reflect.Zero(reflect.TypeOf(value)).Interface() {
			return validator.ErrZeroValue
		}
	}

	return nil
}

// Parses an int in string into int type.
func parseInt(s string) (int64, error) {
	n, err := strconv.ParseInt(s, 0, 64)

	if err != nil {
		return 0, validator.ErrBadParameter
	}

	return n, nil
}
