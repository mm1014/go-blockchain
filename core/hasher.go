package core

import (
	"crypto/sha256"
	"go-blockchain/types"
)

type Hasher[T any] interface {
	Hash(T) types.Hash
}

// 对Hash方法具体实现的结构体
type BlockHasher struct{}

func (BlockHasher) Hash(h *Header) types.Hash {
	h1 := sha256.Sum256(h.Bytes())
	return types.Hash(h1)
}

type TxHasher struct{}

func (TxHasher) Hash(tx *Transaction) types.Hash {
	return types.Hash(sha256.Sum256(tx.Data))
}
