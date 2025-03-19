package network

import (
	"bytes"
	"fmt"
	"sync"
)

type LocalTransport struct {
	addr      NetAddr
	consumeCh chan RPC
	lock      sync.RWMutex
	peers     map[NetAddr]*LocalTransport
}

func NewLocalTransport(addr NetAddr) Transport {
	return &LocalTransport{
		addr:      addr,
		consumeCh: make(chan RPC, 1024),
		peers:     make(map[NetAddr]*LocalTransport),
	}
}

// 返回一个只读通道
func (t *LocalTransport) Consume() <-chan RPC {
	return t.consumeCh
}

// 一般传入的不是Transport而是LocalTransport类型
func (t *LocalTransport) Connect(tr Transport) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	//断言tr为LocalTransport类型
	t.peers[tr.Addr()] = tr.(*LocalTransport)
	return nil
}

func (t *LocalTransport) SendMessage(to NetAddr, payload []byte) error {
	t.lock.RLock()
	defer t.lock.RUnlock()
	peer, ok := t.peers[to]
	if !ok {
		return fmt.Errorf("%s: could not send message to %s", t.addr, to)
	}
	peer.consumeCh <- RPC{ //将接收到的rpc放进对方的consume通道里面
		From: t.addr,
		//io.Reader支持从任意数据源读取数据（如内存、文件、网络流等）
		//bytes.NewReader(payload) 返回一个 *bytes.Reader，它实现了 io.Reader 接口
		Payload: bytes.NewReader(payload),
	}
	return nil
}

func (t *LocalTransport) Broadcast(payload []byte) error {
	//遍历当前节点已连接的所有节点
	for _, peer := range t.peers {
		if err := t.SendMessage(peer.Addr(), payload); err != nil {
			return err
		}
	}
	return nil
}

func (t *LocalTransport) Addr() NetAddr {
	return t.addr
}
