package game

import "sort"

type Roll []uint8

func (r Roll) Score() []*Scoring {
  values := values(r)
  ret := make([]*Scoring, 0)
  for _, s := range []Scorer{
    SixOfAKind(),
    Straight(),
    TwoTriplets(),
    ThreeDoubles(),
    FiveOfAKind(),
    FourOfAKind(),
    ThreeOfAKind(),
    Ones(),
    Fives(),
  } {
    if scoring := s(values); scoring != nil {
      ret = append(ret, scoring...)
    }
    sort.Slice(ret, func(i, j int) bool {
      return ret[i].Score > ret[j].Score
    })
  }
  return ret
}
