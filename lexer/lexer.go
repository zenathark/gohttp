package lexer

import (
	// "errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	// "regexp"
	"unicode"
	"unicode/utf8"
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

const eof = -1

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
	beginState stateFn
}

// TokenIterator creates an iterator of all tokens found on input.
// It is consumed when used
type TokenIterator struct {
	start  int
	pos    int
	width  int
	input  string
	state  stateFn
	tokens chan token
}

type stateFn func(*TokenIterator) stateFn

// NewLexer returns a new instance of a lexer
func lex(name, input string, beginState stateFn) *lexer {
	l := &lexer{
		name:       name,
		input:      input,
		beginState: beginState,
	}
	return l
}

// Iter returns an iterator over all tokens of the lexer
func (l *lexer) Iter() (*TokenIterator, chan token) {
	iter := &TokenIterator{
		state:  l.beginState,
		start:  l.start,
		pos:    l.pos,
		width:  l.width,
		input:  l.input,
		tokens: make(chan token),
	}
	go iter.run()
	return iter, iter.tokens
}

func (ti *TokenIterator) run() {
	for state := ti.state; state != nil; {
		state = state(ti)
	}
	close(ti.tokens)
}

func (ti *TokenIterator) emit(t tokenType) {
	ti.tokens <- token{t, ti.input[ti.start:ti.pos]}
	ti.start = ti.pos
}

func (ti *TokenIterator) next() (r rune) {
	if ti.pos >= len(ti.input) {
		ti.width = 0
		return eof
	}
	r, ti.width = utf8.DecodeRuneInString(ti.input[ti.pos:])
	ti.pos += ti.width
	return r
}

// ------------------- Protocol definition HTTP 1.0-----------------------------

func octetLexer(ti *TokenIterator) stateFn {
	r, _ := utf8.DecodeRuneInString(ti.input[ti.pos:])
	if unicode.In(r, unicode.ASCII_Hex_Digit) {
		if ti.pos > ti.start {
			ti.emit(tokenOctet)
			return octetLexer
		}
	}
	return nil
}
