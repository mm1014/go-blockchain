package network

type NetAddr string

type Transport interface {
	Consume() <-chan RPC     // 获取接收 RPC 的通道
	Connect(Transport) error // 连接到另一个 Transport
	SendMessage(NetAddr, []byte) error
	Broadcast([]byte) error
	Addr() NetAddr
}
