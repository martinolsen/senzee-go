package senzee

import (
	"github.com/martinolsen/cactuskev-go"
	"testing"
)

func BenchmarkSenzee(b *testing.B) {
	table := New()

	for i := 0; i < b.N; i++ {
		cactuskev.AllFive(nil, table.Eval)
	}
}

func TestEval5HandFast(t *testing.T) {
	cactuskev.AllFive(t, func(h cactuskev.Hand) cactuskev.Score {
		return cactuskev.CactusKevScore(eval_5hand_fast(h))
	})
}

func TestEval7Hand(t *testing.T) {
	if testing.Short() {
		t.Skip("not running long tests")
	}

	cactuskev.AllSeven(t, func(h cactuskev.Hand) cactuskev.Score {
		return cactuskev.CactusKevScore(eval_7hand(h))
	})
}

func TestEval7HandFast(t *testing.T) {
	table := New()

	cactuskev.AllSeven(t, func(h cactuskev.Hand) cactuskev.Score {
		return cactuskev.CactusKevScore(table.eval_7hand_fast(h))
	})
}

func TestFive(t *testing.T) {
	if testing.Short() {
		t.Skip("not running long tests")
	}

	t.Logf("Creating table...")
	table := New()

	cactuskev.AllFive(t, table.Eval)
}

func TestSeven(t *testing.T) {
	if testing.Short() {
		t.Skip("not running long tests")
	}

	cactuskev.AllSeven(t, New().Eval)
}

func TestCardToInt(t *testing.T) {
	for s := uint(0); s < 4; s++ {
		for r := 0; r < 13; r++ {
			c := cactuskev.NewCard(cactuskev.Suit(0x1000<<(s)), cactuskev.Rank(r))
			i := CardToInt(c)

			switch {
			// TODO - verify number of set bits
			case i < 1:
				t.Errorf("integer too small for %v, got %d (b%064b)", c, i, i)
			case i > 1<<51:
				t.Errorf("integer too large for %v, got %d (b%064b)", c, i, i)
			default:
				t.Logf(
					"OK %v => b%016b %016b %016b %016b",
					c,
					i>>48,
					(i&0x0000ffff00000000)>>32,
					(i&0x00000000ffff0000)>>16,
					i&0x000000000000ffff,
				)
			}

		}
	}
}

func TestHandToInt(t *testing.T) {
	if !testing.Short() {
		t.Logf("reducing logging during long test to preserve memory")
	}

	var (
		deck = cactuskev.NewDeck()
		h    = cactuskev.NewHand(7)
	)

	nbits := func(n int64) uint {
		var c uint
		for ; n != 0; c++ {
			n &= n - 1
		}

		return c
	}

	for a := 0; a < 46; a++ {
		h.SetCard(0, deck[a])
		for b := a + 1; b < 47; b++ {
			h.SetCard(1, deck[b])
			for c := b + 1; c < 48; c++ {
				h.SetCard(2, deck[c])
				for d := c + 1; d < 49; d++ {
					h.SetCard(3, deck[d])
					for e := d + 1; e < 50; e++ {
						h.SetCard(4, deck[e])
						for f := e + 1; f < 51; f++ {
							h.SetCard(5, deck[f])
							for g := f + 1; g < 52; g++ {
								h.SetCard(6, deck[g])

								i := HandToInt(h)

								switch {
								// TODO - verify number of set bits
								case i < 0x7f:
									t.Errorf("integer too small for %v, got %d (b%064b)", h, i, i)
								case i > 0x7f<<46:
									t.Errorf("integer too large for %v, got %d (b%064b)", h, i, i)
								case nbits(i) != 7:
									t.Errorf("expected 7 bits set, got %d (b%064b)", nbits(i), i)
								default:
									if testing.Short() {
										t.Logf(
											"OK %s => b%016b %016b %016b %016b",
											h,
											i>>48,
											(i&0x0000ffff00000000)>>32,
											(i&0x00000000ffff0000)>>16,
											i&0x000000000000ffff,
										)

									}
								}
							}
						}

						if testing.Short() {
							t.Skip("skipping rest of hands during short test")
						}
					}
				}
			}
		}
	}
}
