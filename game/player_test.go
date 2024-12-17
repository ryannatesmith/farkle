package game_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/ryannatesmith/farkle/game"
	"testing"
)

func TestPlayer_Roll(t *testing.T) {
	t.Parallel()
	type testCase struct {
		name   string
		player func() *game.Player
		err    bool
	}
	for _, c := range []testCase{
		{
			name: "player accepts",
			player: func() *game.Player {
				player := game.NewPlayer("test", func() uint8 { return 3 }, func(dice int, score uint32) {})
				player.Accept(2, 1000)
				return player
			},
			err: false,
		},
		{
			name: "player rejects",
			player: func() *game.Player {
				player := game.NewPlayer("test", func() uint8 { return 3 }, func(dice int, score uint32) {})
				player.Reject()
				return player
			},
			err: false,
		},
		{
			name: "player does not have a current turn",
			player: func() *game.Player {
				return game.NewPlayer("test", func() uint8 { return 3 }, func(dice int, score uint32) {})
			},
			err: true,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			player := c.player()
			err := player.Roll()
			if !cmp.Equal(c.err, err != nil) {
				t.Error("unexpected error", err)
			}
		})
	}
}

func TestPlayer_Next(t *testing.T) {
	t.Parallel()
	type testCase struct {
		name  string
		play  func(func(int, uint32))
		dice  int
		score uint32
	}
	for _, c := range []testCase{
		{
			name: "player keeps 350",
			play: func(next func(int, uint32)) {
				player := game.NewPlayer(
					"test",
					random([]uint8{1,1,1,5,2,3}),
					next,
					)
				player.Reject()
				player.Roll()
				player.Keep(0,1,2,3)
				player.Bank()
			},
			dice: 2,
			score: 350,
		},
		{
			name: "player farkles",
			play: func(next func(int, uint32)) {
				player := game.NewPlayer(
					"test",
					random([]uint8{2,3,4,6,4,3}),
					next,
					)
				player.Reject()
				player.Roll()
			},
			dice: 0,
			score: 0,
		},
	}{
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			next := func(dice int, score uint32) {
				if c.dice != dice {
					t.Errorf("dice: +want -got\n\t+%d\n\t-%d", c.dice, dice)
				}
				if c.score != score {
					t.Errorf("score: +want -got\n\t+%d\n\t-%d", c.dice, dice)
				}
			}
			c.play(next)
		})
	}
}
