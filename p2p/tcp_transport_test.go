package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	listenAddr := ":3000"
	opts := TCPTransportOpts{
		ListenAddress: listenAddr,
		HandshakeFunc: NOPHandshakeFunc,
		Decoder:       DefaultDecoder{},
	}
	tr := NewTCPTransport(opts)

	assert.Equal(t, tr.TCPTransportConfig.ListenAddress, listenAddr)

	assert.Nil(t, tr.ListenAndAccept())
}
