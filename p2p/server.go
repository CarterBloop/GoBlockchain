package p2p

import (
	"fmt"
	"net"
	"encoding/json"
	"io"
)

type Server struct {
	Port string
	Client *Client
}

/* Server */

func (server *Server) Start() {
	fmt.Println("[server] Starting server...")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.Port))
	if err != nil {
		fmt.Println("Failed to start server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server started on port", server.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}

		go server.handleConnection(conn)
	}
}

func (server *Server) handleConnection(conn net.Conn) {
	fmt.Println("[server] Handling connection...")
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break // Connection closed by the client
			}
			fmt.Println("Failed to read from connection:", err)
			return
		}

		var msg Message
		err = json.Unmarshal(buffer[:n], &msg)
		if err != nil {
			fmt.Println("Failed to unmarshal message:", err)
			return
		}

		server.handleMessage(msg)
	}
}


func (server *Server) handleMessage(msg Message) {
	fmt.Println("[server] Handling message...")
	switch msg.Type {
	// Peer List
	case MsgTypePeerList:
		var peers []Peer
		err := json.Unmarshal(msg.Data, &peers)
		if err != nil {
			fmt.Println("Failed to unmarshal peer list:", err)
			return
		}
		// Forward the received peer list to the client component for processing
		server.Client.updatePeerList(peers)
	default:
		fmt.Println("Unknown message type:", msg.Type)
	}
}