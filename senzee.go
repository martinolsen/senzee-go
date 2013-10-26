// Credit: http://senzee.blogspot.com
package senzee

import (
	"encoding/gob"
	"fmt"
	"github.com/cznic/mathutil"
	"github.com/martinolsen/cactuskev-go"
	"log"
	"os"
)

const TableSize = 133784560

type Table [TableSize]int16

func New() *Table {
	t := new(Table)

	var (
		deck = cactuskev.NewDeck()
		hand = cactuskev.NewHand(7)
	)

	for a := 0; a < 46; a++ {
		hand.SetCard(0, deck[a])
		for b := a + 1; b < 47; b++ {
			hand.SetCard(1, deck[b])
			for c := b + 1; c < 48; c++ {
				hand.SetCard(2, deck[c])
				for d := c + 1; d < 49; d++ {
					hand.SetCard(3, deck[d])
					for e := d + 1; e < 50; e++ {
						hand.SetCard(4, deck[e])
						for f := e + 1; f < 51; f++ {
							hand.SetCard(5, deck[f])
							for g := f + 1; g < 52; g++ {
								hand.SetCard(6, deck[g])

								(*t)[index52c7(HandToInt(hand))] = eval_7hand(hand)
							}
						}
					}
				}
			}
		}
	}

	return t
}

func (t *Table) Store(name string) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}

	if err := gob.NewEncoder(file).Encode(t); err != nil {
		return err
	}

	return nil
}

func Load(name string) (*Table, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	t := new(Table)

	if err := gob.NewDecoder(file).Decode(t); err != nil {
		return nil, err
	}

	return t, nil
}

func (t *Table) Eval(h cactuskev.Hand) cactuskev.Score {
	switch h.Len() {
	case 5:
		return cactuskev.CactusKevScore(eval_5hand_fast(h))
	case 7:
		return cactuskev.CactusKevScore(t.eval_7hand_fast(h))
	default:
		return cactuskev.Eval(h)
	}
}

func (t *Table) eval_7hand_fast(h cactuskev.Hand) int16 {
	return (*t)[index52c7(HandToInt(h))]
}

func NewHand(n int) cactuskev.Hand {
	switch n {
	case 5:
		return NewHand5()
	default:
		return cactuskev.NewHand(n)
	}
}

type hand5 uint64

func NewHand5() *hand5 { h := hand5(0); return &h }

func (h hand5) Prime() int {
	var prime = 1

	for i := 0; i < 5; i++ {
		prime *= h.Card(i).Prime()
	}

	return prime
}

func (h hand5) Bit() (bit int) {
	for i := 0; i < 5; i++ {
		bit |= h.Card(i).Bit()
	}

	return bit
}

func (h hand5) Len() int {
	var (
		c int
		n = int64(h)
	)

	for ; n != 0; c++ {
		n &= n - 1
	}

	return c
}

func (h *hand5) SetCard(_ int, c cactuskev.Card) {
	// TODO - do we need some verification of `index`?

	*h |= hand5(CardToInt(c))
}

func (h hand5) Cards() []cactuskev.Card { panic("not implemented") }

func (h hand5) Card(n int) cactuskev.Card {
	if n < 0 || n > 5 {
		panic("card index out of range")
	}

	if int64(h) == 0 {
		panic("empty hand")
	}

	//fmt.Printf("%064b\n%064b\n", int64(h), 1<<uint(mathutil.BitLenUint64(uint64(h))-1))

	for i := hand5(mathutil.BitLenUint64(uint64(h)) - 1); ; i-- {
		/*fmt.Printf("n: % 2d, i: %064b (%d)\n", n, 1<<(i), i)
		fmt.Printf("     , h: %064b\n", h)
		fmt.Printf("     , ?: %064b\n", h&hand5(1<<i))*/

		if n < 1 {
			return IntToCard(uint64(1 << i))
		} else if h&hand5(1<<i) != 0 {
			n--
		}
	}

	panic("card not found")
}

func (h hand5) String() string { return fmt.Sprintf("0x%016x b%064b", int64(h), int64(h)) }

func eval_7hand(h cactuskev.Hand) int16 {
	best := int16(9999)

	for i, sh := 0, hand5(0); i < 21; i++ {
		for j := 0; j < 5; j++ {
			sh.SetCard(j, h.Card(cactuskev.Perm7[i][j]))
		}

		if q := eval_5hand_fast(&sh); q < best {
			best = q
		}
	}

	return best
}

func eval_5hand_fast(h cactuskev.Hand) int16 {
	var q = h.Bit()

	// Flushes and Straight Flushes
	if h.Card(0).Suit()&h.Card(1).Suit()&h.Card(2).Suit()&h.Card(3).Suit()&h.Card(4).Suit() != 0 {
		return cactuskev.Flushes[q]
	}

	// Straights and High Cards
	if s := cactuskev.Unique5[q]; s != 0 {
		return s
	}

	return int16(hash_values[find_fast(uint32(h.Prime()))])
}

func IntToCard(n uint64) cactuskev.Card {
	var i = mathutil.BitLenUint64(n) - 1

	var card = cactuskev.NewCard(
		cactuskev.Suit(0x1000<<uint(i/13)),
		cactuskev.Rank(i%13),
	)

	return card
}

// Convert to a one-bit representation
func CardToInt(c cactuskev.Card) int64 {
	switch c.Suit() {
	case cactuskev.Club:
		return 1 << (uint(c.Rank()) + (0 * 13))
	case cactuskev.Diamond:
		return 1 << (uint(c.Rank()) + (1 * 13))
	case cactuskev.Heart:
		return 1 << (uint(c.Rank()) + (2 * 13))
	case cactuskev.Spade:
		return 1 << (uint(c.Rank()) + (3 * 13))
	default:
		log.Panicf("unknown Suit: 0x%04x", c.Suit())
		return 0
	}
}

func HandToInt(h cactuskev.Hand) int64 {
	var x int64

	for i := 0; i < h.Len(); i++ {
		x |= CardToInt(h.Card(i))
	}

	return x
}

func find_fast(u uint32) uint {
	u += 0xe91aaa35
	u ^= u >> 16
	u += u << 8
	u ^= u >> 4

	var (
		a = (u + (u << 2)) >> 19
		b = (u >> 8) & 0x1ff
	)

	return uint(a ^ uint32(hash_adjust[b]))
}
