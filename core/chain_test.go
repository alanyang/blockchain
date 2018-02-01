package core

import (
	// "crypto/sha256"
	// "math/big"
	"testing"
)

func TestAddBlock(t *testing.T) {
	c, err := NewBlockChain()
	if err != nil {
		t.Fatal(err)
	}
	// c.AddBlock(NewBlock("A block"))
	// c.AddBlock(NewBlock("B block"))
	// c.AddBlock(NewBlock("C block"))
	// c.AddBlock(NewBlock("D block"))
	t.Log(c.String())
}
