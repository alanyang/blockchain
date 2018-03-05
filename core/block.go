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
	Block struct {
		//创建时间戳
		Timestamp int64
		//前一个块的Hash
		PrevHash []byte
		//交易数据
		Transactions []*Transaction
		//自己的Hash
		Hash []byte
		//随机数
		Nonce uint64

		Merkle []byte
	}
)

func (b *Block) hash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	header := bytes.Join([][]byte{timestamp, b.PrevHash, b.hashTransations()}, []byte{})

	h := sha256.Sum256(header)
	b.Hash = h[:]
}

func NewBlock(prevHash []byte, txs []*Transaction) *Block {
	b := &Block{
		PrevHash:     prevHash,
		Timestamp:    time.Now().Unix(),
		Transactions: txs,
	}
	pow := NewProofOfWork(b)
	nonce, hash := pow.run()
	b.Hash = hash
	b.Nonce = nonce
	return b
}

func newGenesisBlock() *Block {
	return NewBlock([]byte{}, []*Transaction{NewCoinbaseTransaction("1Aw2Nd6igX3qK9s2SGx7gAZLsmsHLwf79W", "")})
}

func (b *Block) HashString() string {
	return fmt.Sprintf("%X", b.Hash)
}

func (b *Block) String() string {
	return fmt.Sprintf("Prev: %X\nHash: %X", b.PrevHash, b.Hash)
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

func (b *Block) hashTransations() []byte {
	ids := [][]byte{}
	for _, tx := range b.Transactions {
		ids = append(ids, tx.ID)
	}

	tree := NewMerkleTree(ids)
	b.Merkle = tree.Root.Data

	return b.Merkle
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
