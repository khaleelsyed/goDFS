package p2p

// Peer represents an interface for remote nodes.
type Peer interface {
	Close() error
}

// Transport handles the communication betweeen the nodes in the network.
// This can be TCP, UDP, websockets, ...
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
}
