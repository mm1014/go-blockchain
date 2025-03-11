package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"go-blockchain/types"
	"math/big"
)

type PrivateKey struct {
	key *ecdsa.PrivateKey
}

type PublicKey struct {
	key *ecdsa.PublicKey
}

type Signature struct {
	s, r *big.Int
}

// 不改变k的内部状态 所以不用加*
func (k *PrivateKey) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k.key, data)
	if err != nil {
		return nil, err
	}
	return &Signature{s, r}, nil
}

func GeneratePrivateKey() PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	return PrivateKey{
		key: key,
	}
}

// 从私钥中提取对应的公钥。
func (k *PrivateKey) PublicKey() PublicKey {
	return PublicKey{
		key: &k.key.PublicKey,
	}
}

// 将公钥转换为压缩格式的字节切片
func (k *PublicKey) ToSlice() []byte {
	return elliptic.MarshalCompressed(k.key, k.key.X, k.key.Y)
}

// 公钥本身的后 20 位可能会出现冲突,而哈希计算是计算完整公钥256位，不会出现冲突
func (k *PublicKey) Address() types.Address {
	//使用 SHA-256 哈希算法对公钥的压缩字节切片进行哈希计算
	h := sha256.Sum256(k.ToSlice())
	return types.Address(h[len(h)-20:])
}

func (sig *Signature) Verify(pubKey PublicKey, data []byte) bool {
	return ecdsa.Verify(pubKey.key, data, sig.r, sig.s)
}
