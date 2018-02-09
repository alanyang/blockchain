package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha1"
	"math/big"
	"testing"

	base58 "github.com/jbenet/go-base58"
)

func TestBase58(t *testing.T) {
	d := base58.Encode([]byte("alan"))
	if d != "3VSCm3" {
		t.Fail()
	}
}

func TestAddressGenerate(t *testing.T) {
	w := NewWallet()
	t.Log(w.Address().String())
}

func TestCrypto(t *testing.T) {
	curve := elliptic.P256()
	key, _ := ecdsa.GenerateKey(curve, rand.Reader)

	hash := sha1.Sum([]byte("Alan"))

	r, s, _ := ecdsa.Sign(rand.Reader, key, hash[:])

	pkh := append(key.PublicKey.X.Bytes(), key.PublicKey.Y.Bytes()...)

	var x, y big.Int

	x.SetBytes(pkh[:(len(pkh) / 2)])
	y.SetBytes(pkh[(len(pkh) / 2):])

	pubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}

	if !ecdsa.Verify(&pubKey, hash[:], r, s) {
		t.Fatal("Verify failure")
	} else {
		t.Log("Pass")
	}
}
