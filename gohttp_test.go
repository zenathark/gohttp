package gohttp_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/zenathark/gohttp"
	"testing"
)

func TestNewToken(t *testing.T) {
	k := gohttp.NewToken(gohttp.DIGIT)
	assert.Equal(t, k.GetID(), gohttp.DIGIT, "NewToken should be DIGIT")
}

func TestNewDataToken(t *testing.T) {
	k := gohttp.NewDataToken(gohttp.DIGIT, "3")
	assert.Equal(t, k.GetID(), gohttp.DIGIT)
	assert.Equal(t, k.GetValue(), "3", "Content should be 3")
}

func TestTokenIter(t *testing.T) {
	tk := gohttp.NewHTTPTokenizer("\r\na")
	assert.Equal(t, tk.Message, "\r\na", "Message should content the correct string")
	it := tk.Iter()
	fst, err := it()
	assert.Equal(t, err, nil)
	assert.Equal(t, fst.GetID(), gohttp.CRLF, "Should match CRLF.")
	scd, err := it()
	assert.Equal(t, err, nil)
	assert.Equal(t, scd.GetID(), gohttp.ALPHA, "Shold match a character.")
	lst, err := it()
	assert.Equal(t, err, nil)
	assert.Equal(t, lst, (*gohttp.DataToken)(nil), "Should return nil if the iterator is exhausted.")
}
