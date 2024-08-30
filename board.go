package main

type Grid struct {
	sizeX int
	sizeY int
	board []int
}

func (g *Grid) SetupFromState(state GameState) {
	g.sizeX = state.Board.Width
	g.sizeY = state.Board.Height

	g.board = make([]int, g.sizeX*g.sizeY)

	for i, snake := range state.Board.Snakes {
		for _, segment := range snake.Body {
			g.Set(&segment, i+1)
		}
	}
}

func (g Grid) IsCoordSafe(pos *Coord) bool {
	if pos.X < 0 || pos.Y < 0 || pos.X >= g.sizeX || pos.Y >= g.sizeY {
		return false
	} else {
		return g.Get(pos) == 0
	}
}

func (g *Grid) Set(pos *Coord, value int) {
	g.board[pos.X+pos.Y*g.sizeX] = value
}

func (g Grid) Get(pos *Coord) int {
	return g.board[pos.X+pos.Y*g.sizeX]
}
