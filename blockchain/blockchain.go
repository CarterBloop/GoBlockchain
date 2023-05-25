package blockchain

import (
	"fmt"
	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp/blocks"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database *badger.DB
}

func ChainInit() *BlockChain {
	var lastHash []byte

	// Create a new BadgerDB instance
	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	Handle(err)

	// Add Genesis Block to DB
	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")
			genesis := GenesisBlock()
			fmt.Println("Genesis proved")
			err = txn.Set(genesis.Hash, genesis.Serialize())
			Handle(err)
			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash

			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			Handle(err)
			lastHash, err = item.Value(func(val []byte))
			return err
		}
	})

	Handle(err)

	// Create a new BlockChain instance
	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

func(chain *BlockChain) AddBlock(data string) {
	var lastHash []byte

	// Get Last Hash from DB
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.Value(func(val []byte))
		return err
	})
	Handle(err)

	// Create New Block
	newBlock := CreateBlock(data, lastHash)

	// Update Last Hash in DB
	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)
		chain.LastHash = newBlock.Hash
		return err
	})
	Handle(err)
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	// New BlockChainIterator instance (starts at the last block)
	iter := &BlockChainIterator{chain.LastHash, chain.Database}
	return iter
}

func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		Handle(err)
		encodedBlock, err := item.Value(func(val []byte))
		block = Deserialize(encodedBlock)
		return err
	})
	Handle(err)

	iter.CurrentHash = block.PrevHash
	return block
}