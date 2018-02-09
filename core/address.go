package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"

	base58 "github.com/jbenet/go-base58"
	"golang.org/x/crypto/ripemd160"
)

var (
	checkSumSize   = 4
	addressVersion = []byte{0x0}
)

type (
	Wallet struct {
		PrivateKey *ecdsa.PrivateKey
		PublicKey  PublicKey
	}

	Address   []byte
	PublicKey = []byte
)

func newKeyPair() (*ecdsa.PrivateKey, PublicKey) {
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
	hash := HashPubkey(w.PublicKey)
	versionHash := append(addressVersion, hash...)
	sum := CheckSum(versionHash)
	return Address(append(versionHash, sum...)[:])
}

func HashPubkey(publicKey []byte) []byte {
	hash := sha256.Sum256(publicKey)
	ripemdHash := ripemd160.New()
	ripemdHash.Write(hash[:])
	return ripemdHash.Sum(nil)
}

func (w *Wallet) String() string {
	return fmt.Sprintf("PrivateKey: %X[Very important!]\nPublicKey:  %X\nAddress:    %s\n", w.PrivateKey.D.Bytes(), w.PublicKey, w.Address().String())
}

func CheckSum(b []byte) []byte {
	first := sha256.Sum256(b)
	second := sha256.Sum256(first[:])
	return second[:checkSumSize]
}

func (a Address) String() string {
	return base58.Encode([]byte(a))
}

func (a Address) Bytes() []byte {
	return []byte(a)
}

func PubkeyFromAddress(addr string) []byte {
	b := base58.Decode(addr)
	return b[1 : len(b)-4]
}
