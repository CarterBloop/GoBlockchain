package node

import (
	"GoBlockchain/blockchain"
	"GoBlockchain/p2p"
)

type Node struct {
	Client *p2p.Client
	Server *p2p.Server
	Chain  *blockchain.BlockChain
	// New Block channel
	NewBlockChan chan *blockchain.Block
}

func NewNode(client *p2p.Client, server *p2p.Server, chain *blockchain.BlockChain) *Node {
	return &Node{
		Client: client,
		Server: server,
		Chain:  chain,
	}
}

func (node *Node) Start() {
	go node.Server.Start()
	go node.Client.Start()
	go node.Chain.MineMempool()
}