package p2p

import (
	"log"
	"net"
	"sync"
)

// TCPPeer represents a remote node over a TCP connection
type TCPPeer struct {
	// conn is the underlying connection of the TCPPeer
	conn net.Conn

	// outbound represents the TCPPeer dialing out or not
	// true if dialing out a connection, false if accepting a connection.
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransportOpts struct {
	ListenAddress string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
}

type TCPTransport struct {
	TCPTransportConfig TCPTransportOpts
	listener           net.Listener

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportConfig: opts,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.TCPTransportConfig.ListenAddress)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()
	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			log.Printf("TCP Accept Error: %s\n", err)
		}

		log.Printf("New incoming connection: %+v\n", conn)

		go t.handleConn(conn)

	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	if err := t.TCPTransportConfig.HandshakeFunc(peer); err != nil {
		conn.Close()
		log.Printf("TCP handshake error: %s\n", err)
		return
	}

	// Read loop
	msg := &Message{}
	for {
		if err := t.TCPTransportConfig.Decoder.Decode(conn, msg); err != nil {
			log.Printf("TCP error: %s\n", err)
			continue
		}
		msg.From = conn.RemoteAddr()
		log.Printf("message: %+v\n", msg)
	}
}
