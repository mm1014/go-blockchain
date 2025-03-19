package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"go-blockchain/core"
	"io"

	"github.com/sirupsen/logrus"
)

type RPC struct {
	From    NetAddr
	Payload io.Reader
}
type MessageType byte

const (
	MessageTypeTx MessageType = 0x1
	MessageTypeBlock
)

type Message struct {
	Header MessageType
	Data   []byte
}

func NewMessage(t MessageType, data []byte) *Message {
	return &Message{
		Header: t,
		Data:   data,
	}
}

// 将msg编码后转成字节切片
func (msg *Message) Bytes() []byte {
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(msg)
	return buf.Bytes()
}

type DecodedMessage struct {
	From NetAddr
	Data any
}

type RPCDecodeFunc func(RPC) (*DecodedMessage, error)

func DefaultRPCDecodeFunc(rpc RPC) (*DecodedMessage, error) {
	msg := new(Message)
	//将rpc.payload解码到msg里面
	if err := gob.NewDecoder(rpc.Payload).Decode(&msg); err != nil {
		return nil, fmt.Errorf("failed to decode message from %s:%s", rpc.From, err)
	}
	logrus.WithFields(logrus.Fields{
		"from": rpc.From,
		"type": msg.Header,
	}).Debug("new incoming message")

	switch msg.Header {
	//把 msg.Data 解码成 Transaction 类型的对象。
	case MessageTypeTx:
		tx := new(core.Transaction)
		//调用 tx.Decode 方法进行解码，core.NewGobTxDecoder 用于创建一个解码器，
		// bytes.NewReader(msg.Data) 用于将 msg.Data 转换为可读的字节流。
		if err := tx.Decode(core.NewGobTxDecoder(bytes.NewReader(msg.Data))); err != nil {
			return nil, err
		}
		return &DecodedMessage{From: rpc.From, Data: tx}, nil
	default:
		return nil, fmt.Errorf("invalid message header %x", msg.Header)
	}

}

type RPCProcessor interface {
	// ProcessTransaction(NetAddr, *core.Transaction) error
	ProcessMessage(*DecodedMessage) error
}

//	type RPCHandler interface {
//		HandleRPC(rpc RPC) error
//	}
// type DefaultRPCHandler struct {
// 	p RPCProcessor
// }

// func NewDefaultRPCHandler(p RPCProcessor) *DefaultRPCHandler {
// 	return &DefaultRPCHandler{
// 		p: p,
// 	}
// }

// func (h *DefaultRPCHandler) HandleRPC(rpc RPC) error {
// 	msg := Message{}
// 	if err := gob.NewDecoder(rpc.Payload).Decode(&msg); err != nil {
// 		return fmt.Errorf("failed to decode message from %s:%s", rpc.From, err)
// 	}
// 	switch msg.Header {
// 	case MessageTypeTx:
// 		tx := new(core.Transaction)
// 		if err := tx.Decode(core.NewGobTxDecoder(bytes.NewReader(msg.Data))); err != nil {
// 			return err
// 		}
// 		return h.p.ProcessTransaction(rpc.From, tx)
// 	default:
// 		return fmt.Errorf("invalid message header %x", msg.Header)
// 	}

// 	return nil
// }
