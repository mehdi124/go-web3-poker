package main

import (
	"fmt"

	"github.com/mehdi124/go-web3-poker"
)

func main() {

	card := go-web3-poker.deck.NewCard(deck.Spades, 1)
	fmt.Println(card)

}
