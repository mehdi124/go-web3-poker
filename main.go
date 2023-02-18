package main

import (
	//"github.com/mehdi124/go-web3-poker/deck"
	"fmt"
	"time"

	"github.com/mehdi124/go-web3-poker/deck"
	"github.com/mehdi124/go-web3-poker/p2p"
)

func main() {

	cfg := &p2p.ServerConfig{
		Version:     "mehdi124_V0.1-alpha",
		ListenAddr:  ":3000",
		GameVariant: p2p.TexasHoldem,
	}

	server := p2p.NewServer(*cfg)
	go server.Start()

	time.Sleep(2 * time.Second)

	remoteCfg := &p2p.ServerConfig{
		Version:     "mehdi124_V0.1-alpha",
		ListenAddr:  ":4000",
		GameVariant: p2p.TexasHoldem,
	}

	remoteServer := p2p.NewServer(*remoteCfg)
	if err := remoteServer.Connect(":3000"); err != nil {
		panic(err)
	}

	go remoteServer.Start()

	fmt.Println(deck.New())

}
