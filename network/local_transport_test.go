package network

import (
	"io"
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
	msg := []byte("Hello world!")
	assert.Nil(t, tra.SendMessage(trb.Addr(), msg))
	//从trb的通道里面读一个rpc出来
	rpc := <-trb.Consume()
	b, err := io.ReadAll(rpc.Payload)
	assert.Nil(t, err)
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
	//Payload不是[]byte  所以不能直接赋值
	b, err := io.ReadAll(rpcB.Payload)
	assert.Nil(t, err)
	assert.Equal(t, b, msg)

	rpcC := <-trc.Consume()
	b, err = io.ReadAll(rpcC.Payload)
	assert.Nil(t, err)
	assert.Equal(t, b, msg)
}
