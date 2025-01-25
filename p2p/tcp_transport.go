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

// Close implements the Peer interface
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

type TCPTransportOpts struct {
	ListenAddress string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
}

type TCPTransport struct {
	TCPTransportConfig TCPTransportOpts
	rpcCh              chan RPC
	listener           net.Listener

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportConfig: opts,
		rpcCh:              make(chan RPC),
	}
}

// Consume implements the Transport interface,
// returning a read-only channel for reading the
// incoming messages received from another peer
// in the network.
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcCh
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
	var peer Peer = NewTCPPeer(conn, true)

	if err := t.TCPTransportConfig.HandshakeFunc(peer); err != nil {
		conn.Close()
		log.Printf("TCP handshake error: %s\n", err)
		return
	}

	// Read loop
	rpc := RPC{}
	for {
		if err := t.TCPTransportConfig.Decoder.Decode(conn, &rpc); err != nil {
			log.Printf("TCP error: %s\n", err)
			continue
		}
		rpc.From = conn.RemoteAddr()
		t.rpcCh <- rpc
	}
}
