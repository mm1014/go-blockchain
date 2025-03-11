package network

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestConnect(t *testing.T){
// 	tra := NewLocalTransport("A")
// 	trb := NewLocalTransport("B")

// 	tra.Connect(trb)  //先链接才能断言  因为链接需要的是LocalTransport
// 	trb.Connect(tra)

// 	aTr := tra.(*LocalTransport)
// 	bTr := trb.(*LocalTransport)

// 	assert.Equal(t,aTr.peers[bTr.Addr()],bTr)
// 	assert.Equal(t,bTr.peers[aTr.Addr()],aTr)
// }

func TestSendMessage(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")

	tra.Connect(trb)
	trb.Connect(tra)

	aTr := tra.(*LocalTransport)
	bTr := trb.(*LocalTransport)

	msg := []byte("Hello world!")
	var receivedRPC RPC
	var wg sync.WaitGroup
	wg.Add(1)   // 增加等待的协程数量
	go func() { //开一个协程来监听数据
		defer wg.Done() // 协程完成任务后通知 WaitGroup
		receivedRPC = <-bTr.Consume()
	}()
	assert.Nil(t, aTr.SendMessage(bTr.Addr(), msg))
	// 等待接收协程完成
	wg.Wait()
	assert.Equal(t, receivedRPC.Payload, msg)
	assert.Equal(t, receivedRPC.From, aTr.Addr())
}
