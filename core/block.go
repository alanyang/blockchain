package core

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

type (
	Block struct {
		Timestamp int64
		PreHash   []byte
		Data      []byte
		Hash      []byte
		Nonce     uint64
	}
)

func (b *Block) hash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	header := bytes.Join([][]byte{timestamp, b.PreHash, b.Data}, []byte{})

	h := sha256.Sum256(header)
	b.Hash = h[:]
}

func NewBlock(data string) *Block {
	b := &Block{
		Timestamp: time.Now().Unix(),
		Data:      []byte(data),
	}
	b.hash()
	return b
}

func newGenesisBlock() *Block {
	return NewBlock("Genesis block")
}

func (b *Block) HashString() string {
	return fmt.Sprintf("%X", b.Hash)
}

func (b *Block) String() string {
	return fmt.Sprintf("Prev: %X\nData: %s\nHash: %X", b.PreHash, string(b.Data), b.Hash)
}
