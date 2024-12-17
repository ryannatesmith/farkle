package game

import (
  "fmt"
  "slices"
  "sort"
)

const (
  startDice = 6
)

type Opt func(*Turn)

func WithStart(dice int, score uint32) Opt {
  return func(turn *Turn) {
    turn.available = dice
    turn.score = score
  }
}

type Turn struct {
  available   int
  currentRoll Roll
  rolls       []Roll
  random      func() uint8
  score       uint32
  farkle      bool
}

func (t *Turn) Roll() {
  dice := make([]uint8, t.available)
  for idx := range t.available {
    dice[idx] = t.random()
  }
  t.currentRoll = dice
  if scores := t.currentRoll.Score(); len(scores) == 0 {
    t.available = 0
    t.score = 0
    t.farkle = true
  }
}

func (t *Turn) Farkle() bool {
  return t.farkle
}

func (t *Turn) Keep(i ...int) error {
  kept := len(i)
  if kept > t.available {
    return fmt.Errorf("can only keep %d dice", t.available)
  }
  candidates := make([]*candidate, 0)
  sort.Ints(i)
  scorings := t.currentRoll.Score()
  for _, scoring := range scorings {
    if slices.Equal(i, scoring.Set) {
      roll := make(Roll, len(i))
      for idx, j := range i {
        roll[idx] = t.currentRoll[j]
      }
      t.rolls = append(t.rolls, roll)
      t.score += scoring.Score
      for _, c := range candidates {
        t.rolls = append(t.rolls, c.roll)
        t.score += c.score
      }
      t.available -= kept
      if t.available == 0 {
        t.available = startDice
      }
      return nil
    }
    c, truncated := t.checkSubset(scoring, i...)
    if c != nil {
      candidates = append(candidates, c)
      i = truncated
    }
    if len(i) == 0 {
      t.available -= kept
      if t.available == 0 {
        t.available = startDice
      }
      for _, c := range candidates {
        t.rolls = append(t.rolls, c.roll)
        t.score += c.score
      }
      return nil
    }
  }
  return fmt.Errorf("invalid keep sequence: %v", i)
}

func (t *Turn) Result() uint32 {
  return t.score
}

func (t *Turn) checkSubset(scoring *Scoring, i ...int) (*candidate, []int) {
  if len(scoring.Set) > len(i) {
    return nil, i
  }
  for _, index := range scoring.Set {
    if !slices.Contains(i, index) {
      return nil, i
    }
  }
  roll := make(Roll, len(scoring.Set))
  ret := make([]int, 0, len(i)-len(roll))
  for i, j := range scoring.Set {
    roll[i] = t.currentRoll[j]
  }
  for j := range i {
    if !slices.Contains(scoring.Set, i[j]) {
      ret = append(ret, i[j])
    }
  }
  return &candidate{roll: roll, score: scoring.Score}, ret
}

func NewTurn(random Random, opts ...Opt) *Turn {
  turn := &Turn{available: startDice, random: random}
  for _, opt := range opts {
    opt(turn)
  }
  return turn
}
