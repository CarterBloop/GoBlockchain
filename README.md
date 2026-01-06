# GoBlockchain

GoBlockchain is a simple blockchain-based voting system implemented in Go. It allows users to create accounts, vote on proposals, and view the blockchain. The system uses a peer-to-peer network to propagate and synchronize data across multiple nodes, ensuring a decentralized and secure voting process.

## Features

- Blockchain-based Voting: Users can create accounts and vote on proposals. Each vote is recorded as a transaction on the blockchain, ensuring transparency and immutability.
- Proof of Work: The system uses a proof-of-work algorithm to validate and mine new blocks, adding an additional layer of security.
- Peer-to-Peer Network: Nodes in the network communicate with each other to propagate and synchronize data, ensuring a decentralized system.
- Mempool: Transactions are temporarily stored in a mempool before being mined into a new block, allowing for efficient transaction processing.
- CLI Interface: Users can interact with the system through a simple command-line interface.

## Local Testnet

```shell
git clone github.com/carterbloop/goblockchain
cd goblockchain
go run main.go -data-dir tmp -server-port 8080 -entry-ip localhost -entry-port 8081
go run main.go -data-dir tmp2 -server-port 8081 -entry-ip localhost -entry-port 8080
go run main.go -data-dir tmp3 -server-port 8082 -entry-ip localhost -entry-port 8081
```
