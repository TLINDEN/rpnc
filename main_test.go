package main

import (
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	testscript.Main(m, map[string]func(){
		"rpn": main,
	})
}

func TestRpn(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "t",
	})
}
