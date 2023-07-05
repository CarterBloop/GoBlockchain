# GoBlockchain
## Introduction
GoBlockchain is a simple implementation of a blockchain in Go. It provides an accessible platform to understand the basic working principles of a blockchain system.

## Features
- Creation of a genesis block upon initialization.
- Addition of new blocks with user-specified data.
- Implementation of a proof-of-work system, validating each new block.
- Serialization of the blockchain to persist data.
- Command-line interface for interaction.

## How-to
```shell
git clone github.com/carterbloop/goblockchain
cd goblockchain
go run main.go add -block <your text here>
go run main.go print
```