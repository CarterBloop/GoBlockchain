package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	Hash         []byte
	Transactions []*Transaction
	PrevHash     []byte
	Nonce        int
}

func CreateBlock(transactions []*Transaction, prevHash []byte) *Block {
	block := &Block{[]byte{}, transactions, prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func Genesis() *Block {
	return CreateBlock([]*Transaction{}, []byte{})
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	Handle(err)

	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	Handle(err)

	return &block
}

func (b *Block) ValidateBlock(prevHash []byte) bool {
	// Verify the proof of work
	pow := NewProof(b)
	if !pow.Validate() {
		return false
	}

	// Check the previous hash
	if !bytes.Equal(b.PrevHash, prevHash) {
		return false
	}

	// Validate all transactions in the block
	for _, tx := range b.Transactions {
		if !tx.ValidateTransaction() {
			return false
		}
	}

	return true
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}