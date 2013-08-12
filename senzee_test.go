package senzee

import (
	"github.com/martinolsen/cactuskev-go"
	"testing"
)

func TestSenzee(t *testing.T) {
	cactuskev.AllFive(Eval)
}

func BenchmarkSenzee(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cactuskev.AllFive(Eval)
	}
}
