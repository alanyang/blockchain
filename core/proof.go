package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"math/big"
)

const targetBits = 12

type ProofOfWork struct {
	b      *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(int64(1))
	target.Lsh(target, uint(256-targetBits))
	return &ProofOfWork{b, target}
}

func (p *ProofOfWork) prepare(nonce uint64) []byte {
	ts := bytes.NewBuffer([]byte{})
	binary.Write(ts, binary.BigEndian, p.b.Timestamp)
	binary.Write(ts, binary.BigEndian, int32(nonce))

	return bytes.Join([][]byte{
		p.b.Data, p.b.PreHash, ts.Bytes(), []byte{targetBits},
	}, []byte{})
}

func (p *ProofOfWork) run() (uint64, []byte) {
	var nonce uint64 = 0
	var hashInt big.Int

	for {
		hash := sha256.Sum256(p.prepare(nonce))
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(p.target) == -1 {
			return nonce, hash[:]
		} else {
			// fmt.Printf("%X - %d\n", hash, nonce)
		}
		nonce++
	}
}
