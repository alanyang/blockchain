package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"math/big"

	"encoding/gob"
	"encoding/hex"
	"fmt"
)

const (
	subsidy = 25
)

type (
	TransactionID = []byte

	//Transation struct
	Transaction struct {
		ID     TransactionID
		Input  []*TransactionInput
		Output []*TransactionOutput
	}

	//TransationInput struct
	TransactionInput struct {
		//prev transation id
		ID TransactionID
		//prev index of transations
		Index int
		//owner sign
		Sign      []byte
		PublicKey PublicKey
	}

	//TransationOutput struct
	TransactionOutput struct {
		//amount
		Value int
		//just a address or script
		PubKeyHash []byte
	}
)

func (tx *Transaction) IsCoinbase() bool {
	return tx.Input[0].Index == -1 && bytes.Compare(tx.Input[0].PublicKey, []byte{}) == 0 && len(tx.Input[0].ID) == 0
}

func (tx *Transaction) TrimmedCopy() *Transaction {
	txi := []*TransactionInput{}
	txo := []*TransactionOutput{}

	for _, tx := range tx.Input {
		txi = append(txi, &TransactionInput{tx.ID, tx.Index, nil, nil})
	}

	for _, tx := range tx.Output {
		txo = append(txo, &TransactionOutput{tx.Value, tx.PubKeyHash})
	}

	return &Transaction{tx.ID, txi, txo}
}

func (tx *Transaction) Sign(privateKey *ecdsa.PrivateKey, prevTxs map[string]*Transaction) error {
	if tx.IsCoinbase() {
		return errors.New("The is coinbase transaction.")
	}

	txCopy := tx.TrimmedCopy()
	for i, txi := range txCopy.Input {
		prevTx := prevTxs[hex.EncodeToString(txi.ID)]
		txCopy.Input[i].Sign = nil
		txCopy.Input[i].PublicKey = prevTx.Output[txi.Index].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Input[i].PublicKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, privateKey, txCopy.ID)
		if err != nil {
			return err
		}
		signature := append(r.Bytes(), s.Bytes()...)
		txCopy.Input[i].Sign = signature
	}
	return nil
}

func (tx *Transaction) Verify(prevTxs map[string]*Transaction) bool {
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for i, txi := range txCopy.Input {
		prevTx := prevTxs[hex.EncodeToString(txi.ID)]
		txCopy.Input[i].Sign = nil
		txCopy.Input[i].PublicKey = prevTx.Output[txi.Index].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Input[i].PublicKey = nil

		var r, s, x, y big.Int

		sigLen := len(txCopy.Input[i].Sign)
		r.SetBytes(txCopy.Input[i].Sign[:(sigLen / 2)])
		s.SetBytes(txCopy.Input[i].Sign[(sigLen / 2):])

		keyLen := len(txi.PublicKey)
		x.SetBytes(txi.PublicKey[:(keyLen / 2)])
		y.SetBytes(txi.PublicKey[(keyLen / 2):])

		raw := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if !ecdsa.Verify(&raw, txCopy.ID, &r, &s) {
			return false
		}
	}
	return true
}

func NewCoinbaseTransaction(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	input := &TransactionInput{
		[]byte{}, -1, []byte(data), []byte{},
	}
	output := &TransactionOutput{Value: subsidy}
	output.Lock(to)

	tx := &Transaction{Input: []*TransactionInput{input}, Output: []*TransactionOutput{output}}
	tx.ID = tx.Hash()
	return tx
}

func (tx *Transaction) Hash() []byte {
	txCopy := *tx
	txCopy.ID = []byte{}
	hash := sha256.Sum256(tx.Serialize())
	return hash[:]
}

func (tx *Transaction) Serialize() []byte {
	w := new(bytes.Buffer)
	en := gob.NewEncoder(w)
	err := en.Encode(*tx)
	if err != nil {
		return nil
	}
	return w.Bytes()
}

func (txo *TransactionOutput) Lock(address string) {
	txo.PubKeyHash = PubkeyFromAddress(address)
}

func (txo *TransactionOutput) IsLockedWithKey(key []byte) bool {
	return bytes.Compare(key, txo.PubKeyHash) == 0
}

func (txi *TransactionInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubkey(txi.PublicKey)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
