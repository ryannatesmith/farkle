package game

import "math/rand/v2"

type Random func() uint8

func NewRandom() Random {
	return func() uint8 {
		return uint8(rand.Uint32N(6) + 1)
	}
}
