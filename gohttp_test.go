package gohttp_test

import (
	"github.com/zenathark/gohttp"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestNewToken(c *C) {
	k := gohttp.NewToken(gohttp.DIGIT)
	c.Assert(k.GetID(), Equals, gohttp.DIGIT)
}

func (s *MySuite) TestNewDataToken(c *C) {
	k := gohttp.NewDataToken(gohttp.DIGIT, "3")
	c.Assert(k.GetID(), Equals, gohttp.DIGIT)
	c.Assert(k.GetValue(), Equals, "3")
}

func (s *MySuite) TestNextToken(c *C) {
	t := gohttp.NewHTTPTokenizer("a")
	token := t.NextToken()
	c.Assert(t.Message, Equals, "a")
	c.Assert(token, Not(Equals), (*gohttp.DataToken)(nil))
	c.Assert(token, Equals, gohttp.CRLF)
}
