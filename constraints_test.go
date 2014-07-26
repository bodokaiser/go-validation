package validation

import (
	"testing"

	"gopkg.in/check.v1"

	"gopkg.in/validator.v1"
)

func TestConstraints(t *testing.T) {
	check.Suite(&ConstraintsSuite{})
	check.TestingT(t)
}

type ConstraintsSuite struct{}

func (s *ConstraintsSuite) TestEmail(c *check.C) {
	c.Check(Email("", ""), check.IsNil)
	c.Check(Email("foo@bar.org", ""), check.IsNil)

	c.Check(Email(2, ""), check.Equals, validator.ErrUnsupported)
	c.Check(Email("foo", ""), check.Equals, validator.ErrInvalid)
	c.Check(Email("foo@", ""), check.Equals, validator.ErrInvalid)
	c.Check(Email("foo@bar", ""), check.Equals, validator.ErrInvalid)
	c.Check(Email("foo@bar.", ""), check.Equals, validator.ErrInvalid)
}

func (s *ConstraintsSuite) TestNested(c *check.C) {
	type nested struct {
		Name string `validate:"required"`
	}

	c.Check(Nested(nil, ""), check.IsNil)
	c.Check(Nested([]nested{}, ""), check.IsNil)

	c.Check(Nested(nested{}, ""), check.DeepEquals, Error{
		"Name": []error{
			validator.ErrZeroValue,
		},
	})

	c.Check(Nested([]nested{
		nested{},
		nested{},
	}, ""), check.DeepEquals, Errors{
		Error{
			"Name": []error{
				validator.ErrZeroValue,
			},
		},
		Error{
			"Name": []error{
				validator.ErrZeroValue,
			},
		},
	})
}

func (s *ConstraintsSuite) TestMinimum(c *check.C) {
	c.Check(Minimum(0, "1"), check.IsNil)
	c.Check(Minimum(5, "4"), check.IsNil)
	c.Check(Minimum("abc", "2"), check.IsNil)
	c.Check(Minimum([]int{1, 2}, "1"), check.IsNil)

	c.Check(Minimum(5, "a"), check.Equals, validator.ErrBadParameter)
	c.Check(Minimum(5, "6"), check.Equals, validator.ErrMin)
	c.Check(Minimum("abc", "4"), check.Equals, validator.ErrMin)
	c.Check(Minimum([]int{1, 2}, "3"), check.Equals, validator.ErrMin)
}

func (s *ConstraintsSuite) TestRequired(c *check.C) {
	c.Check(Required(1, ""), check.IsNil)
	c.Check(Required("a", ""), check.IsNil)
	c.Check(Required([]int{1}, ""), check.IsNil)

	c.Check(Required(nil, ""), check.Equals, validator.ErrZeroValue)
	c.Check(Required(0, ""), check.Equals, validator.ErrZeroValue)
	c.Check(Required("", ""), check.Equals, validator.ErrZeroValue)
	c.Check(Required([]int{}, ""), check.Equals, validator.ErrZeroValue)
}
