package validation

import (
	"fmt"
	"testing"

	"gopkg.in/check.v1"
	"gopkg.in/validator.v1"
)

func TestValidation(t *testing.T) {
	check.Suite(&Suite{})
	check.TestingT(t)
}

type Suite struct{}

func (s *Suite) TestError(c *check.C) {
	c.Check(Error{}.Error(), check.Equals, "")
	c.Check(Error{"foo": Errors{validator.ErrMin}}.Error(),
		check.Matches, fmt.Sprintf("foo has %s", validator.ErrMin.Error()))
}

func (s *Suite) TestErrors(c *check.C) {
	c.Check(Errors{}.Error(), check.Equals, "")
	c.Check(Errors{validator.ErrMin}.Error(),
		check.Matches, validator.ErrMin.Error())
}

func (s *Suite) TestValid(c *check.C) {
	c.Check(Valid("", "email"), check.IsNil)
	c.Check(Valid("", "min=4"), check.IsNil)
	c.Check(Valid("foo", "required,min=2"), check.IsNil)
	c.Check(Valid("foo@bar.org", "required,email"), check.IsNil)

	c.Check(Valid("", "required,email"), check.DeepEquals, Errors{
		validator.ErrZeroValue,
	})
	c.Check(Valid("", "required,min=3"), check.DeepEquals, Errors{
		validator.ErrZeroValue,
	})
	c.Check(Valid("foo@", "required,email"), check.DeepEquals, Errors{
		validator.ErrInvalid,
	})
	c.Check(Valid("fobo", "required,min=5"), check.DeepEquals, Errors{
		validator.ErrMin,
	})
}

func (s *Suite) TestValidate(c *check.C) {
	type Model struct {
		Name    string `validate:"required,min=5"`
		Email   string `validate:"required,email,min=5"`
		Company string `validate:"min=5,max=20"`
	}

	c.Check(Validate(Model{"Bob Joe", "bob@joe.tx", "Vac Inc."}), check.IsNil)
	c.Check(Validate(Model{"Bob Joe", "bob@joe.tx", ""}), check.IsNil)

	c.Check(Validate(Model{}), check.DeepEquals, Error{
		"Name":  Errors{validator.ErrZeroValue},
		"Email": Errors{validator.ErrZeroValue},
	})
	c.Check(Validate(Model{"Bob", "b@tx", ""}), check.DeepEquals, Error{
		"Name":  Errors{validator.ErrMin},
		"Email": Errors{validator.ErrInvalid, validator.ErrMin},
	})
}

func (s *Suite) TestEmail(c *check.C) {
	c.Check(Email("", ""), check.IsNil)
	c.Check(Email("foo@bar.org", ""), check.IsNil)

	c.Check(Email(2, ""), check.Equals, validator.ErrUnsupported)
	c.Check(Email("foo", ""), check.Equals, validator.ErrInvalid)
	c.Check(Email("foo@", ""), check.Equals, validator.ErrInvalid)
	c.Check(Email("foo@bar", ""), check.Equals, validator.ErrInvalid)
	c.Check(Email("foo@bar.", ""), check.Equals, validator.ErrInvalid)
}

func (s *Suite) TestNested(c *check.C) {
	type nested struct {
		Name string `validate:"required"`
	}

	c.Check(Nested(nil, ""), check.IsNil)
	c.Check(Nested([]nested{}, ""), check.IsNil)

	c.Check(Nested(nested{}, ""), check.DeepEquals, Error{
		"Name": []error{validator.ErrZeroValue},
	})

	c.Check(Nested([]nested{
		nested{},
		nested{},
	}, ""), check.DeepEquals, Errors{
		Error{"Name": []error{validator.ErrZeroValue}},
		Error{"Name": []error{validator.ErrZeroValue}},
	})
}

func (s *Suite) TestMinimum(c *check.C) {
	c.Check(Minimum(0, "1"), check.IsNil)
	c.Check(Minimum(5, "4"), check.IsNil)
	c.Check(Minimum("abc", "2"), check.IsNil)
	c.Check(Minimum([]int{1, 2}, "1"), check.IsNil)

	c.Check(Minimum(5, "a"), check.Equals, validator.ErrBadParameter)
	c.Check(Minimum(5, "6"), check.Equals, validator.ErrMin)
	c.Check(Minimum("abc", "4"), check.Equals, validator.ErrMin)
	c.Check(Minimum([]int{1, 2}, "3"), check.Equals, validator.ErrMin)
}

func (s *Suite) TestRequired(c *check.C) {
	c.Check(Required(1, ""), check.IsNil)
	c.Check(Required("a", ""), check.IsNil)
	c.Check(Required([]int{1}, ""), check.IsNil)

	c.Check(Required(nil, ""), check.Equals, validator.ErrZeroValue)
	c.Check(Required(0, ""), check.Equals, validator.ErrZeroValue)
	c.Check(Required("", ""), check.Equals, validator.ErrZeroValue)
	c.Check(Required([]int{}, ""), check.Equals, validator.ErrZeroValue)
}
