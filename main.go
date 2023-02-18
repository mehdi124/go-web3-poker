package main

import (
	//"github.com/mehdi124/go-web3-poker/deck"

	"log"
	"time"

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

	go remoteServer.Start()

	if err := remoteServer.Connect(":3000"); err != nil {
		log.Fatal(err)
	}

	//	fmt.Println(deck.New())

	select {}
}
