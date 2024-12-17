package game

var (
  all = []int{0, 1, 2, 3, 4, 5}
)

// Scorer returns a set of indexes that score according to the criteria in Scorer
type Scorer func(map[uint8][]int) []*Scoring

type Scoring struct {
  Score uint32
  Set   []int
}

func SixOfAKind() Scorer {
  return xOfAKind(6, func(_ uint8) uint32 { return 3_000 })
}

func FiveOfAKind() Scorer {
  return xOfAKind(5, func(_ uint8) uint32 { return 2_000 })
}

func FourOfAKind() Scorer {
  return xOfAKind(4, func(_ uint8) uint32 { return 1_000 })
}

func ThreeOfAKind() Scorer {
  return xOfAKind(3, func(n uint8) uint32 {
    if n == 1 {
      return 300
    }
    return uint32(n) * 100
  })
}

func TwoTriplets() Scorer {
  return func(values map[uint8][]int) []*Scoring {
    if len(values) != 2 {
      return nil
    }
    for _, v := range values {
      if len(v) != 3 {
        return nil
      }
    }
    return []*Scoring{{Score: 2_500, Set: all}}
  }
}

func ThreeDoubles() Scorer {
  return func(values map[uint8][]int) []*Scoring {
    switch len(values) {
    case 2:
      for _, v := range values {
        if len(v) != 2 && len(v) != 4 {
          return nil
        }
      }
      return []*Scoring{{Score: 1_500, Set: all}}
    case 3:
      for _, v := range values {
        if len(v) != 2 {
          return nil
        }
      }
      return []*Scoring{{Score: 1_500, Set: all}}
    default:
      return nil
    }
  }
}

func Ones() Scorer {
  return func(values map[uint8][]int) []*Scoring {
    if v, ok := values[1]; ok {
      ret := make([]*Scoring, len(v))
      for i, j := range v {
        ret[i] = &Scoring{Score: 100, Set: []int{j}}
      }
      return ret
    }
    return nil
  }
}

func Fives() Scorer {
  return func(values map[uint8][]int) []*Scoring {
    if v, ok := values[5]; ok {
      ret := make([]*Scoring, len(v))
      for i, j := range v {
        ret[i] = &Scoring{Score: 50, Set: []int{j}}
      }
      return ret
    }
    return nil
  }
}

func Straight() Scorer {
  return func(m map[uint8][]int) []*Scoring {
    if len(m) < 6 {
      return nil
    }
    for _, v := range m {
      if len(v) != 1 {
        return nil
      }
    }
    return []*Scoring{{Score: 1500, Set: all}}
  }
}

func xOfAKind(n int, score func(k uint8) uint32) Scorer {
  return func(values map[uint8][]int) []*Scoring {
    ret := make([]*Scoring, 0)
    for k, v := range values {
      if len(v) == n {
        ret = append(ret, &Scoring{Set: v, Score: score(k)})
      }
    }
    if len(ret) == 0 {
      return nil
    }
    return ret
  }
}

func values(r Roll) map[uint8][]int {
  values := make(map[uint8][]int)
  for idx, c := range r {
    values[c] = append(values[c], idx)
  }
  return values
}
