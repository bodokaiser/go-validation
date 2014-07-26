package validation

import (
	"net/mail"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/validator.v1"
)

// Returns error if given value is not a bson object id.
//
// NOTE: It may be wise to remove the required check as the most common case for
// invalid object ids are uninitialized one which currently just are skipped.
func Id(v interface{}, _ string) error {
	if err := Required(v, ""); err != nil {
		return nil
	}

	switch t := v.(type) {
	case bson.ObjectId:
		if !t.Valid() {
			return validator.ErrInvalid
		}
	default:
		return validator.ErrUnsupported
	}
	return nil
}

// Returns error if the given value is not a valid mongo DBRef.
//
// NOTE: This only validates the basic requirements for a DBRef there is no
// check if the defined reference actually exists. This would need to be done
// seperated - maybe in the models Validate() method.
// TODO: Remove required check. See above for more.
func Ref(v interface{}, _ string) error {
	if err := Required(v, ""); err != nil {
		return nil
	}

	switch t := v.(type) {
	case mgo.DBRef:
		if id, ok := t.Id.(bson.ObjectId); ok && !id.Valid() {
			return validator.ErrInvalid
		}
		if len(t.Collection) == 0 {
			return validator.ErrInvalid
		}
	default:
		return validator.ErrUnsupported
	}
	return nil
}

// Returns error if given value is not an email.
func Email(v interface{}, _ string) error {
	if err := Required(v, ""); err != nil {
		return nil
	}

	switch r := reflect.ValueOf(v); r.Kind() {
	case reflect.String:
		s := r.String()

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
//
// NOTE: This actually just allows to call validation recursively and to
// integrate the errors into the existing Error map. It also supports recursive
// validation of Validatable slices.
func Nested(v interface{}, _ string) error {
	if v == nil {
		return nil
	}
	if r := reflect.ValueOf(v); r.Kind() == reflect.Slice {
		errs := Errors{}
		for i := 0; i < r.Len(); i++ {
			if n, ok := r.Index(i).Interface().(Validatable); ok {
				if err := n.Validate(); err != nil {
					errs = errs.Add(err)
				}
			}
		}
		if len(errs) > 0 {
			return errs
		}
		return nil
	}
	if n, ok := v.(Validatable); ok {
		return n.Validate()
	}
	return nil
}

// Returns error if given value exceeds minimum.
//
// NOTE: In difference to the built-in min validator this one will NOT validate
// zero or nil values. To get the old behaviour you will need to combine this
// with the required constraint.
func Minimum(v interface{}, s string) error {
	if err := Required(v, ""); err != nil {
		return nil
	}

	switch r := reflect.ValueOf(v); r.Kind() {
	case reflect.String, reflect.Slice, reflect.Map, reflect.Array:
		i, err := parseInt(s)

		if err != nil {
			return err
		}

		if int64(r.Len()) < i {
			return validator.ErrMin
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := parseInt(s)

		if err != nil {
			return err
		}

		if r.Int() < i {
			return validator.ErrMin
		}
	default:
		return validator.ErrUnsupported
	}
	return nil
}

// Returns error if given value is nil or zero value.
func Required(v interface{}, _ string) error {
	if v == nil {
		return validator.ErrZeroValue
	}

	switch r := reflect.ValueOf(v); r.Kind() {
	case reflect.Slice, reflect.Map, reflect.Array:
		if r.Len() == 0 {
			return validator.ErrZeroValue
		}
	default:
		if v == reflect.Zero(reflect.TypeOf(v)).Interface() {
			return validator.ErrZeroValue
		}
	}
	return nil
}

// Parses an int in string into int type.
func parseInt(s string) (int64, error) {
	i, err := strconv.ParseInt(s, 0, 64)

	if err != nil {
		return 0, validator.ErrBadParameter
	}

	return i, nil
}
