package "deck"


import (
	"fmt"
)


type Suit int

func (s Suit) String() string {

	switch s {
	case Spades:
		return "SPADES"
	case Harts:
		return "HARTES"
	case Diamonds:
		return "DIAMONDS"
	case Clubs:
		return "CLUBS"
	default:
		panic("invalid card suit")
	}
}

const (

	Spades Suit = iota
	Harts
	Diamonds
	Clubs
)


type Card struct{
	value int
	suit Suit
}

func 

func NewCard(s Suit,v int) Card {
	
	if v > 13 {
		panic("the value of card can not be bigger than 13")
	}

	return Card{
		suit : s,
		value : v,
	}
}


func (c Card) String() string{
	return fmt.Sprintf("%d of %s %s",c.value,c.suit,suitToUnicode(c.suit))
}


func suitToUnicode(s Suit) string{
	
	sqitch s{
	case Spades:
		return "♠"
	case Harts:
		return "♥"
	case Diamonds:
		return "♦"
	case Clubs:
		return "♣"
	default:
		panic("invalid card suit")
	}

}


func main(){

	

}
