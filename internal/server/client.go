package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"net-cat/internal/config"
	"net-cat/internal/utils"
)

type Message struct {
	Time    time.Time
	From    string
	Content string
}

type Client struct {
	conn     net.Conn
	name     string
	server   *Server
	outgoing chan Message
	workers  int
}

// NewClient creates and initializes a new Client instance.
func NewClient(conn net.Conn, server *Server) *Client {
	c := &Client{
		conn:     conn,
		server:   server,
		outgoing: make(chan Message),
		workers:  10, // Number of concurrent message handlers
	}

	// Start worker pool
	for i := 0; i < c.workers; i++ {
		go c.messageWorker()
	}

	return c
}

func (c *Client) messageWorker() {
	for msg := range c.outgoing {
		// Process message
		data := fmt.Sprintf("[%s][%s]: %s\n", msg.Time.Format("2006-01-02 15:04:05"), msg.From, msg.Content)
		c.conn.Write([]byte(data))
	}
}

// Write continuously sends messages from the outgoing channel to the client.
func (c *Client) getName() (string, error) {
	reader := bufio.NewReader(c.conn)

	for {
		name, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("failed to read name")
		}

		originalName := strings.TrimSpace(name)
		if !utils.ValidateName(originalName) {
			c.conn.Write([]byte(config.EmptyNameError + "\n"))
			c.conn.Write([]byte(config.NamePrompt))
			continue
		}

		normalizedName := utils.NormalizeName(originalName)

		c.server.mutex.RLock()
		exists := false
		for existingName := range c.server.clients {
			if utils.NormalizeName(existingName) == normalizedName {
				exists = true
				break
			}
		}
		c.server.mutex.RUnlock()

		if exists {
			c.conn.Write([]byte("This name is already taken\n"))
			c.conn.Write([]byte(config.NamePrompt))
			continue
		}

		return originalName, nil
	}
}

// Write continuously sends messages from the outgoing channel to the client.
func (c *Client) Read() {
	defer c.server.removeClient(c)

	reader := bufio.NewReader(c.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		message = strings.TrimSpace(message)
		if message == "" {
			continue
		}

		msg := Message{
			Time:    time.Now(),
			From:    c.name,
			Content: message,
		}

		c.server.broadcast(msg)
	}
}

// Write continuously sends messages from the outgoing channel to the client.
func (c *Client) Write() {
	for msg := range c.outgoing {
		formattedMsg := utils.FormatMessage(
			msg.Time.Format(config.TimeFormat),
			msg.From,
			msg.Content,
		)

		if _, err := c.conn.Write([]byte(formattedMsg)); err != nil {
			return
		}
	}
}
