package main

import (
	"bytes"
	"go-blockchain/core"
	"go-blockchain/crypto"
	"go-blockchain/network"
	"math/rand"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	trLocal := network.NewLocalTransport("Local") //trLocal是Transport接口，所以要实现该接口所有方法
	trRemote := network.NewLocalTransport("REMOTE")

	trLocal.Connect(trRemote)
	trRemote.Connect(trLocal)

	// var wg sync.WaitGroup
	// wg.Add(2) // 增加等待的协程数量
	// go func() {
	// 	defer wg.Done() // 协程完成任务后通知 WaitGroup
	// 	s.Start()
	// }()
	go func() {
		// defer wg.Done() // 协程完成任务后通知 WaitGroup
		for {
			// trRemote.SendMessage(trLocal.Addr(), []byte("hello world"))

			if err := sendTransaction(trRemote, trLocal.Addr()); err != nil {
				logrus.Error(err)
			}
			time.Sleep(1 * time.Second)
			// // 直接读取 s.rpcCh 的长度
			// length := len(s.GetRPCCh())
			// fmt.Printf("s.rpcCh 的长度: %d\n", length)
		}
	}()
	opts := network.ServerOpts{
		Transports: []network.Transport{trLocal}, //远程给本地发消息
	}
	s := network.NewServer(opts)
	// wg.Wait()
	s.Start()
}

func sendTransaction(tr network.Transport, to network.NetAddr) error {
	privKey := crypto.GeneratePrivateKey()
	data := []byte(strconv.FormatInt(int64(rand.Intn(10)), 10))
	tx := core.NewTransaction(data)
	tx.Sign(privKey)
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())

	return tr.SendMessage(to, msg.Bytes())
}
