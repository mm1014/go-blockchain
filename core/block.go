package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"go-blockchain/crypto"
	"go-blockchain/types"
	"io"
)

type Header struct {
	Version       uint32
	DataHash      types.Hash
	PrevBlockHash types.Hash
	Timestamp     int64
	Height        uint32
	Nonce         uint64
}

type Block struct {
	*Header
	Transactions []Transaction
	Validator    crypto.PublicKey
	Signature    *crypto.Signature

	//Cached version of the header hash
	hash types.Hash
}

func NewBlock(h *Header, txs []Transaction) *Block {
	return &Block{
		Header:       h,
		Transactions: txs,
	}
}

// 将Header字段的数据编码字节切片的形式返回
func (b *Block) HeaderData() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(b.Header)
	return buf.Bytes()
}

func (b *Block) Sign(privKey crypto.PrivateKey) error {
	sig, err := privKey.Sign(b.HeaderData())
	if err != nil {
		return err
	}
	b.Validator = privKey.PublicKey()
	b.Signature = sig
	return nil
}

func (b *Block) Verify() error {
	if b.Signature == nil {
		return fmt.Errorf("block has no signature")
	}
	if !b.Signature.Verify(b.Validator, b.HeaderData()) {
		return fmt.Errorf("block has invalid signature")
	}
	return nil
}

// 专门用于解码 *Block 类型的对象
func (b *Block) Decode(r io.Reader, dec Decoder[*Block]) error {
	return dec.Decode(r, b)
}

func (b *Block) Encode(w io.Writer, dec Encoder[*Block]) error {
	return dec.Encode(w, b)
}

func (b *Block) Hash(hasher Hasher[*Block]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b)
	}
	return b.hash
}

// func (b *Block) Hash() types.Hash { //计算区块哈希
// 	buf := &bytes.Buffer{}     //创建一个内存缓冲区（bytes.Buffer），用于临时存储序列化后的二进制数据
// 	b.Header.EncodeBinary(buf) //将区块头（Header）序列化为二进制数据，并写入缓冲区 buf。
// 	if b.hash.IsZero() {
// 		b.hash = types.Hash(sha256.Sum256(buf.Bytes())) //对 data 计算 SHA-256 哈希，返回一个 [32]byte 类型的哈希值。
// 		// 将 [32]byte 转换为 types.Hash 类型（假设 types.Hash 是 [32]byte 的别名）。 底层类型相同的值进行显式类型转换
// 	}
// 	return b.hash
// }

// func (h *Header) EncodeBinary(w io.Writer) error {
// 	if err := binary.Write(w, binary.LittleEndian, &h.Version); err != nil {
// 		return err
// 	}
// 	if err := binary.Write(w, binary.LittleEndian, &h.PrevBlock); err != nil {
// 		return err
// 	}
// 	if err := binary.Write(w, binary.LittleEndian, &h.Timestamp); err != nil {
// 		return err
// 	}
// 	if err := binary.Write(w, binary.LittleEndian, &h.Height); err != nil {
// 		return err
// 	}
// 	return binary.Write(w, binary.LittleEndian, &h.Nonce)

// }

// func (h *Header) DecodeBinary(r io.Reader) error {
// 	if err := binary.Read(r, binary.LittleEndian, &h.Version); err != nil {
// 		return err
// 	}
// 	if err := binary.Read(r, binary.LittleEndian, &h.PrevBlock); err != nil {
// 		return err
// 	}
// 	if err := binary.Read(r, binary.LittleEndian, &h.Timestamp); err != nil {
// 		return err
// 	}
// 	if err := binary.Read(r, binary.LittleEndian, &h.Height); err != nil {
// 		return err
// 	}
// 	return binary.Read(r, binary.LittleEndian, &h.Nonce)

// }

// // Block将Header和Transactions的每一笔交易都序列化
// func (b *Block) EncodeBinary(w io.Writer) error {
// 	if err := b.Header.EncodeBinary(w); err != nil {
// 		return err
// 	}
// 	for _, tx := range b.Transactions {
// 		if err := tx.EncodeBinary(w); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func (b *Block) DecodeBinary(r io.Reader) error {
// 	if err := b.Header.DecodeBinary(r); err != nil {
// 		return err
// 	}
// 	for _, tx := range b.Transactions {
// 		if err := tx.DecodeBinary(r); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
