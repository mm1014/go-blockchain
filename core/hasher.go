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

func (BlockHasher) Hash(b *Block) types.Hash {
	// buf := &bytes.Buffer{}
	// enc := gob.NewEncoder(buf)
	// if err := enc.Encode(b.Header); err != nil {
	// 	panic(err)
	// }
	h := sha256.Sum256(b.HeaderData())
	return types.Hash(h)
}
