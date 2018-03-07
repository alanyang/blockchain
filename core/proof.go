package core

import (
	"bytes"
	"encoding/binary"
	"log"
	"math/big"
	"time"
)

const (
	targetBits  = 12
	posMaxCount = 7200
)

const (
	TraditionConsensus = iota
	PrimeConsensus
)

//ProofOfWork ...
type ProofOfWork struct {
	b         *Block
	target    *big.Int
	hasher    Hasher
	consensus int
}

type ProofOfStake struct {
	pb      *Block
	account *Account
	hasher  Hasher
}

type PosResult struct {
	pb        *Block
	account   *Account
	timestamp int64
}

func NewProofOfWork(b *Block, t int) *ProofOfWork {
	target := big.NewInt(int64(1))
	target.Lsh(target, uint(256-targetBits))
	return &ProofOfWork{b: b, target: target, hasher: NewScryptHasher(), consensus: t}
}

func NewProofOfStake(pb *Block, account *Account) *ProofOfStake {
	return &ProofOfStake{pb: pb, account: account, hasher: NewSha256Hasher()}
}

func (p *ProofOfWork) prepare(nonce uint64) []byte {
	ts := bytes.NewBuffer([]byte{})
	binary.Write(ts, binary.BigEndian, p.b.Timestamp)
	binary.Write(ts, binary.BigEndian, int32(nonce))

	return bytes.Join([][]byte{
		p.b.hashTransations(), p.b.PrevHash, ts.Bytes(), []byte{targetBits},
	}, []byte{})
}

func (p *ProofOfWork) bitConsensus() (uint64, []byte) {
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

func (p *ProofOfWork) primeConsensus() (uint64, []byte) {
	var (
		nonce   uint64 = 0
		hashInt big.Int
	)

	for {
		hash := p.hasher.Sum(p.prepare(nonce))
		hashInt.SetBytes(hash[:])

		if hashInt.ProbablyPrime(0) {
			return nonce, hash[:]
		}
		nonce++
	}
}

//H(b) <= M
func (p *ProofOfWork) verifyBitConsensus() bool {
	var (
		hashInt      big.Int
		givenHashInt big.Int
	)
	hash := p.hasher.Sum(p.prepare(p.b.Nonce))
	hashInt.SetBytes(hash[:])
	givenHashInt.SetBytes(p.b.Hash)

	return givenHashInt.Cmp(&hashInt) == 0
}

//P(H(b)) ?
func (p *ProofOfWork) verifyPrimeConsensus() bool {
	var (
		hashInt      big.Int
		givenHashInt big.Int
	)
	hash := p.hasher.Sum(p.prepare(p.b.Nonce))
	hashInt.SetBytes(hash[:])
	givenHashInt.SetBytes(p.b.Hash)

	return givenHashInt.ProbablyPrime(0)
}

func (p *ProofOfWork) Consensus() (uint64, []byte) {
	if p.consensus == TraditionConsensus {
		return p.bitConsensus()
	} else if p.consensus == PrimeConsensus {
		return p.primeConsensus()
	} else {
		log.Fatal("Unsupport consensus")
	}
	return 0, nil
}

func (p *ProofOfWork) Verify() bool {
	if p.consensus == TraditionConsensus {
		return p.verifyBitConsensus()
	} else if p.consensus == PrimeConsensus {
		return p.verifyPrimeConsensus()
	} else {
		log.Fatal("Unsuuprt verify")
	}
	return false
}

//todo
//big.Int.ProbablyPrime is unsafe when x > (2 << 64)
func isPrime(a *big.Int) bool {
	return false
}

// H(H(pb) + t + A) <= balance(A) * M
func (pos *ProofOfStake) Consensus(ch chan *PosResult) {
	pbh := pos.hasher.Sum(pos.pb.Hash)
	var (
		target  big.Int
		hashInt big.Int
	)

	for i := 0; i < posMaxCount; i++ {
		ts := time.Now().Unix()
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.BigEndian, &ts)
		d := append(pbh, bytes.Join([][]byte{pos.account.Hash(), buf.Bytes()}, []byte{})...)

		hash := pos.hasher.Sum(d)
		hashInt.SetBytes(hash)

		target.Mul(big.NewInt(pos.account.Balance), big.NewInt(targetBits))

		if hashInt.Cmp(&target) >= 0 {
			ch <- &PosResult{
				pb:        pos.pb,
				account:   pos.account,
				timestamp: ts,
			}
			return
		}
		time.Sleep(time.Second)
	}

	ch <- nil
}
