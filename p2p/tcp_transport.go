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

type TCPTransport struct {
	listenAddress string
	listener      net.Listener

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{listenAddress: listenAddr}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.listenAddress)
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
			continue
		}

		go t.handleConn(conn)

	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	log.Printf("New incoming connection: %+v\n", peer)
}
