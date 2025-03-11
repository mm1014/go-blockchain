package main

import (
	"fmt"
	"go-blockchain/network"
	"sync"
	"time"
)

func main() {
	trLocal := network.NewLocalTransport("Local") //trLocal是Transport接口，所以要实现该接口所有方法
	trRemote := network.NewLocalTransport("REMOTE")

	trLocal.Connect(trRemote)
	trRemote.Connect(trLocal)

	opts := network.ServerOpts{
		Transports: []network.Transport{trLocal}, //远程给本地发消息
	}
	s := network.NewServer(opts)
	var wg sync.WaitGroup
	wg.Add(2) // 增加等待的协程数量
	go func() {
		defer wg.Done() // 协程完成任务后通知 WaitGroup
		s.Start()
	}()
	go func() {
		defer wg.Done() // 协程完成任务后通知 WaitGroup
		for {
			trRemote.SendMessage(trLocal.Addr(), []byte("hello world"))
			time.Sleep(1 * time.Second)
			// 直接读取 s.rpcCh 的长度
			length := len(s.GetRPCCh())
			fmt.Printf("s.rpcCh 的长度: %d\n", length)
		}
	}()
	wg.Wait()

}
