package types

import (
	"crypto/rand"
	"encoding/hex"
)

type Hash [32]byte

func (h Hash) IsZero() bool {
	for i := 0; i < 32; i++ {
		if h[i] != 0 {
			return false
		}
	}
	return true
}

// 字节转字符串
func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

// 生成指定长度的随机字节切片。
func RandomHash(size int) Hash {
	token := make([]byte, size)
	rand.Read(token) //rand.Read() 会将生成的伪随机字节填充到这个token中。
	return Hash(token)
}
