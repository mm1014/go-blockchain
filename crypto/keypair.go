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

type PublicKey []byte
type Signature struct {
	S, R *big.Int
}

// 不改变k的内部状态 所以不用加*
func (k PrivateKey) Sign(data []byte) (*Signature, error) {
	//rand.Reader用于生成签名所需的随机数
	r, s, err := ecdsa.Sign(rand.Reader, k.key, data)
	if err != nil {
		return nil, err
	}
	return &Signature{R: r, S: s}, nil
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
func (k PrivateKey) PublicKey() PublicKey {
	//给定X、Y得到PublicKey
	return elliptic.MarshalCompressed(k.key.PublicKey, k.key.PublicKey.X, k.key.PublicKey.Y)
}

// 公钥本身的后 20 位可能会出现冲突,而哈希计算是计算完整公钥256位，不会出现冲突
func (k PublicKey) Address() types.Address {
	// 使用 SHA-256 哈希算法对公钥的压缩字节切片进行哈希计算
	h := sha256.Sum256(k)

	return types.AddressFromBytes(h[len(h)-20:])
}

func (sig *Signature) Verify(pubKey PublicKey, data []byte) bool {
	//给定公钥得到 x 和 y
	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), pubKey)
	key := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	return ecdsa.Verify(key, data, sig.R, sig.S)
}
