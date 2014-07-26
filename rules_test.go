package validation

import (
	"errors"
	"testing"

	"gopkg.in/check.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/validator.v1"
)

func TestRules(t *testing.T) {
	check.Suite(&RulesSuite{})
	check.TestingT(t)
}

type RulesSuite struct{}

func (s *RulesSuite) TestId(c *check.C) {
	c.Check(Id("", ""), check.IsNil)
	c.Check(Id(nil, ""), check.IsNil)
	c.Check(Id(bson.NewObjectId(), ""), check.IsNil)

	c.Check(Id("abcd", ""), check.Equals, validator.ErrUnsupported)
	c.Check(Id(bson.ObjectId("1234"), ""), check.Equals, validator.ErrInvalid)
}

func (s *RulesSuite) TestRef(c *check.C) {
	c.Check(Ref(nil, ""), check.IsNil)
	c.Check(Ref(mgo.DBRef{
		Id:         bson.NewObjectId(),
		Collection: "collection1",
	}, ""), check.IsNil)

	c.Check(Ref(123, ""), check.Equals, validator.ErrUnsupported)
	c.Check(Ref(mgo.DBRef{
		Id: bson.NewObjectId(),
	}, ""), check.Equals, validator.ErrInvalid)
}

func (s *RulesSuite) TestEmail(c *check.C) {
	c.Check(Email("", ""), check.IsNil)
	c.Check(Email("foo@bar.org", ""), check.IsNil)

	c.Check(Email(2, ""), check.Equals, validator.ErrUnsupported)
	c.Check(Email("foo", ""), check.Equals, validator.ErrInvalid)
	c.Check(Email("foo@", ""), check.Equals, validator.ErrInvalid)
	c.Check(Email("foo@bar", ""), check.Equals, validator.ErrInvalid)
	c.Check(Email("foo@bar.", ""), check.Equals, validator.ErrInvalid)
}

func (s *RulesSuite) TestNested(c *check.C) {
	c.Check(Nested(nil, ""), check.IsNil)
	c.Check(Nested([]nested{}, ""), check.IsNil)

	c.Check(Nested(nested{}, ""), check.ErrorMatches, "error")

	c.Check(Nested([]nested{
		nested{},
		nested{},
	}, ""), check.DeepEquals, Errors{
		errors.New("error"),
		errors.New("error"),
	})
}

func (s *RulesSuite) TestMinimum(c *check.C) {
	c.Check(Minimum(0, "1"), check.IsNil)
	c.Check(Minimum(5, "4"), check.IsNil)
	c.Check(Minimum("abc", "2"), check.IsNil)
	c.Check(Minimum([]int{1, 2}, "1"), check.IsNil)

	c.Check(Minimum(5, "a"), check.Equals, validator.ErrBadParameter)
	c.Check(Minimum(5, "6"), check.Equals, validator.ErrMin)
	c.Check(Minimum("abc", "4"), check.Equals, validator.ErrMin)
	c.Check(Minimum([]int{1, 2}, "3"), check.Equals, validator.ErrMin)
}

func (s *RulesSuite) TestRequired(c *check.C) {
	c.Check(Required(1, ""), check.IsNil)
	c.Check(Required("a", ""), check.IsNil)
	c.Check(Required([]int{1}, ""), check.IsNil)

	c.Check(Required(nil, ""), check.Equals, validator.ErrZeroValue)
	c.Check(Required(0, ""), check.Equals, validator.ErrZeroValue)
	c.Check(Required("", ""), check.Equals, validator.ErrZeroValue)
	c.Check(Required([]int{}, ""), check.Equals, validator.ErrZeroValue)
}

type nested struct {
	Name string
}

func (n nested) Validate() error {
	return errors.New("error")
}
