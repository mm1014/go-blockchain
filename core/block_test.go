package core

import (
	"go-blockchain/crypto"
	"go-blockchain/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func randomBlock(height uint32, prevBlockHash types.Hash) *Block {
	header := &Header{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		Height:        height,
		Timestamp:     time.Now().UnixNano(),
	}
	// tx := Transaction{
	// 	Data: []byte("too"),
	// }
	return NewBlock(header, []Transaction{})
}

func randomBlockWithSignature(t *testing.T, height uint32, prevBlockHash types.Hash) *Block {
	privKey := crypto.GeneratePrivateKey()
	b := randomBlock(height, prevBlockHash)
	tx := randomTxWithSignature(t)
	b.AddTransaction(tx)
	assert.Nil(t, b.Sign(privKey))
	return b
}
func TestSignBlock(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	b := randomBlock(0, types.Hash{})
	assert.Nil(t, b.Sign(privKey))
	assert.NotNil(t, b.Signature)
}

func TestVerifyBlock(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	b := randomBlock(0, types.Hash{})
	assert.Nil(t, b.Sign(privKey))
	assert.Nil(t, b.Verify())

	otherPrivKey := crypto.GeneratePrivateKey()
	b.Validator = otherPrivKey.PublicKey()
	assert.NotNil(t, b.Verify())

	b.Height = 100
	assert.NotNil(t, b.Verify())
}

// func TestHeader_Encode_Decode(t *testing.T) {
// 	h := &Header{
// 		Version:   1,
// 		PrevBlock: types.RandomHash(32),
// 		Timestamp: time.Now().UnixNano(),
// 		Height:    10,
// 		Nonce:     989394,
// 	}
// 	buf := &bytes.Buffer{}
// 	assert.Nil(t, h.EncodeBinary(buf)) //若不报错返回的是Nil

// 	hDecode := &Header{}
// 	assert.Nil(t, hDecode.DecodeBinary(buf)) //从buf拿数据赋值到hDecode中
// 	assert.Equal(t, h, hDecode)
// }

// func TestBlock_Encode_Decode(t *testing.T) {
// 	b := &Block{
// 		Header: Header{
// 			Version:   1,
// 			PrevBlock: types.RandomHash(32),
// 			Timestamp: time.Now().UnixNano(),
// 			Height:    10,
// 			Nonce:     989394,
// 		},
// 		Transactions: nil,
// 	}
// 	buf := &bytes.Buffer{}
// 	assert.Nil(t, b.EncodeBinary(buf))

// 	bDecode := &Block{}
// 	assert.Nil(t, bDecode.DecodeBinary(buf))
// 	assert.Equal(t, b, bDecode)
// 	fmt.Printf("%+v", bDecode)
// }

// func TestBlockHash(t *testing.T) {
// 	b := &Block{
// 		Header: Header{
// 			Version:   1,
// 			PrevBlock: types.RandomHash(32),
// 			Timestamp: time.Now().UnixNano(),
// 			Height:    10,
// 		},
// 		Transactions: nil,
// 	}
// 	h := b.Hash()
// 	fmt.Println(h)
// 	assert.False(t, h.IsZero())
// }
