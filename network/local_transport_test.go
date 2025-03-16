package network

import (
	"io/ioutil"
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

	// aTr := tra.(*LocalTransport)
	// bTr := trb.(*LocalTransport)

	msg := []byte("Hello world!")
	// var receivedRPC RPC
	// var wg sync.WaitGroup
	// wg.Add(1)   // 增加等待的协程数量
	// go func() { //开一个协程来监听数据
	// 	defer wg.Done() // 协程完成任务后通知 WaitGroup
	// 	receivedRPC = <-bTr.Consume()
	// }()
	assert.Nil(t, tra.SendMessage(trb.Addr(), msg))
	// // 等待接收协程完成
	// wg.Wait()
	rpc := <-trb.Consume()
	b, err := ioutil.ReadAll(rpc.Payload)
	// buf := make([]byte, len(msg))
	// n, err := rpc.Payload.Read(buf)
	assert.Nil(t, err)
	// assert.Equal(t, n, len(msg))
	assert.Equal(t, b, msg)
	assert.Equal(t, rpc.From, tra.Addr())
}

func TestBroadcast(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")
	trc := NewLocalTransport("C")

	tra.Connect(trb)
	tra.Connect(trc)

	msg := []byte("too")
	assert.Nil(t, tra.Broadcast(msg))

	rpcB := <-trb.Consume()
	b, err := ioutil.ReadAll(rpcB.Payload)
	assert.Nil(t, err)
	assert.Equal(t, b, msg)

	rpcC := <-trc.Consume()
	b, err = ioutil.ReadAll(rpcC.Payload)
	assert.Nil(t, err)
	assert.Equal(t, b, msg)
}
