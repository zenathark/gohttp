package lexer_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/zenathark/lexer/lexer"
	"testing"
)

func TestNewToken(t *testing.T) {
	k := lexer.NewToken(lexer.DIGIT)
	assert.Equal(t, k.GetID(), lexer.DIGIT, "NewToken should be DIGIT")
}
