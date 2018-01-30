package core

import (
	"github.com/boltdb/bolt"
	"strings"
)

const (
	DB_PATH = "./db/___cb.db"
)

var (
	CHAIN_BLOCK = []byte("amy-chain")
	LAST_TIP    = []byte{0x03, 0x00}
)

type (
	Chain struct {
		Tip []byte
		DB  *bolt.DB
	}

	ChainIterator struct {
		c  *Chain
		it []byte
	}
)

func NewBlockChain() (*Chain, error) {
	last := []byte{}
	db, err := bolt.Open(DB_PATH, 0600, nil)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(CHAIN_BLOCK)
		if b == nil {
			g := newGenesisBlock()
			b, _ := tx.CreateBucket(CHAIN_BLOCK)
			b.Put(g.Hash, g.Serialize())
			b.Put(LAST_TIP, g.Hash)
			last = g.Hash
		} else {
			last = b.Get(LAST_TIP)
		}
		return nil
	})
	return &Chain{last, db}, nil

}

func (c *Chain) AddBlock(b *Block) *Block {
	last := []byte{}

	c.DB.View(func(tx *bolt.Tx) error {
		bt := tx.Bucket(CHAIN_BLOCK)
		last = bt.Get(LAST_TIP)
		return nil
	})
	b.PreHash = last
	pow := NewProofOfWork(b)
	nonce, hash := pow.run()

	b.Hash = hash[:]
	b.Nonce = nonce
	c.Tip = b.Hash

	c.DB.Update(func(tx *bolt.Tx) error {
		bt := tx.Bucket(CHAIN_BLOCK)
		bt.Put(b.Hash, b.Serialize())
		bt.Put(LAST_TIP, c.Tip)
		return nil
	})

	return b
}

func (c *Chain) Iterator() *ChainIterator {
	return NewChainIterator(c)
}

func (c *Chain) String() string {
	s := []string{}
	it := c.Iterator()
	for {
		b := it.Next()
		if b == nil {
			break
		}
		s = append(s, b.String(), "👆")
	}
	return strings.Join(s, "\n")
}

func NewChainIterator(c *Chain) *ChainIterator {
	return &ChainIterator{c, c.Tip}
}

func (ci *ChainIterator) Next() (b *Block) {
	ci.c.DB.View(func(tx *bolt.Tx) error {
		bt := tx.Bucket(CHAIN_BLOCK)
		bl := bt.Get(ci.it)
		if bl != nil {
			b = Unserialize(bl)
		}
		return nil
	})
	if b != nil {
		ci.it = b.PreHash
	}
	return
}
