package lexer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"unicode/utf8"
)

func TestNewlexer(t *testing.T) {
	lex := NewLexer("Test1", "AO53", octetLexer)
	assert.Equal(t, lex.name, "Test1", "Name should be initialized.")
	assert.Equal(t, lex.input, "AO53", "Input should be initialized.")
}

func TestEmit(t *testing.T) {
	lex := NewLexer("Test1", "AO53", octetLexer)
	lex.emit(tokenOctet)
	tk := <-lex.tokens
	assert.Equal(t, reflect.DeepEqual(tk, token{tokenOctet, ""}), true, "Token should be emmited")
}

func TestEmitFirstCharacter(t *testing.T) {
	lex := NewLexer("Test1", "AO53", octetLexer)
	lex.next()
	lex.emit(tokenOctet)
	tk := <-lex.tokens
	assert.Equal(t, reflect.DeepEqual(tk, token{tokenOctet, "A"}), true, "First character should be captured by emit")
}

func TestPeekFirstCharacter(t *testing.T) {
	lex := NewLexer("Test1", "AO53", octetLexer)
	peeked := lex.peek()
	assert.Equal(t, 'A', peeked, "First character should be return on peek")
	_, validWidth := utf8.DecodeRuneInString("A")
	assert.Equal(t, lex.width, validWidth, "Width should be the size of the rune")
	peeked = lex.peek()
	assert.Equal(t, 'A', peeked, "Peeking again should return the same character")
}

func TestNext(t *testing.T) {
	lex := NewLexer("Test1", "AO53", octetLexer)
	for i, c := range "AO53" {
		poped, _ := lex.next()
		assert.Equal(t, c, poped, fmt.Sprintf("%d-th character [%c] should be peeked", i, c))
	}
}

func TestNextPeek(t *testing.T) {
	lex := NewLexer("Test1", "AO53", octetLexer)
	for range "AO53" {
		peeked := lex.peek()
		poped, _ := lex.next()
		assert.Equal(t, peeked, poped, fmt.Sprintf("peeked character [%c] should be equal to poped character [%c]", peeked, poped))
	}
}

func TestPeekNextLoop(t *testing.T) {
	lex := NewLexer("Test1", "AO53", octetLexer)
	for i, c := range "AO53" {
		peeked := lex.peek()
		assert.Equal(t, c, peeked, fmt.Sprintf("%d-th character [%c] should be peeked", i, c))
		lex.next()
	}
}

func TestPeekEmitLoop(t *testing.T) {
	lex := newLexer("Test1", "AO53", octetLexer, true, false)
	for range "AO53" {
		c, _ := lex.next()
		tk := token{tokenOctet, fmt.Sprintf("%c", c)}
		lex.emit(tokenOctet)
		emitted := <-lex.tokens
		assert.Equal(t, reflect.DeepEqual(tk, emitted), true, "A token must be emitted")
	}
}

func TestEmitLongLoop(t *testing.T) {
	lex := newLexer("Test1", "AO53", octetLexer, true, false)
	var c []rune
	c = make([]rune, 2)
	c[0], _ = lex.next()
	c[1], _ = lex.next()
	tk := token{tokenOctet, string(c)}
	lex.emit(tokenOctet)
	emitted := <-lex.tokens
	assert.Equal(t, reflect.DeepEqual(tk, emitted), true, fmt.Sprintf("The token %s must be emitted", string(c)))
}

func TestAccept(t *testing.T) {
	lex := newLexer("Test1", "AO53", emptyState, false, false)
	for i := range "AO53" {
		c := lex.peek()
		accepted1 := lex.accept("53AO")
		assert.Equal(t, accepted1, true, fmt.Sprintf("%d-th Char [%c] must be accepted", i, c))
	}
	lex.emit(tokenNone)
	emmited := <-lex.tokens
	assert.Equal(t, emmited.Val, "AO53", "Emitted token must contain all string")
}

func TestRun(t *testing.T) {
	lex := newLexer("Test1", "ABABA", coloopAState, true, true)
	for tk := <-lex.tokens; tk.Typ != tokenEOF; tk = <-lex.tokens {
		fmt.Printf("%s is the symbol\n", tk.Val)
		assert.Equal(t, tk.Typ, tokenNone, "Emitted a None token")
		assert.Equal(t, tk.Val, "A", "Took just an A")
		tk := <-lex.tokens
		fmt.Printf("%s is the symbol\n", tk.Val)
		if tk.Typ == tokenEOF {
			break
		}
		assert.Equal(t, tk.Typ, tokenNone, "Emitted a None token")
		assert.Equal(t, tk.Val, "B", "Took just a B")
	}
}

func TestRunAndChannel(t *testing.T) {
	lex := newLexer("Test1", "ABABA", coloopAState, true, true)
	s := ""
	for tk := range lex.tokens {
		if tk.Typ != tokenEOF {
			assert.Equal(t, tk.Typ, tokenNone, "Emitted a None token")
			s += tk.Val
		}
	}
	assert.Equal(t, s, "ABABA", "Emitted all characters right")
}
