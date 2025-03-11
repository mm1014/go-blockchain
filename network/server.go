package network

import (
	"fmt"
	"time"
)

type ServerOpts struct {
	Transports []Transport
}

type Server struct {
	ServerOpts
	rpcCh  chan RPC
	quitCh chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	return &Server{
		ServerOpts: opts,
		rpcCh:      make(chan RPC, 1024),
		quitCh:     make(chan struct{}, 1),
	}
}

// 暴露 rpcCh 的只读通道
func (s *Server) GetRPCCh() <-chan RPC {
	return s.rpcCh
}

func (s *Server) Start() {
	s.initTransports()
	ticker := time.NewTicker(5 * time.Second)
free:
	for {
		select { //监听多个通道的操作
		case rpc := <-s.rpcCh:
			fmt.Printf("%+v\n", rpc)
			time.Sleep(2 * time.Second) // 增加延迟，模拟消息处理时间
		case <-s.quitCh:
			break free
		case <-ticker.C:
			fmt.Println("do stuff every x seconds")
		}
	}
	fmt.Println("Server shutdown")
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
