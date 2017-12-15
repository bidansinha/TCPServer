package tcp_server

import (
//	"bufio"
	"log"
	"net"
	"io"
	"fmt"
	//"time"
)

// Client holds info about connection
type Client struct {
	conn   net.Conn
	Server *server
}

// TCP server
type server struct {
	address                  string // Address to open connection: localhost:9999
	onNewClientCallback      func(c *Client)
	onClientConnectionClosed func(c *Client, err error)
	onNewMessage             func(c *Client, message string)
	onNewMessages			func(c *Client, byte []byte, size int)
}

// Read client data from channel
func (c *Client) listen() {

	buf := make([]byte, 0, 4096) // big buffer
	tmp := make([]byte, 1024)     // using small tmo buffer for demonstrating
	size := 0

	defer c.Close();

	//c.conn.SetReadDeadline(time.Now().Add(100000))
	for {

		n, err := c.conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
				c.Close();
				return
			}
			//fmt.Println(" error breaking ", err)
			//c.Close()
			break
		} else {
			//fmt.Println("got", n, "bytes.")
			buf = append(buf, tmp[:n]...)
		}
		size = size + n
	}
	c.Server.onNewMessages(c, buf, size)
	//reader := bufio.NewReader(c.conn)
	//for {
	//	message, err := reader.ReadString('\n')
	//	if err != nil {
	//		c.conn.Close()
	//		c.Server.onClientConnectionClosed(c, err)
	//		return
	//	}
	//	c.Server.onNewMessage(c, message)
	//}
}

// Send text message to client
func (c *Client) Send(message string) error {
	_, err := c.conn.Write([]byte(message))
	return err
}

// Send bytes to client
func (c *Client) SendBytes(b []byte) error {
	_, err := c.conn.Write(b)
	return err
}

func (c *Client) Conn() net.Conn {
	return c.conn
}

func (c *Client) Close() error {
	c.Server.onClientConnectionClosed(c, nil)
	return c.conn.Close()
}

// Called right after server starts listening new client
func (s *server) OnNewClient(callback func(c *Client)) {
	s.onNewClientCallback = callback
}

// Called right after connection closed
func (s *server) OnClientConnectionClosed(callback func(c *Client, err error)) {
	s.onClientConnectionClosed = callback
}

// Called when Client receives new message
func (s *server) OnNewMessage(callback func(c *Client, message string)) {
	s.onNewMessage = callback
}

func (s *server) OnNewMessages(callback func(c *Client, bytes [] byte, size int)) {
	s.onNewMessages = callback
}

// Start network server
func (s *server) Listen() {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatal("Error starting TCP server.")
	}
	defer listener.Close()

	for {
		conn, _ := listener.Accept()
		client := &Client{
			conn:   conn,
			Server: s,
		}
		go client.listen()
		s.onNewClientCallback(client)
	}
}

// Creates new tcp server instance
func New(address string) *server {
	log.Println("Creating server with address", address)
	server := &server{
		address: address,
	}

	server.OnNewClient(func(c *Client) {})
	server.OnNewMessage(func(c *Client, message string) {})
	server.OnClientConnectionClosed(func(c *Client, err error) {})
	server.OnNewMessages(func(c *Client, byte[] byte, size int) {})

	return server
}
