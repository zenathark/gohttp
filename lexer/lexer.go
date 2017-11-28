package lexer

import (
	// "errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	// "regexp"
	// "unicode/utf8"
)

// Logger info
var log = logrus.New()

func init() {
	log.Out = os.Stdout
	log.SetLevel(logrus.DebugLevel)
}

// TokenType is the general type of token IDs
type tokenType int

// Types of tokens
const (
	tokenError tokenType = iota

	tokenOctet
	tokenEOF
)

// Token holds all information of a processed symbol
type token struct {
	typ tokenType
	val string
}

// String returns a string representation of a token
func (t token) String() string {
	switch t.typ {
	case tokenEOF:
		return "EOF"
	case tokenError:
		return t.val
	}
	if len(t.val) > 10 {
		return fmt.Sprintf("%.10q...", t.val)
	}
	return fmt.Sprintf("%q", t.val)
}

type lexer struct {
	name       string
	input      string
	start      int
	pos        int
	width      int
	tokens     chan token
	beginState stateFn
}

// TokenIterator creates an iterator of all tokens found on input.
// It is consumed when used
type TokenIterator struct {
	state stateFn
}

type stateFn func(*lexer) stateFn

// NewLexer returns a new instance of a lexer
func lex(name, input string, beginState stateFn) (*lexer, chan token) {
	l := &lexer{
		name:       name,
		input:      input,
		tokens:     make(chan token),
		beginState: beginState,
	}
	return l, l.tokens
}

func (l *lexer) Iter() *TokenIterator {
	return &TokenIterator{
		state: l.beginState,
	}
}

// ------------------- Protocol definition HTTP 1.0-----------------------------
