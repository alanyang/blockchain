package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"strconv"
	"time"
)

type (
	Trans struct {
	}
	Block struct {
		Timestamp     int64
		PreHash       []byte
		Data          []byte
		Hash          []byte
		Nonce         uint64
		Transtication []*Trans
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

func (b *Block) Serialize() []byte {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(b)
	if err != nil {
		return nil
	}
	return w.Bytes()
}

func Unserialize(d []byte) *Block {
	b := new(Block)
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(b)
	if err != nil {
		return nil
	}
	return b
}
