package core

import (
	"strings"
)

type (
	Chain struct {
		Blocks    []*Block
		blocksMap map[string]*Block
	}
)

func NewBlockChain() *Chain {
	g := newGenesisBlock()
	return &Chain{
		Blocks: []*Block{g},
		blocksMap: map[string]*Block{
			g.HashString(): g,
		},
	}
}

func (c *Chain) AddBlock(b *Block) {
	prev := c.Blocks[len(c.Blocks)-1]
	b.PreHash = prev.Hash
	pow := NewProofOfWork(b)
	nonce, hash := pow.run()

	b.Hash = hash[:]
	b.Nonce = nonce

	c.Blocks = append(c.Blocks, b)
	c.blocksMap[b.HashString()] = b
}

func (c *Chain) String() string {
	s := []string{}
	for _, b := range c.Blocks {
		s = append(s, b.String(), "==============================")
	}
	return strings.Join(s, "\n")
}
