package senzee

import (
	"math/rand"
	"testing"
)

func TestIndex52c7(t *testing.T) {
	for i := 0; i < 10e4; i++ {
		var n int64
		for j := rand.Intn(7) + 1; j > 0; j-- {
			jn := int64(1 << (rand.Uint32()%13 + ((rand.Uint32() % 4) * 13)))
			if jn&n == jn {
				j++
				continue
			}

			n += jn
		}

		x := index52c7(n)

		if x > TableSize {
			t.Errorf("above table limit: %016x => %d vs %d", n, x, TableSize)
		}
	}
}
