package game_test

import (
  "testing"

  "github.com/google/go-cmp/cmp"
  "github.com/ryannatesmith/farkle/game"
)

func TestRoll_Score(t *testing.T) {
  t.Parallel()
  type testCase struct {
    name string
    roll game.Roll
    want []*game.Scoring
  }
  for _, c := range []testCase{
    {
      name: "three fifty",
      roll: []uint8{1, 1, 1, 5, 4, 3},
      want: []*game.Scoring{
        {Score: 300, Set: []int{0, 1, 2}},
        {Score: 100, Set: []int{0}},
        {Score: 100, Set: []int{1}},
        {Score: 100, Set: []int{2}},
        {Score: 50, Set: []int{3}}},
    },
    {
      name: "six of a kind",
      roll: []uint8{4, 4, 4, 4, 4, 4},
      want: []*game.Scoring{{Score: 3000, Set: []int{0, 1, 2, 3, 4, 5}}},
    },
    {
      name: "five of a kind",
      roll: []uint8{4, 4, 4, 4, 4, 3},
      want: []*game.Scoring{{Score: 2000, Set: []int{0, 1, 2, 3, 4}}},
    },
    {
      name: "five of a kind plus 1",
      roll: []uint8{4, 4, 4, 4, 4, 1},
      want: []*game.Scoring{{Score: 2000, Set: []int{0, 1, 2, 3, 4}}, {Score: 100, Set: []int{5}}},
    },
    {
      name: "five of a kind plus 5",
      roll: []uint8{4, 4, 4, 4, 4, 5},
      want: []*game.Scoring{{Score: 2000, Set: []int{0, 1, 2, 3, 4}}, {Score: 50, Set: []int{5}}},
    },
    {
      name: "four of a kind",
      roll: []uint8{4, 4, 4, 4, 2, 3},
      want: []*game.Scoring{{Score: 1000, Set: []int{0, 1, 2, 3}}},
    },
    {
      name: "two triplets",
      roll: []uint8{3, 2, 3, 2, 3, 2},
      want: []*game.Scoring{
        {Score: 2500, Set: []int{0, 1, 2, 3, 4, 5}},
        {Score: 300, Set: []int{0, 2, 4}},
        {Score: 200, Set: []int{1, 3, 5}},
      },
    },
    {
      name: "three pairs",
      roll: []uint8{3, 3, 4, 4, 5, 5},
      want: []*game.Scoring{
        {Score: 1500, Set: []int{0, 1, 2, 3, 4, 5}},
        {Score: 50, Set: []int{4}},
        {Score: 50, Set: []int{5}},
      },
    },
    {
      name: "four and two",
      roll: []uint8{3, 3, 3, 3, 5, 5},
      want: []*game.Scoring{
        {Score: 1500, Set: []int{0, 1, 2, 3, 4, 5}},
        {Score: 1000, Set: []int{0, 1, 2, 3}},
        {Score: 50, Set: []int{4}},
        {Score: 50, Set: []int{5}},
      },
    },
    {
      name: "three sixes",
      roll: []uint8{6, 6, 6, 4, 3, 2},
      want: []*game.Scoring{{Score: 600, Set: []int{0, 1, 2}}},
    },
    {
      name: "three fives",
      roll: []uint8{5, 5, 5, 4, 3, 2},
      want: []*game.Scoring{
        {Score: 500, Set: []int{0, 1, 2}},
        {Score: 50, Set: []int{0}},
        {Score: 50, Set: []int{1}},
        {Score: 50, Set: []int{2}},
      },
    },
    {
      name: "three fours",
      roll: []uint8{6, 4, 4, 4, 3, 2},
      want: []*game.Scoring{{Score: 400, Set: []int{1, 2, 3}}},
    },
    {
      name: "three threes",
      roll: []uint8{4, 4, 3, 3, 3, 2},
      want: []*game.Scoring{{Score: 300, Set: []int{2, 3, 4}}},
    },
    {
      name: "three twos",
      roll: []uint8{2, 2, 4, 4, 3, 2},
      want: []*game.Scoring{{Score: 200, Set: []int{0, 1, 5}}},
    },
    {
      name: "three ones",
      roll: []uint8{1, 1, 1, 4, 3, 2},
      want: []*game.Scoring{
        {Score: 300, Set: []int{0, 1, 2}},
        {Score: 100, Set: []int{0}},
        {Score: 100, Set: []int{1}},
        {Score: 100, Set: []int{2}},
      },
    },
    {
      name: "three fours and a one",
      roll: []uint8{4, 4, 4, 1, 3, 2},
      want: []*game.Scoring{{Score: 400, Set: []int{0, 1, 2}}, {Score: 100, Set: []int{3}}},
    },
    {
      name: "three fives and a one",
      roll: []uint8{5, 5, 5, 1, 3, 2},
      want: []*game.Scoring{
        {Score: 500, Set: []int{0, 1, 2}},
        {Score: 100, Set: []int{3}},
        {Score: 50, Set: []int{0}},
        {Score: 50, Set: []int{1}},
        {Score: 50, Set: []int{2}},
      },
    },
    {
      name: "four fives and a one",
      roll: []uint8{5, 5, 5, 5, 1, 2},
      want: []*game.Scoring{
        {Score: 1000, Set: []int{0, 1, 2, 3}},
        {Score: 100, Set: []int{4}},
        {Score: 50, Set: []int{0}},
        {Score: 50, Set: []int{1}},
        {Score: 50, Set: []int{2}},
        {Score: 50, Set: []int{3}},
      },
    },
    {
      name: "re-roll three, get one and five",
      roll: []uint8{2, 1, 5},
      want: []*game.Scoring{{Score: 100, Set: []int{1}}, {Score: 50, Set: []int{2}}},
    },
    {
      name: "re-roll five, get five threes",
      roll: []uint8{3, 3, 3, 3, 3},
      want: []*game.Scoring{{Score: 2000, Set: []int{0, 1, 2, 3, 4}}},
    },
    {
      name: "re-roll five, get four of a kind, and a one",
      roll: []uint8{3, 3, 3, 3, 1},
      want: []*game.Scoring{{Score: 1000, Set: []int{0, 1, 2, 3}}, {Score: 100, Set: []int{4}}},
    },
    {
      name: "re-roll four, get four twos",
      roll: []uint8{2, 2, 2, 2},
      want: []*game.Scoring{{Score: 1000, Set: []int{0, 1, 2, 3}}},
    },
    {
      name: "straight",
      roll: []uint8{1, 2, 3, 4, 5, 6},
      want: []*game.Scoring{
        {Score: 1500, Set: []int{0, 1, 2, 3, 4, 5}},
        {Score: 100, Set: []int{0}},
        {Score: 50, Set: []int{4}},
      },
    },
  } {
    t.Run(c.name, func(t *testing.T) {
      t.Parallel()
      got := c.roll.Score()
      if diff := cmp.Diff(c.want, got); diff != "" {
        t.Error("+want -got", diff)
      }
    })
  }
}
