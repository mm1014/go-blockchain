package network

import (
	"bytes"
	"go-blockchain/core"
	"go-blockchain/crypto"

	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

var defaultBlockTime = 5 * time.Second

type ServerOpts struct {
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor  RPCProcessor
	Transports    []Transport
	BlockTime     time.Duration //区块生成时间间隔
	PrivateKey    *crypto.PrivateKey
}

type Server struct {
	ServerOpts
	memPool     *TxPool
	isValidator bool
	rpcCh       chan RPC
	quitCh      chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}
	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}
	s := &Server{
		ServerOpts: opts,
		// blockTime:   opts.BlockTime,
		memPool:     NewTxPool(),
		isValidator: opts.PrivateKey != nil,
		rpcCh:       make(chan RPC, 1024),
		quitCh:      make(chan struct{}, 1),
	}
	// if opts.RPCHandler == nil {
	// 	opts.RPCHandler = NewDefaultRPCHandler(s)
	// }
	// s.ServerOpts = opts

	//If we dont got any processor from the server options, we going to use the server as default
	if s.RPCProcessor == nil {
		s.RPCProcessor = s
	}
	return s
}

// 暴露 rpcCh 的只读通道
func (s *Server) GetRPCCh() <-chan RPC {
	return s.rpcCh
}
func (s *Server) initTransports() { //为每个传输层启动一个 Goroutine，监听其 RPC 消息，并将消息转发到 s.rpcCh 通道
	for _, tr := range s.Transports {
		go func(tr Transport) { //每个传输层启动一个 Goroutine，监听其 RPC 消息
			for rpc := range tr.Consume() { //获取当前传输层的 RPC 消息通道，并不断从中读取消息
				s.rpcCh <- rpc //将接收到的 RPC 消息发送到服务器的 rpcCh 通道，供服务器统一处理。
			}
		}(tr)
	}
}

func (s *Server) Start() {
	s.initTransports()
	// 创建一个定时器，每隔 s.BlockTime 时间触发一次
	ticker := time.NewTicker(s.BlockTime)
free:
	for {
		select { //监听多个通道的操作
		case rpc := <-s.rpcCh:
			// if err := s.RPCHandler.HandleRPC(rpc); err != nil {
			// 	logrus.Error(err)
			// }
			msg, err := s.RPCDecodeFunc(rpc)
			if err != nil {
				logrus.Error(err)
			}
			if err := s.RPCProcessor.ProcessMessage(msg); err != nil {
				logrus.Error(err)
			}
		case <-s.quitCh:
			break free
		case <-ticker.C:
			if s.isValidator {
				s.createNewBlock()
			}
		}
	}
	fmt.Println("Server shutdown")
}

func (s *Server) ProcessMessage(msg *DecodedMessage) error {
	switch t := msg.Data.(type) {
	case *core.Transaction:
		return s.processTransaction(t)
	}
	return nil
}

func (s *Server) processTransaction(tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})
	if s.memPool.Has(hash) {
		logrus.WithFields(logrus.Fields{
			"hash": hash,
		}).Info("transaction already in mempool")

		return nil
	}

	if err := tx.Verify(); err != nil {
		return err
	}

	tx.SetFirstSeen(time.Now().UnixNano())

	logrus.WithFields(logrus.Fields{
		"hash":           hash,
		"mempool length": s.memPool.Len(),
	}).Info("adding new tx to the mempool")

	go s.broadcastTx(tx)
	return s.memPool.Add(tx)
}

// 广播最基本的函数
func (s *Server) broadcast(payload []byte) error {
	for _, tr := range s.Transports {
		if err := tr.Broadcast(payload); err != nil {
			return err
		}
	}
	return nil
}

// 广播Tx 调用broadcast
func (s *Server) broadcastTx(tx *core.Transaction) error {
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}
	msg := NewMessage(MessageTypeTx, buf.Bytes())
	return s.broadcast(msg.Bytes())
}

func (s *Server) createNewBlock() error {
	fmt.Println("do stuff every x seconds")
	return nil
}
