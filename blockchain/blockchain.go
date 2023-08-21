package blockchain

import (
	"fmt"
	"log"
	"os"
	"sync"
	"strings"
	"strconv"

	"github.com/dgraph-io/badger"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func InitBlockChain(dataDir string) *BlockChain {
	var lastHash []byte

	opts := badger.DefaultOptions(dataDir)

	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		err := os.MkdirAll(dataDir, 0755)
		if err != nil {
			log.Fatal("Error Creating Dir: ", dataDir, " ", err)
		}
	}
	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")
			genesis := Genesis()
			fmt.Println("Genesis proved")
			err = txn.Set(genesis.Hash, genesis.Serialize())
			Handle(err)
			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash

			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			Handle(err)
			err = item.Value(func(val []byte) error {
				lastHash = append([]byte{}, val...)
				return nil
			})
			return err
		}
	})

	Handle(err)

	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

func (chain *BlockChain) AddBlock(transactions []*Transaction) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = append([]byte{}, val...)
			return nil
		})

		return err
	})
	Handle(err)

	newBlock := CreateBlock(transactions, lastHash)

	// Validate the new block
	if !newBlock.ValidateBlock(lastHash) {
		log.Println("Invalid block")
		return
	}

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
	iter := &BlockChainIterator{chain.LastHash, chain.Database}

	return iter
}

func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			block = Deserialize(val)
			return nil
		})
		return err
	})
	Handle(err)

	iter.CurrentHash = block.PrevHash

	return block
}

func (chain *BlockChain) ToString() string {
	var result strings.Builder

	iter := chain.Iterator()

	for {
		block := iter.Next()

		result.WriteString(fmt.Sprintf("Prev. hash: %x\n", block.PrevHash))
		var i = 1
		for _, tx := range block.Transactions {
			result.WriteString(fmt.Sprintf("  (%d)\n",i))
			i++
			result.WriteString(fmt.Sprintf("    - Transaction ID: %s\n", tx.ID))
			result.WriteString(fmt.Sprintf("    - Voter ID: %s\n", tx.VoterID))
			result.WriteString(fmt.Sprintf("    - Proposal: %s\n", tx.Proposal))
			result.WriteString(fmt.Sprintf("    - Signature: %x\n", tx.Signature))
		}
		result.WriteString(fmt.Sprintf("Hash: %x\n", block.Hash))
		pow := NewProof(block)
		result.WriteString(fmt.Sprintf("PoW: %s\n", strconv.FormatBool(pow.Validate())))
		result.WriteString("\n")

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return result.String()
}

// Mempool

var Mempool = make([]*Transaction, 0)
var MempoolLock = sync.Mutex{}

func (chain *BlockChain) AddTransactionToMempool(transaction *Transaction) {
	MempoolLock.Lock()
	Mempool = append(Mempool, transaction)
	MempoolLock.Unlock()
}

const MempoolSizeThreshold = 10

func (chain *BlockChain) MineMempool() {
	for {
		MempoolLock.Lock()
		if len(Mempool) < MempoolSizeThreshold {
			MempoolLock.Unlock()
			continue
		}
		transactions := Mempool[:]
		Mempool = make([]*Transaction, 0)
		MempoolLock.Unlock()

		chain.AddBlock(transactions)
		fmt.Println("Mined new block with", len(transactions), "transactions")
	}
}