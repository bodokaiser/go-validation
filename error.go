package validation

import "fmt"

// Represents a map of validation errors.
type Error map[string][]error

// Returns the first error of the map.
//
// NOTE: This will only expose one error at the time though multiple errors
// could occur. In this case you need to cast your error to Error and work with
// the map. Reason for this approach is that it is reporting multiple errors is
// not supported by our front-end nor our REST API.
func (err Error) Error() string {
	for k, errs := range err {
		return fmt.Sprintf("%s has %s", k, Errors(errs))
	}
	return ""
}

// Represents a collection of errors on the same type.
//
// NOTE: I am not 100% sure about the need of this type. May be removed in the
// future.
type Errors []error

// Returns an Error collection containing previous and provided errors.
//
// NOTE: This is a convenient method which is used when we iterate through
// slices so that we can merge multiple errors together.
// NOTE: We need to return a new type as we else would need to use Errors as
// pointer which MAY not be supported with the error type.
func (errs Errors) Add(err error) Errors {
	return append(errs, err)
}

// Returns the first error or an empty string if no errors exist.
func (errs Errors) Error() string {
	if len(errs) > 0 {
		return errs[0].Error()
	}
	return ""
}
