package game

import "fmt"

type Player struct {
	name    string
	random  func() uint8
	turns   []*Turn
	current *Turn
	next    func(dice int, score uint32)
}

func (p *Player) Score() uint32 {
	var sum uint32
	for _, turn := range p.turns {
		sum += turn.Result()
	}
	return sum
}

// Accept starts a new turn with the remaining dice and
// score from the previous turn
func (p *Player) Accept(dice int, score uint32) {
	p.current = NewTurn(p.random, WithStart(dice, score))
}

// Reject starts a new turn with six dice and no score
func (p *Player) Reject() {
	p.current = NewTurn(p.random)
}

// Roll rolls the available dice in turn
func (p *Player) Roll() error {
	if p.current == nil {
		return fmt.Errorf("no current turn for player %q", p.name)
	}
	p.current.Roll()
	if p.current.Farkle() {
		defer p.next(p.current.available, p.current.score)
		p.turns = append(p.turns, p.current)
		p.current = nil
	}
	return nil
}

// Keep keeps the given dice
func (p *Player) Keep(dice ...int) error {
	return p.current.Keep(dice...)
}

// Bank concludes the current turn
func (p *Player) Bank() {
	defer p.next(p.current.available, p.current.score)
	p.turns = append(p.turns, p.current)
	p.current = nil
}

func NewPlayer(name string, random Random, next func(dice int, score uint32)) *Player {
	return &Player{name: name, random: random, next: next}
}
