package gohttp_test

import (
	"github.com/zenathark/gohttp"
	"testing"
)

func TestNewToken(t *testing.T) {
	k := gohttp.NewToken(gohttp.DIGIT)
	if k.GetId() != gohttp.DIGIT {
		t.Error("Expected token as DIGIT")
	}
}

func TestNewDataToken(t *testing.T) {
	k := gohttp.NewDataToken(gohttp.DIGIT, "3")
	if k.GetId() != gohttp.DIGIT || k.GetValue() != "3" {
		t.Error("Expected token as DIGIT")
	}
}
