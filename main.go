package main

import (
	"log"
	"os"
	"pbft-impl-go/consensus"
	"pbft-impl-go/network"
	"sync"
)

// Hard-coded for test.
var viewID = int64(10000000000)
var nodeTableForTest = []*network.NodeInfo{
	{NodeID: "Node1", Url: "localhost:49152"},
	{NodeID: "Node2", Url: "localhost:49153"},
	{NodeID: "Node3", Url: "localhost:49154"},
	{NodeID: "Node4", Url: "localhost:49155"},
}

var NodesPS = make(map[string]map[string]string)

func main() {
	var nodeTable []*network.NodeInfo
	nodeTable = nodeTableForTest
	for _, nodeInfo := range nodeTable {
		// 生成公私钥
		sk, pk := consensus.BLSMgr.GenerateKey()
		var nodePS = make(map[string]string)
		nodePS["PublicKey"] = pk.Compress().String()
		nodePS["PrivateKey"] = sk.Compress().String()
		NodesPS[nodeInfo.NodeID] = nodePS
		nodeInfo.PubKey = pk
	}

	var wg sync.WaitGroup // 使用sync.WaitGroup来等待所有goroutine完成

	for _, nodeInfo := range nodeTable {
		wg.Add(1)                             // 在启动goroutine之前增加WaitGroup的计数器
		go func(nodeInfo *network.NodeInfo) { // 用nodeInfo作为参数启动goroutine
			defer wg.Done() // 确保在goroutine结束时调用wg.Done()
			decodePrivKey, err := consensus.BLSMgr.DecSecretKeyHex(NodesPS[nodeInfo.NodeID]["PrivateKey"])
			if err != nil {
				AssertError(err) // 确保你有一种方式来处理错误
				return
			}
			server := network.NewServer(nodeInfo.NodeID, nodeTable, viewID, decodePrivKey)
			if server != nil {
				server.Start() // 这将在其自己的goroutine中执行
			}
		}(nodeInfo) // 将当前循环的nodeInfo作为参数传递给goroutine
	}
	wg.Wait() // 等待所有goroutine完成
}

func AssertError(err error) {
	if err == nil {
		return
	}

	log.Println(err)
	os.Exit(1)
}
