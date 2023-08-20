package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/rand"
	"fmt"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
)

type Transaction struct {
	ID        string
	VoterID   string
	Proposal  string
	Signature []byte
}

func NewTransaction(voterID, proposal string, signature []byte) *Transaction {
	tx := &Transaction{
		ID:        GenerateTxID(),
		VoterID:   voterID,
		Proposal:  proposal,
		Signature: signature,
	}

	return tx
}

func (tx *Transaction) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(tx)
	Handle(err)

	return res.Bytes()
}

func DeserializeTransaction(data []byte) *Transaction {
	var tx Transaction

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&tx)
	Handle(err)

	return &tx
}

func GenerateTxID() string {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		log.Fatal(err)
	}

	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

func (tx *Transaction) ValidateTransaction() bool {
	// Check the transaction format
	if tx.ID == "" || tx.VoterID == "" || tx.Proposal == "" || len(tx.Signature) == 0 {
		return false
	}

	// Verify the signature
	publicKeyBytes, err := hex.DecodeString(tx.VoterID)
	if err != nil {
		return false
	}

	curve := elliptic.P256()
	x := big.Int{}
	y := big.Int{}
	keyLen := len(publicKeyBytes)
	x.SetBytes(publicKeyBytes[:(keyLen / 2)])
	y.SetBytes(publicKeyBytes[(keyLen / 2):])

	rawPublicKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}

	r := big.Int{}
	s := big.Int{}
	sigLen := len(tx.Signature)
	r.SetBytes(tx.Signature[:(sigLen / 2)])
	s.SetBytes(tx.Signature[(sigLen / 2):])

	hash := sha256.Sum256([]byte(tx.ID + tx.VoterID + tx.Proposal))

	return ecdsa.Verify(&rawPublicKey, hash[:], &r, &s)
}
