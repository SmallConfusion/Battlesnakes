package main

import (
	"math"
)

type Grid struct {
	sizeX  int
	sizeY  int
	board  []int
	you    int
	snakes []Battlesnake
	food   []Coord
}

const (
	Empty int = iota
	Hazard
	OutOfBounds
	Player
)

func (g *Grid) SetupFromState(state GameState) {
	g.sizeX = state.Board.Width
	g.sizeY = state.Board.Height

	g.board = make([]int, g.sizeX*g.sizeY)

	g.snakes = state.Board.Snakes
	g.food = state.Board.Food

	for i, snake := range g.snakes {
		if snake.ID == state.You.ID {
			g.you = i
		}

		for _, segment := range snake.Body {
			g.Set(&segment, Player+i)
		}
	}

	for _, hazard := range state.Board.Hazards {
		g.Set(&hazard, Hazard)
	}
}

func (g Grid) IsCoordSafe(pos *Coord) bool {
	if pos.X < 0 || pos.Y < 0 || pos.X >= g.sizeX || pos.Y >= g.sizeY {
		return false
	} else {
		return g.Get(pos) == Empty
	}
}

func (g *Grid) Set(pos *Coord, value int) {
	g.board[pos.X+pos.Y*g.sizeX] = value
}

func (g Grid) Get(pos *Coord) int {
	if pos.X < 0 || pos.Y < 0 || pos.X >= g.sizeX || pos.Y >= g.sizeY {
		return OutOfBounds
	} else {
		return g.board[pos.X+pos.Y*g.sizeX]
	}
}

func (g *Grid) Move() Direction {
	evals := make([]float64, 4)
	check := Coord{}

	for i, dir := range directions {
		evals[i] = g.quickEval(check.AddDir(&g.snakes[g.you].Head, dir), g.you)
	}

	max := math.Inf(-1)
	var dir Direction

	for i, eval := range evals {
		if eval > max {
			max = eval
			dir = Direction(i)
		}
	}

	return dir
}

func (g Grid) quickEval(pos *Coord, player int) float64 {
	eval := 0.0
	val := g.Get(pos)

	if val >= Player {
		return -90
	}

	if val == OutOfBounds {
		return -100
	}

	if val == Hazard {
		eval -= 0.5
	}

	dist := g.headMinDist(pos, player)

	if dist == 1 {
		eval -= 0.25
	} else if dist != math.Inf(1) {
		eval += dist * 0.001
	}

	foodMultiplier := 0.001
	if g.snakes[player].Health < 10 {
		foodMultiplier = 0.02
	}

	foodDist := g.foodMinDist(pos)
	if foodDist != math.Inf(1) {
		eval += -foodDist * foodMultiplier
	}

	return eval
}

func (g Grid) headMinDist(pos *Coord, ignorePlayer int) float64 {
	min := math.Inf(1)

	for i, snake := range g.snakes {
		if i == ignorePlayer {
			break
		}

		d := snake.Head.Dist(pos)

		if min > d {
			min = d
		}
	}

	return min
}

func (g Grid) foodMinDist(pos *Coord) float64 {
	min := math.Inf(1)

	for _, food := range g.food {
		d := food.Dist(pos)

		if d < min {
			min = d
		}
	}

	return min
}
