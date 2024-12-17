package game_test

import (
	"github.com/ryannatesmith/farkle/game"
	"testing"
)

func TestRandom(t *testing.T) {
	t.Parallel()
	random := game.NewRandom()
	for range 1_000 {
		if got := random(); got < 1 || got > 6 {
			t.Errorf("unexpected random %d", got)
		}
	}
}
