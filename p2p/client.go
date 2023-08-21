package p2p

import (
	"fmt"
	"net"
	"encoding/json"
	"time"
)

type Peer struct {
	IP   string
	Port string
}

type Client struct {
	Peers []Peer
}

const (
	MsgTypePeerList      = "peer_list"
)

type Message struct {
	Type string
	Data []byte
}

/* Client */

func (client *Client) Start() {
	fmt.Println("[client] Starting...")
	for {
		fmt.Println("[client] Dialing peers...")
		for _, peer := range client.Peers {
			conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", peer.IP, peer.Port))
			if err != nil {
				fmt.Println("Failed to connect to peer:", err)
				continue
			}
			defer conn.Close()

			go client.handleConnection(conn)
		}
		time.Sleep(10 * time.Second)
	}
}

func (client *Client) handleConnection(conn net.Conn) {
	fmt.Println("[client] Handling connection...")

	client.sharePeerList(conn)

	// Listen for incoming messages
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Failed to read from connection:", err)
			return
		}

		var msg Message
		err = json.Unmarshal(buffer[:n], &msg)
		if err != nil {
			fmt.Println("Failed to unmarshal message:", err)
			return
		}

		client.handleMessage(msg)
	}
}

func (client *Client) handleMessage(msg Message) {
	fmt.Println("[client] Handling msg...")
	switch msg.Type {
	// Peer List
	case MsgTypePeerList:
		var peers []Peer
		err := json.Unmarshal(msg.Data, &peers)
		if err != nil {
			fmt.Println("Failed to unmarshal peer list:", err)
			return
		}
		client.updatePeerList(peers)

	// Unknown
	default:
		fmt.Println("Unknown message type:", msg.Type)
	}
}

/* Peer Discovery */

func (client *Client) sharePeerList(conn net.Conn) {
	fmt.Println("[client] Sharing peer list...")
	fmt.Println("[client] Peers:", client.Peers)
	peerList, err := json.Marshal(client.Peers)
	if err != nil {
		fmt.Println("Failed to marshal peer list:", err)
		return
	}

	msg := Message{
		Type: MsgTypePeerList,
		Data: peerList,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Failed to marshal message:", err)
		return
	}

	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("Failed to send peer list:", err)
		return
	}
}

func (client *Client) updatePeerList(peers []Peer) {
	fmt.Println("[client] Updating peer list...")
	peerMap := make(map[string]bool)
	for _, peer := range client.Peers {
		peerMap[fmt.Sprintf("%s:%s", peer.IP, peer.Port)] = true
	}

	for _, peer := range peers {
		if !peerMap[fmt.Sprintf("%s:%s", peer.IP, peer.Port)] {
			client.Peers = append(client.Peers, peer)
		}
	}
}