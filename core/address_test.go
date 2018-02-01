package core

import (
	base58 "github.com/jbenet/go-base58"
	"testing"
)

func TestBase58(t *testing.T) {
	d := base58.Encode([]byte("alan"))
	if d != "3VSCm3" {
		t.Fail()
	}
}

func TestAddressGenerate(t *testing.T) {
	w := NewWallet()
	t.Log(w.Address().String())
}
