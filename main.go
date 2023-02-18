package main

import (
	//"github.com/mehdi124/go-web3-poker/deck"
	"github.com/mehdi124/go-web3-poker/p2p"
)

func main() {

	cfg := &p2p.ServerConfig{
		ListenAddr: ":3000",
	}

	server := p2p.NewServer(*cfg)
	server.Start()

	//	card := deck.New()

}
