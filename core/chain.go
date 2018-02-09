package core

import (
	"encoding/hex"
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
		Last []byte
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

	chain.Last = last

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

func (c *Chain) AddBlock(b *Block) (*Block, bool) {
	pow := NewProofOfWork(b)
	valid := pow.Verify()
	if valid {
		c.storeBlock(b)
	}
	return b, valid
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

func (c *Chain) FindUnspentTransactions(address string) (txs []*Transaction) {
	spentTx := make(map[string][]int)
	pubKey := PubkeyFromAddress(address)
	it := c.Iterator()
	for {
		block := it.Next()
		if block == nil {
			break
		}

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		OutPut:
			for outputIndex, output := range tx.Output {
				_, ok := spentTx[txID]
				if ok {
					for _, spentOut := range spentTx[txID] {
						if spentOut == outputIndex {
							continue OutPut
						}
					}
				}

				if output.IsLockedWithKey(pubKey) {
					txs = append(txs, tx)
				}
			}

			if !tx.IsCoinbase() {
				for inputIndex, input := range tx.Input {
					if input.UsesKey(pubKey) {
						inputID := hex.EncodeToString(input.ID)
						spentTx[inputID] = append(spentTx[inputID], inputIndex)
					}
				}
			}
		}
	}
	return
}

func (c *Chain) FindUTXO(address string) (utxo []*TransactionOutput) {
	txs := c.FindUnspentTransactions(address)
	for _, tx := range txs {
		for _, txo := range tx.Output {
			if txo.IsLockedWithKey(PubkeyFromAddress(address)) {
				utxo = append(utxo, txo)
			}
		}
	}
	return
}

func (c *Chain) GetBalance(address string) (b int) {
	utxo := c.FindUTXO(address)
	for _, txo := range utxo {
		b += txo.Value
	}
	return
}

//NewChainIterator ...
func NewChainIterator(c *Chain) *ChainIterator {
	return &ChainIterator{c, c.Last}
}

//Next ...
func (ci *ChainIterator) Next() (b *Block) {
	b = ci.c.getBlock(ci.it)
	if b != nil {
		ci.it = b.PrevHash
	}
	return
}
