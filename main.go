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

//	var NodesPS = map[string]map[string]string{
//		"Node1": {
//			"PrivateKey": "7388ab5ea6442ea952c3734c7644bcf677286418b55607398343a2cb2e94393b",
//			"PublicKey":  "b165511a1cd1d8d110ad47c0ac8fb4d340a028b925a7630d7c803d4be3721a02a4766ff2257c9f4bb866e81a1612a5df05e4a876adc20cbc3bee593efe170827e17c84e879c64f324e49545eeae62ea4cb4693039eba0e390ec91cf110e9df18",
//		},
//		"Node2": {
//			"PrivateKey": "5abc593e02780a836bbbeaf4fc5a467a0a25f8d65c71c3cf5e73087a5b9eb4e5",
//			"PublicKey":  "a77e0cf38d5c4542a740da858083781cd5ed09b1c1c2b8a5acb66c8eeb7e6498b9561fb6e795dfa5001019181b4d9c7a0f15d800aff57660d7a62205aa05a1fbf7157060d5e514bc3e4fe6933576d2a46ccabfa6928ef4dbbbb4d4f132d7f844",
//		},
//		"Node3": {
//			"PrivateKey": "00174bfd12b1121ace57f2f3bd40ad006061542a89d5230a1eb5fc5c33fe6bd8",
//			"PublicKey":  "88b2adb7c1d71dea05a2e9c30482a23068ff0698e91e70ca9ce373da8c999ff855953393dc919e42ea58d48d03f2c5f90e3af642ecfe12f6d9f31cef6c1f774c60ed9482e19f55df3ec62c383dc5eb542f3213161f988d85f514580249795c45",
//		},
//		"Node4": {
//			"PrivateKey": "3bbb2dacb81b9945df7e0ff13b40792a10e0d81fe7b0e92ca2bd347ed2dc7f33",
//			"PublicKey":  "b44bd99d1c6c9cd858f136aa92f6a93205b546642873b47dfed403488e7af39771a09137c58413c90cdd2ba004dc71210bc6b71bd6204a6fdfab95d84654de6cc2587c4f90bfd7820f1489d962843a079863f7073f86d1fc50f8411b00d26ded",
//		},
//	}
var NodesPS = make(map[string]map[string]string)

func main() {
	var nodeTable []*network.NodeInfo
	//if len(os.Args) < 2 {
	//	fmt.Println("Usage:", os.Args[0], "<nodeID> [node.list]")
	//	return
	//}
	//
	//nodeID := os.Args[1]
	//if len(os.Args) == 2 {
	//	fmt.Println("Node list are not specified")
	//	fmt.Println("Embedded list is used for test")
	//	nodeTable = nodeTableForTest
	//} else {
	//	nodeListFile := os.Args[2]
	//	jsonFile, err := os.Open(nodeListFile)
	//	AssertError(err)
	//	defer jsonFile.Close()
	//
	//	err = json.NewDecoder(jsonFile).Decode(&nodeTable)
	//	AssertError(err)
	//}
	//
	//for _, nodeInfo := range nodeTable {
	//	PublicKey, err := consensus.BLSMgr.DecPublicKeyHex(NodesPS[nodeInfo.NodeID]["PublicKey"])
	//	AssertError(err)
	//	nodeInfo.PubKey = PublicKey
	//}
	//decodePrivKey, err := consensus.BLSMgr.DecSecretKeyHex(NodesPS[nodeID]["PrivateKey"])
	//if err != nil {
	//	AssertError(err)
	//}
	//server := network.NewServer(nodeID, nodeTable, viewID, decodePrivKey)
	//if server != nil {
	//	server.Start()
	//}

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
	//for _, nodeInfo := range nodeTable {
	//	decodePrivKey, err := consensus.BLSMgr.DecSecretKeyHex(NodesPS[nodeInfo.NodeID]["PrivateKey"])
	//	if err != nil {
	//		AssertError(err)
	//	}
	//	server := network.NewServer(nodeInfo.NodeID, nodeTable, viewID, decodePrivKey)
	//	if server != nil {
	//		server.Start()
	//	}
	//}

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
