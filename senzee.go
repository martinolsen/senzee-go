// Credit: http://senzee.blogspot.com
package senzee

import (
	"github.com/martinolsen/cactuskev-go"
)

func Eval(h cactuskev.Hand) cactuskev.Score {
	switch h.Len() {
	case 5:
		return cactuskev.Score(eval_5hand_fast(h))
	default:
		return cactuskev.Eval(h)
	}
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
