package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	base58 "github.com/jbenet/go-base58"
	"golang.org/x/crypto/ripemd160"
	"math/big"
)

var (
	checkSumSize   = 4
	addressVersion = []byte{0x0}
)

type (
	Wallet struct {
		PrivateKey *ecdsa.PrivateKey
		PublicKey  []byte
	}

	Address []byte
)

func newKeyPair() (*ecdsa.PrivateKey, []byte) {
	el := elliptic.P256()
	pKey, _ := ecdsa.GenerateKey(el, rand.Reader)
	return pKey, append(pKey.PublicKey.X.Bytes(), pKey.PublicKey.Y.Bytes()...)
}

func NewWallet() *Wallet {
	privateKey, publicKey := newKeyPair()
	return &Wallet{privateKey, publicKey}
}

func FromPrivateKey(key []byte) *Wallet {
	b := new(big.Int).SetBytes(key)
	pKey := &ecdsa.PrivateKey{D: b}
	return &Wallet{pKey, append(pKey.PublicKey.X.Bytes(), pKey.PublicKey.Y.Bytes()...)}
}

func FromPrivateKeyString(s string) *Wallet {
	b, ok := new(big.Int).SetString(s, 16)
	if !ok {
		return nil
	}
	pKey := &ecdsa.PrivateKey{D: b}
	return &Wallet{pKey, append(pKey.PublicKey.X.Bytes(), pKey.PublicKey.Y.Bytes()...)}
}

func (w *Wallet) Address() Address {
	hash := w.hashPubkey()
	versionHash := append(addressVersion, hash...)
	sum := CheckSum(versionHash)
	return Address(append(versionHash, sum...)[:])
}

func (w *Wallet) hashPubkey() []byte {
	hash := sha256.Sum256(w.PublicKey)
	ripemdHash := ripemd160.New()
	ripemdHash.Write(hash[:])
	return ripemdHash.Sum(nil)
}

func CheckSum(b []byte) []byte {
	first := sha256.Sum256(b)
	second := sha256.Sum256(first[:])
	return second[:checkSumSize]
}

func (a Address) String() string {
	return base58.Encode([]byte(a))
}
