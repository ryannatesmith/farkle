package game_test

import (
  "testing"

  "github.com/ryannatesmith/farkle/game"
  "github.com/google/go-cmp/cmp"
)

func TestTurn_Play(t *testing.T) {
  t.Parallel()
  type testCase struct {
    name  string
    play  func() *game.Turn
    want  game.Roll
    score uint32
  }
  for _, c := range []testCase{
    {
      name: "one roll keep all",
      play: func() *game.Turn {
        turn := game.NewTurn(func() uint8 {
          return 5
        })
        turn.Roll()
        if err := turn.Keep(0, 1, 2, 3, 4, 5); err != nil {
          t.Fatal(err)
        }
        return turn
      },
      want:  game.Roll{5, 5, 5, 5, 5, 5},
      score: 3000,
    },
    {
      name: "start by continuing roll",
      play: func() *game.Turn {
        turn := game.NewTurn(random([]uint8{5, 4}), game.WithStart(2, 1000))
        turn.Roll()
        if err := turn.Keep(0); err != nil {
          t.Fatal(err)
        }
        return turn
      },
      score: 1050,
    },
    {
      name: "one roll keep four",
      play: func() *game.Turn {
        turn := game.NewTurn(random([]uint8{1, 1, 1, 5, 4, 2}))
        turn.Roll()
        if err := turn.Keep(0, 1, 2, 3); err != nil {
          t.Fatal(err)
        }
        return turn
      },
      want:  game.Roll{1, 1, 1, 5},
      score: 350,
    },
    {
      name: "complicated roll",
      play: func() *game.Turn {
        turn := game.NewTurn(random([]uint8{1, 3, 4, 1, 1, 1, 5, 6, 5}))
        turn.Roll()
        if err := turn.Keep(0, 3, 4, 5); err != nil {
          t.Fatal(err)
        }
        turn.Roll()
        if err := turn.Keep(0); err != nil {
          t.Fatal(err)
        }
        turn.Roll()
        if err := turn.Keep(0); err != nil {
          t.Fatal(err)
        }
        return turn
      },
      want:  []uint8{1, 1, 1, 1, 5, 5},
      score: 1100,
    },
    {
      name: "re-roll all dice and add scores",
      play: func() *game.Turn {
        turn := game.NewTurn(random([]uint8{1, 2, 3, 4, 5, 6, 1, 1, 1, 5, 4, 3}))
        turn.Roll()
        if err := turn.Keep(0, 1, 2, 3, 4, 5); err != nil {
          t.Fatal(err)
        }
        turn.Roll()
        if err := turn.Keep(0, 1, 2, 3); err != nil {
          t.Fatal(err)
        }
        return turn
      },
      want:  []uint8{1, 1, 1, 5},
      score: 1850,
    },
    {
      name: "farkle",
      play: func() *game.Turn {
        turn := game.NewTurn(random([]uint8{3, 3, 3, 5, 4, 2, 4, 3}))
        turn.Roll()
        if err := turn.Keep(0, 1, 2, 3); err != nil {
          t.Fatal(err)
        }
        turn.Roll()
        return turn
      },
      score: 0,
    },
  } {
    t.Run(c.name, func(t *testing.T) {
      t.Parallel()
      turn := c.play()
      gotScore := turn.Result()
      if !cmp.Equal(c.score, gotScore) {
        t.Error(cmp.Diff(c.score, gotScore))
      }
    })
  }
}

func TestTurn_PlayError(t *testing.T) {
  t.Parallel()
  type testCase struct {
    name string
    play func(t *testing.T)
  }
  for _, c := range []testCase{
    {
      name: "too many dice kept",
      play: func(t *testing.T) {
        turn := game.NewTurn(func() uint8 { return 6 })
        turn.Roll()
        if err := turn.Keep(1, 2, 3, 4, 5, 6, 7); err == nil {
          t.Error("should have got error")
        }
        if turn.Result() > 0 {
          t.Error("turn result should be 0")
        }
      },
    },
    {
      name: "index out of range",
      play: func(t *testing.T) {
        turn := game.NewTurn(func() uint8 { return 6 })
        turn.Roll()
        if err := turn.Keep(1, 2, 7); err == nil {
          t.Error("should have got error")
        }
        if turn.Result() > 0 {
          t.Error("turn result should be 0")
        }
      },
    },
    {
      name: "non-scoring di kept",
      play: func(t *testing.T) {
        turn := game.NewTurn(random([]uint8{1, 1, 1, 5, 4, 3}))
        turn.Roll()
        if err := turn.Keep(1, 2, 3, 4, 5); err == nil {
          t.Error("should have got error")
        }
        if turn.Result() > 0 {
          t.Error("turn result should be 0")
        }
      },
    },
    {
      name: "non-scoring di kept with scoring dice",
      play: func(t *testing.T) {
        turn := game.NewTurn(random([]uint8{1, 1, 1, 5, 4, 3}))
        turn.Roll()
        if err := turn.Keep(0, 4); err == nil {
          t.Error("should have got error")
        }
        if turn.Result() > 0 {
          t.Error("turn result should be 0")
        }
      },
    },
  } {
    t.Run(c.name, func(t *testing.T) {
      t.Parallel()
      c.play(t)
    })
  }
}

func random(dice []uint8) func() uint8 {
  roll := -1
  return func() uint8 {
    roll++
    return dice[roll]
  }
}
