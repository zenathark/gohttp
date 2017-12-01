package lexer

import (
	// "errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	// "regexp"
	"strings"
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
	tokenNone

	tokenOctet
	tokenEOF
)

// Token holds all information of a processed symbol
type token struct {
	Typ tokenType
	Val string
}

// String returns a string representation of a token
func (t token) String() string {
	switch t.Typ {
	case tokenEOF:
		return "EOF"
	case tokenError:
		return t.Val
	}
	if len(t.Val) > 10 {
		return fmt.Sprintf("%.10q...", t.Val)
	}
	return fmt.Sprintf("%q", t.Val)
}

type lexer struct {
	name   string
	input  string
	start  int
	pos    int
	width  int
	state  stateFn
	tokens chan token
}

type stateFn func(*lexer) stateFn

// NewLexer returns a new instance of a lexer
func NewLexer(name, input string, beginState stateFn) *lexer {
	l := &lexer{
		name:   name,
		input:  input,
		state:  beginState,
		tokens: make(chan token, 2),
	}
	go l.run(true)
	return l
}

// newLexer returns a new instance of a lexer
// This method is mainly for test purpouses
// if true, runTokenizer instructs the function to start the lexer's
// go routine.
// The closeOnFinish argument instructs the function to close the
// channel of the lexer on finish. This is the normal operation but can
// be left open for testing.
func newLexer(name, input string, beginState stateFn, runTokenizer, closeOnFinish bool) *lexer {
	l := &lexer{
		name:   name,
		input:  input,
		state:  beginState,
		tokens: make(chan token, 2),
	}
	if runTokenizer {
		go l.run(closeOnFinish)
	}
	return l
}

func (ti *lexer) run(closeOnFinish bool) {
	for state := ti.state; state != nil; {
		state = state(ti)
	}
	if closeOnFinish {
		close(ti.tokens)
	}
}

func (ti *lexer) emit(t tokenType) {
	ti.tokens <- token{t, ti.input[ti.start:ti.pos]}
	ti.start = ti.pos
}

func (ti *lexer) next() (r rune, eof bool) {
	if ti.pos >= len(ti.input) {
		ti.width = 0
		return 0, true
	}
	r, ti.width = utf8.DecodeRuneInString(ti.input[ti.pos:])
	ti.pos += ti.width
	return r, false
}

func (ti *lexer) peek() rune {
	r, _ := ti.next()
	ti.backward()
	return r
}

func (ti *lexer) backward() {
	ti.pos -= ti.width
}

func (ti *lexer) ignore() {
	ti.start = ti.pos
}

func (ti *lexer) accept(valid string) bool {
	r, _ := ti.next()
	if strings.IndexRune(valid, r) >= 0 {
		return true
	}
	ti.backward()
	return false
}

func (ti *lexer) acceptRun(valid string) {
	for ti.accept(valid) {
	}
}

// NextToken return the following token of the string
// func (l *lexer) NextToken() token {
//	for {
//		select {
//		case item := <-l.tokens:
//			return item
//		default:
//			l.state = l.state(l)
//		}
//	}
//	panic("Should not be reached")
// }

// ------------------- Protocol definition HTTP 1.0-----------------------------

func emptyState(ti *lexer) stateFn {
	if ti.pos > ti.start {
		ti.emit(tokenNone)
		return emptyState
	}
	_, eof := ti.next()
	if eof {
		ti.emit(tokenEOF)
		return nil
	}
	return nil
}

func coloopAState(ti *lexer) stateFn {
	ti.acceptRun("A")
	if ti.pos > ti.start {
		ti.emit(tokenNone)
		return coloopBState
	}
	_, eof := ti.next()
	if eof {
		ti.emit(tokenEOF)
		return nil
	}
	return nil
}

func coloopBState(ti *lexer) stateFn {
	ti.acceptRun("B")
	if ti.pos > ti.start {
		ti.emit(tokenNone)
		return coloopAState
	}
	_, eof := ti.next()
	if eof {
		ti.emit(tokenEOF)
		return nil
	}
	return nil
}

func octetLexer(ti *lexer) stateFn {
	r, _ := utf8.DecodeRuneInString(ti.input[ti.pos:])
	if unicode.In(r, unicode.ASCII_Hex_Digit) {
		if ti.pos > ti.start {
			ti.emit(tokenOctet)
			return octetLexer
		}
		_, eof := ti.next()
		if eof {
			ti.emit(tokenEOF)
			return nil
		}
	}
	return nil
}
