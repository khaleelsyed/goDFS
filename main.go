package main

import (
	"fmt"
	"log"

	"github.com/khaleelsyed/goDFS/p2p"
)

func main() {
	tcpOpts := &p2p.TCPTransportOpts{
		ListenAddress: ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	var tr p2p.Transport = p2p.NewTCPTransport(*tcpOpts)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("%+v\n", msg)
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
