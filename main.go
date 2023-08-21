package main

import (
	"flag"
	"os"
	"bufio"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"fmt"

	"GoBlockchain/blockchain"
	"GoBlockchain/p2p"
	"GoBlockchain/node"
)

type CommandLine struct {
	blockchain *blockchain.BlockChain
}

func (cli *CommandLine) printChain() {
	fmt.Println(cli.blockchain.ToString())
}

func main() {
	// Flags
	entryIP := flag.String("entry-ip", "", "IP address of the entry node")
	entryPort := flag.String("entry-port", "", "Port of the entry node")
	serverPort := flag.String("server-port", "", "Port of the current node")
	dataDir := flag.String("data-dir", "", "Data directory for the blockchain")
	flag.Parse()

	if *dataDir == "" {
		*dataDir = os.TempDir()
	}

	if *entryIP == "" || *entryPort == "" || *serverPort == "" || *dataDir == "" {
		fmt.Println("Missing flags. Please try again.")
		return
	}

	// Init client and server
	client := p2p.Client{
		Peers: []p2p.Peer{
			{IP: *entryIP, Port: *entryPort},
			{IP: "localhost", Port: *serverPort},
		},
	}
	server := p2p.Server{
		Port:   *serverPort,
		Client: &client,
	}

	// Init the blockchain
	chain := blockchain.InitBlockChain(*dataDir)
	defer chain.Database.Close()

	// Create the Node object
	node := node.NewNode(&client, &server, chain)

	// Start the Node
	node.Start()

	// Init the reader
	reader := bufio.NewReader(os.Stdin)

	// Init node account
	var privateKey *ecdsa.PrivateKey
	var publicKey ecdsa.PublicKey

	// Start the CLI
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