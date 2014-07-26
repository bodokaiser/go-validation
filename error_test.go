package validation

import (
	"fmt"
	"testing"

	"gopkg.in/check.v1"
	"gopkg.in/validator.v1"
)

func TestError(t *testing.T) {
	check.Suite(&ErrorSuite{})
	check.TestingT(t)
}

type ErrorSuite struct{}

func (s *ErrorSuite) TestError(c *check.C) {
	c.Check(Error{}.Error(), check.Equals, "")

	c.Check(Error{
		"foo": Errors{
			validator.ErrMin,
		},
	}.Error(), check.Matches, fmt.Sprintf("foo has %s", validator.ErrMin.Error()))
}

func (s *ErrorSuite) TestErrors(c *check.C) {
	c.Check(Errors{}.Error(), check.Equals, "")

	c.Check(Errors{
		validator.ErrMin,
	}.Error(), check.Matches, validator.ErrMin.Error())
}

func (s *ErrorSuite) TestErrorsAdd(c *check.C) {
	err := Errors{}

	c.Check(err.Add(validator.ErrMin), check.DeepEquals, Errors{
		validator.ErrMin,
	})
}
