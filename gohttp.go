package gohttp

import (
// "fmt"
)

type tokenType int

const (
	DIGIT tokenType = iota
	ALPHA
)

type simpleToken struct {
	id tokenType
}

type dataToken struct {
	id    tokenType
	value string
}

type Token interface {
	getId() tokenType
}

type TokenWithData interface {
	getValue() string
}

func (t *simpleToken) GetId() tokenType {
	return t.id
}

func (t *dataToken) GetId() tokenType {
	return t.id
}

func (t *dataToken) GetValue() string {
	return t.value
}

func NewToken(id tokenType) *simpleToken {
	t := new(simpleToken)
	t.id = id
	return t
}

func NewDataToken(id tokenType, data string) *dataToken {
	t := new(dataToken)
	t.id = id
	t.value = data
	return t
}
