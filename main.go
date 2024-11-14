package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	mrand "math/rand"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
	"github.com/multiformats/go-multiaddr"
)

func main() {
	seed := int64(0)
	port := 0

	flag.Int64Var(&seed, "seed", 0, "Seed value for generating a PeerID, 0 is random")
	flag.IntVar(&port, "port", 0, "")
	flag.Parse()

	_, cancel := context.WithCancel(context.Background())

	var r io.Reader
	if seed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(seed))
	}

	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		log.Fatalf("failed to create priv key %v\n", err)
	}

	addr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))

	h, err := libp2p.New(
		libp2p.ListenAddrs(addr),
		libp2p.Identity(priv),
	)
	if err != nil {
		log.Fatalf("failed to create nodehost %v\n", err)
	}

	_, err = relay.New(h)
	if err != nil {
		log.Fatalf("failed to instantiate relay %v\n", err)
	}

	for _, addr := range h.Addrs() {
		log.Printf("  %s/p2p/%s", addr, h.ID().String())
	}

	waitSignal(h, cancel)
}

func waitSignal(h host.Host, cancel func()) {
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-c

	cancel()

	if err := h.Close(); err != nil {
		panic(err)
	}
	os.Exit(0)
}
