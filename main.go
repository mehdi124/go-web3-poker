package main

import (
	"fmt"

	"github.com/mehdi124/go-web3-poker/deck"
)

func main() {

	card := deck.NewCard(deck.Spades, 1)
	fmt.Println(card)

}
