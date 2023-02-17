package deck

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
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

type Card struct {
	value int
	suit  Suit
}

func NewCard(s Suit, v int) Card {

	if v > 13 {
		panic("the value of card can not be bigger than 13")
	}

	return Card{
		suit:  s,
		value: v,
	}
}

func (c Card) String() string {

	value := strconv.Itoa(c.value)
	if value == "1" {
		value = "ACE"
	}

	if value == "11" {
		value = "Jack"
	}

	if value == "12" {
		value = "Queen"
	}

	if value == "13" {
		value = "King"
	}

	return fmt.Sprintf("%s of %s %s", value, c.suit, suitToUnicode(c.suit))
}

func Shuffle(d Deck) Deck {

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < len(d); i++ {

		r := rand.Intn(i + 1)

		if r != i {
			d[i], d[r] = d[r], d[i]
		}
	}

	return d
}

func suitToUnicode(s Suit) string {

	switch s {
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

type Deck [52]Card

func New() Deck {

	nSuits := 4
	nCards := 13

	d := [52]Card{}

	x := 0
	for i := 0; i < nSuits; i++ {

		for j := 0; j < nCards; j++ {

			d[x] = NewCard(Suit(i), j+1)
			x++
		}
	}

	return Shuffle(d)
}
