package core

import (
	"fmt"
	"go-blockchain/crypto"
	"go-blockchain/types"
)

// Transaction要验证Data签名
type Transaction struct {
	Data      []byte
	From      crypto.PublicKey
	Signature *crypto.Signature

	//cached version of the tx data hash
	hash types.Hash
	// firstSeen is the timestamp of when this tx is first seen locally
	firstSeen int64
}

func NewTransaction(data []byte) *Transaction {
	return &Transaction{
		Data: data,
	}
}

func (tx *Transaction) Hash(hasher Hasher[*Transaction]) types.Hash {
	if tx.hash.IsZero() {
		tx.hash = hasher.Hash(tx)
	}
	return hasher.Hash(tx)
}

// 用发送方的私钥对Data签名
func (tx *Transaction) Sign(privKey crypto.PrivateKey) error {
	sig, err := privKey.Sign(tx.Data)
	if err != nil {
		return err
	}
	tx.From = privKey.PublicKey()
	tx.Signature = sig
	return nil
}

func (tx *Transaction) Verify() error {
	if tx.Signature == nil {
		return fmt.Errorf("transaction has no signature")
	}
	if !tx.Signature.Verify(tx.From, tx.Data) {
		return fmt.Errorf("invalid transaction signature")

	}
	return nil
}

func (tx *Transaction) Decode(dec Decoder[*Transaction]) error {
	return dec.Decode(tx)
}
func (tx *Transaction) Encode(enc Encoder[*Transaction]) error {
	return enc.Encode(tx)
}

func (tx *Transaction) SetFirstSeen(t int64) {
	tx.firstSeen = t
}
func (tx *Transaction) FirstSeen() int64 { return tx.firstSeen }

// func (tx *Transaction) DecodeBinary(r io.Reader) error { return nil }

// func (tx *Transaction) EncodeBinary(w io.Writer) error { return nil }
