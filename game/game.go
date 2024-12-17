package game

type Game struct {
	players []*Player
	currentPlayer int
}

func (g *Game) Next(dice int, score uint32) {
	// TODO: fill this in to move to the next player
}

func (g *Game) Join(player string) {
	g.players = append(g.players, NewPlayer(player, NewRandom(), g.Next))
}

func (g *Game) Start() {
	
}
