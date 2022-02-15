package tcp_server

import (
	"bufio"
	"net"
	"os"
	"time"
)

type Client struct {
	lastSeen int64
	conn     net.Conn
	Server   *server
}

// Read client data from channel
func (c *Client) listen() {
	var delimiter string
	var ok bool
	if delimiter, ok = os.LookupEnv("DELIMITER"); !ok {
		delimiter = "\n"
	}

	c.Server.onNewClientCallback(c)
	reader := bufio.NewReader(c.conn)
	for {
		message, err := reader.ReadString(byte(delimiter[0]))
		if err != nil {
			c.conn.Close()
			c.Server.onClientConnectionClosed(c, err)
			return
		}
		c.lastSeen = time.Now().Unix()
		c.Server.onNewMessage(c, message)
	}
}

// Send text message to client
func (c *Client) Send(message string) error {
	return c.SendBytes([]byte(message))
}

// Send bytes to client
func (c *Client) SendBytes(b []byte) error {
	_, err := c.conn.Write(b)
	if err != nil {
		c.conn.Close()
		c.Server.onClientConnectionClosed(c, err)
	}
	return err
}

func (c *Client) Conn() net.Conn {
	return c.conn
}

func (c *Client) Close() error {
	return c.conn.Close()
}
