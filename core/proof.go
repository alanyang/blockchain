package core

import (
	"bytes"
	"encoding/binary"
	"math/big"
)

const targetBits = 12

//ProofOfWork ...
type ProofOfWork struct {
	b      *Block
	target *big.Int
	hasher Hasher
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(int64(1))
	target.Lsh(target, uint(256-targetBits))
	return &ProofOfWork{b, target, NewScryptHasher()}
}

func (p *ProofOfWork) prepare(nonce uint64) []byte {
	ts := bytes.NewBuffer([]byte{})
	binary.Write(ts, binary.BigEndian, p.b.Timestamp)
	binary.Write(ts, binary.BigEndian, int32(nonce))

	return bytes.Join([][]byte{
		p.b.hashTransations(), p.b.PrevHash, ts.Bytes(), []byte{targetBits},
	}, []byte{})
}

func (p *ProofOfWork) run() (uint64, []byte) {
	var nonce uint64 = 0
	var hashInt big.Int

	for {
		hash := p.hasher.Sum(p.prepare(nonce))
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(p.target) == -1 {
			return nonce, hash[:]
		}
		nonce++
	}
}

func (p *ProofOfWork) Verify() bool {
	var (
		hashInt      big.Int
		givenHashInt big.Int
	)
	hash := p.hasher.Sum(p.prepare(p.b.Nonce))
	hashInt.SetBytes(hash[:])
	givenHashInt.SetBytes(p.b.Hash)

	return givenHashInt.Cmp(&hashInt) == 0
}
