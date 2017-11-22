package gohttp

import (
	"errors"
	"github.com/sirupsen/logrus"
	"os"
	"regexp"
	"unicode/utf8"
)

// Logger info
var log = logrus.New()

func init() {
	log.Out = os.Stdout
	log.SetLevel(logrus.DebugLevel)
}

// TokenType is the general type of token IDs
type TokenType int

// Types of tokens
const (
	DIGIT TokenType = iota
	ALPHA
	CRLF
)

// SimpleToken represents a token with only an ID.
// For terminal symbols such as ( )
type SimpleToken struct {
	id TokenType
}

// DataToken represents a token with its watching string
// for undetermined terminals as ALHA
type DataToken struct {
	id     TokenType
	value  string
	Offset int
	Len    int
}

// Token interface only returns its type
type Token interface {
	getID() TokenType
}

// TokenWithData is a token that carries on its token hit
type TokenWithData interface {
	getValue() string
}

// GetID returns the ID of the SimpleToken
func (t *SimpleToken) GetID() TokenType {
	return t.id
}

// GetID returns the ID of the DataToken
func (t *DataToken) GetID() TokenType {
	return t.id
}

// GetValue returns the matching string of the DataToken
func (t *DataToken) GetValue() string {
	return t.value
}

// NewToken creates a new simple token
func NewToken(id TokenType) *SimpleToken {
	t := new(SimpleToken)
	t.id = id
	return t
}

// NewDataToken represents a token with its matching string
func NewDataToken(id TokenType, data string) *DataToken {
	t := new(DataToken)
	t.id = id
	t.value = data
	return t
}

// SimpleTokenMatcher is the basic Interface of a token matcher.
type SimpleTokenMatcher interface {
	// Match must return true if the following token matches its rule
	Match() bool
	// Gets the token that matches the rule, produces an error if the following
	// token is not of this type
	GetToken() (SimpleToken, error)
}

// DataTokenMatcher is the basic Interface of a data token matcher.
type DataTokenMatcher interface {
	// Match must return true if the following token matches its rule
	Match() bool
	// Gets the token that matches the rule, produces an error if the following
	// token is not of this type
	GetToken() (DataToken, error)
}

// ------------------- Protocol definition HTTP 1.0-----------------------------

type tokenRegex int

var terminalSymbols = map[TokenType]string{
	CRLF:  "\r\n",
	ALPHA: "a",
}

func compileHTTPTokens() *regexp.Regexp {
	//build the full regex
	fullRegex := ""
	for _, exp := range terminalSymbols {
		fullRegex += "(" + exp + ")|"
	}
	log.Debugf("Full regex <%s>", fullRegex[:utf8.RuneCountInString(fullRegex)])
	return regexp.MustCompile(fullRegex[:utf8.RuneCountInString(fullRegex)])
}

var terminalTokenizer = compileHTTPTokens()

type HTTPTokenizer struct {
	Message string
}

//Iter returns an iterator over the tokens of an HTTPTokenizer
// The iterator will return the next valid token of the string, nil if there
// is no more tokens or error if the string is not empty but could not be parsed.
//The iterator is consumed after used.
func (h *HTTPTokenizer) Iter() func() (*DataToken, error) {
	// This cursor is captured to keep track of the current token
	cursor := 0
	return func() (*DataToken, error) {
		if cursor >= utf8.RuneCountInString(h.Message) {
			return nil, nil
		} else {
			log.Debugf("Current Message <%s>", h.Message[cursor:])
			idx := terminalTokenizer.FindStringIndex(h.Message[cursor:])
			if idx == nil {
				return (*DataToken)(nil), errors.New("Error parsing message")
			} else {
				ts := h.Message[idx[0]+cursor : idx[1]+cursor]
				log.Debugf("Checking for Token Type of %s", ts)
				for k, v := range terminalSymbols {
					x := regexp.MustCompile(v)
					m := x.MatchString(ts)
					if m {
						tok := NewDataToken(k, ts)
						tok.Offset = idx[0] + cursor
						tok.Len = idx[1] + cursor
						cursor += idx[1]
						return tok, nil
					}
				}
				return nil, errors.New("Unmatched token, this must not happe" +
					"unless there is an error on the parser")
			}
		}
	}
}

// var regexSymbols = map[tokenRegex]*regexp.Regexp{
//	crlf:  regexp.MustCompile(terminalSymbols[crlf]),
//	alpha: regexp.MustCompile(terminalSymbols[alpha]),
// }

// // HTTPTokenizer digests a http request message
// type HTTPTokenizer struct {
//	Message string
//	cursor  int
// }

// NewHTTPTokenizer is the basic contructor
func NewHTTPTokenizer(request string) *HTTPTokenizer {
	tok := new(HTTPTokenizer)
	tok.Message = request
	return tok
}

// // NextToken gets the next valid token of the request
// func (tok *HTTPTokenizer) NextToken() *DataToken {
//	for _, value := range regexSymbols {
//		pos := value.FindStringIndex(tok.Message[tok.cursor:])
//		if pos != nil && pos[0] == 0 {
//			return NewDataToken(CRLF, tok.Message[tok.cursor:])
//		}
//	}
//	// for i := tok.cursor + 1; i < utf8.RuneCountInString(tok.Message[tok.cursor:]); i++ {
//	//	window := tok.Message[tok.cursor:i]
//	//	for _, value := range regexSymbols {
//	//		pos := value.FindStringIndex(window)
//	//		if pos != nil && pos[0] == 0 {
//	//			return NewDataToken(CRLF, window)
//	//		}
//	//	}
//	// }
//	return nil
// }
