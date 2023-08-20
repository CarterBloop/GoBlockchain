package main

import (
	"os"
	"bufio"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"fmt"

	"GoBlockchain/blockchain"
)

type CommandLine struct {
	blockchain *blockchain.BlockChain
}

func (cli *CommandLine) printChain() {
	fmt.Println(cli.blockchain.ToString())
}

func main() {
	chain := blockchain.InitBlockChain()
	defer chain.Database.Close()

	go chain.MineMempool()

	reader := bufio.NewReader(os.Stdin)

	var privateKey *ecdsa.PrivateKey
	var publicKey ecdsa.PublicKey

	for {
		fmt.Println("Commands:")
		fmt.Println("1) createaccount - Create a new account")
		fmt.Println("2) vote - Vote on a proposal")
		fmt.Println("3) print - Print the blockchain")
		fmt.Println("4) exit - Exit the program")
		fmt.Print("Enter command: ")

		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)

		switch cmd {
		case "createaccount":
			privateKey, _ = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			publicKey = privateKey.PublicKey
			fmt.Println("Account created!")
			fmt.Println("Voter ID (Public Key):", hex.EncodeToString(elliptic.Marshal(elliptic.P256(), publicKey.X, publicKey.Y)))

		case "vote":
			if privateKey == nil {
				fmt.Println("Please create an account first.")
				continue
			}
			fmt.Print("Enter proposal: ")
			proposal, _ := reader.ReadString('\n')
			proposal = strings.TrimSpace(proposal)
			signature, _ := ecdsa.SignASN1(rand.Reader, privateKey, []byte(proposal))
			transaction := blockchain.NewTransaction(hex.EncodeToString(elliptic.Marshal(elliptic.P256(), publicKey.X, publicKey.Y)), proposal, signature)
			chain.AddTransactionToMempool(transaction)
			fmt.Println("Vote added to mempool!")

		case "print":
			cli := CommandLine{chain}
			cli.printChain()

		case "exit":
			os.Exit(0)

		default:
			fmt.Println("Invalid command. Please try again.")
		}
	}
}