package server

import (
	"fmt"
	"net"
	"sync"
	"time"

	"net-cat/internal/config"
)

// Server represents a chat server
type Server struct {
	port     int
	clients  map[string]*Client
	messages []Message
	mutex    sync.RWMutex
}

// NewServer creates a new server instance
func NewServer(port int) *Server {
	return &Server{
		port:     port,
		clients:  make(map[string]*Client),
		messages: make([]Message, 0),
	}
}

// Start starts the server
func (s *Server) Start() error {
	// Listen on all available network interfaces (0.0.0.0)
	listener, err := net.Listen(config.Protocol, fmt.Sprintf("0.0.0.0:%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	defer listener.Close()

	// Get all network interfaces to display available IP addresses
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("Warning: Could not get network interfaces: %v\n", err)
	} else {
		fmt.Println("Server is available on:")
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil { // Only show IPv4 addresses
					fmt.Printf("http://%s:%d\n", ipnet.IP.String(), s.port)
				}
			}
		}
	}

	fmt.Printf("Listening on port :%d\n", s.port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		s.mutex.RLock()
		clientCount := len(s.clients)
		s.mutex.RUnlock()

		if clientCount >= config.MaxClients {
			conn.Write([]byte(config.FullServerError + "\n"))
			conn.Close()
			continue
		}

		go s.handleConnection(conn)
	}
}

// Rest of the code remains the same...
func (s *Server) handleConnection(conn net.Conn) {
	client := NewClient(conn, s)

	// Send welcome banner
	conn.Write([]byte(config.WelcomeBanner))
	conn.Write([]byte(config.NamePrompt))

	// Get client name
	name, err := client.getName()
	if err != nil {
		conn.Write([]byte(err.Error() + "\n"))
		conn.Close()
		return
	}

	client.name = name

	s.mutex.Lock()
	s.clients[client.name] = client
	s.mutex.Unlock()

	s.sendHistory(client)

	joinMsg := Message{
		Time:    time.Now(),
		From:    "system",
		Content: client.name + config.JoinedMessage,
	}
	s.broadcast(joinMsg)

	go client.Read()
	go client.Write()
}

func (s *Server) broadcast(msg Message) {
	s.mutex.Lock()
	s.messages = append(s.messages, msg)
	clients := make([]*Client, 0, len(s.clients))
	for _, client := range s.clients {
		clients = append(clients, client)
	}
	s.mutex.Unlock()

	for _, client := range clients {
		select {
		case client.outgoing <- msg:
		default:
		}
	}
}

func (s *Server) removeClient(client *Client) {
	s.mutex.Lock()
	if _, exists := s.clients[client.name]; exists {
		leaveMsg := Message{
			Time:    time.Now(),
			From:    "system",
			Content: client.name + config.LeftMessage,
		}

		delete(s.clients, client.name)
		close(client.outgoing)
		client.conn.Close()

		s.mutex.Unlock()
		s.broadcast(leaveMsg)
		return
	}
	s.mutex.Unlock()
}

func (s *Server) sendHistory(client *Client) {
	s.mutex.RLock()
	history := make([]Message, len(s.messages))
	copy(history, s.messages)
	s.mutex.RUnlock()

	for _, msg := range history {
		client.outgoing <- msg
	}
}
