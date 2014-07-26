package validation

import (
	"testing"

	"gopkg.in/check.v1"
	"gopkg.in/validator.v1"
)

func TestValidation(t *testing.T) {
	check.Suite(&ValidationSuite{})
	check.TestingT(t)
}

type ValidationSuite struct{}

func (s *ValidationSuite) TestValid(c *check.C) {
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

func (s *ValidationSuite) TestValidate(c *check.C) {
	type Model struct {
		Name    string `validate:"required,min=5"`
		Email   string `validate:"required,email,min=5"`
		Company string `validate:"min=5,max=20"`
	}

	c.Check(Validate(Model{
		Name:    "Bobby Joe",
		Email:   "bobby@joe.tx",
		Company: "Vacumizer Ind.",
	}), check.IsNil)
	c.Check(Validate(Model{
		Name:  "Bobby Joe",
		Email: "bobby@joe.tx",
	}), check.IsNil)

	c.Check(Validate(Model{}), check.DeepEquals, Error{
		"Name": Errors{
			validator.ErrZeroValue,
		},
		"Email": Errors{
			validator.ErrZeroValue,
		},
	})

	c.Check(Validate(Model{
		Name:  "Bob",
		Email: "b@tx",
	}), check.DeepEquals, Error{
		"Name": Errors{
			validator.ErrMin,
		},
		"Email": Errors{
			validator.ErrInvalid,
			validator.ErrMin,
		},
	})
}
