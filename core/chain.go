package core

import (
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
)

const (
	dbPath = "./db/blockchain"
)

var (
	bucketChainBlock = []byte("amy-chain")
	bucketLastTip    = []byte{0x03, 0x00}
)

type (
	//Chain ...
	Chain struct {
		last []byte
		db   *leveldb.DB
	}

	//ChainIterator ...
	ChainIterator struct {
		c  *Chain
		it []byte
	}
)

//NewBlockChain ...
func NewBlockChain() (*Chain, error) {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, err
	}
	chain := &Chain{db: db}
	last := chain.getLastTip()

	if last == nil {
		g := newGenesisBlock()
		chain.storeBlock(g)
		last = g.Hash
	}

	chain.last = last

	return chain, nil
}

func (c *Chain) getLastTip() (tip []byte) {
	tip, _ = c.db.Get(bucketChainBlock, nil)
	return
}

func (c *Chain) getLastBlock() (b *Block) {
	t := c.getLastTip()
	d, _ := c.db.Get(t, nil)
	return Unserialize(d)
}

func (c *Chain) getBlock(hash []byte) (b *Block) {
	d, _ := c.db.Get(hash, nil)
	b = Unserialize(d)
	return
}

func (c *Chain) storeBlock(b *Block) {
	c.db.Put(b.Hash, b.Serialize(), nil)
	c.db.Put(bucketChainBlock, b.Hash, nil)
}

func (c *Chain) AddBlock(b *Block) *Block {
	last := c.getLastTip()

	b.PreHash = last
	pow := NewProofOfWork(b)
	nonce, hash := pow.run()

	b.Hash = hash[:]
	b.Nonce = nonce
	c.last = b.Hash

	c.storeBlock(b)
	return b
}

func (c *Chain) Close() {
	c.db.Close()
}

//Iterator ...
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
		s = append(s, b.String(), "*")
	}
	return strings.Join(s, "\n")
}

//NewChainIterator ...
func NewChainIterator(c *Chain) *ChainIterator {
	return &ChainIterator{c, c.last}
}

//Next ...
func (ci *ChainIterator) Next() (b *Block) {
	b = ci.c.getBlock(ci.it)
	if b != nil {
		ci.it = b.PreHash
	}
	return
}
