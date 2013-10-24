package senzee

import (
	"github.com/martinolsen/cactuskev-go"
	"testing"
)

func TestHand5(t *testing.T) {
	h := NewHand5()

	if l := h.Len(); l != 0 {
		t.Errorf("invalid length. expected %d, got: %d", 5, l)
	}

	t.Logf("h: %s", h)

	h.SetCard(0, cactuskev.NewCard(cactuskev.Spade, cactuskev.Ace))

	t.Logf("h: %s", h)

	if c := h.Card(0); c.Rank() != cactuskev.Ace || c.Suit() != cactuskev.Spade {
		t.Errorf("unexpected card: %s", c)
	} else {
		t.Logf("h[0]: %s", c)
	}

	t.Logf("Card(0): %s", h.Card(0))
	t.Logf("Card(4): %s", h.Card(4))

	if h.Len() != 1 {
		t.Errorf("unexpected length")
	}
}
