package main

import (
	"log"
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

	for _, dir := range directions {
		evals[dir] = g.quickEval(check.AddDir(&g.snakes[g.you].Head, dir), g.you)
	}

	max := math.Inf(-1)
	var dir Direction

	for i, eval := range evals {
		if eval > max {
			max = eval
			dir = Direction(i)
		}
	}

	log.Println(evals)

	return dir
}

func (g Grid) quickEval(pos *Coord, player int) float64 {
	eval := 0.0
	val := g.Get(pos)

	if val >= Player {
		eval -= 999999
	}

	if val == OutOfBounds {
		eval -= 9999999
	}

	if val == Hazard {
		eval -= 0.5
	}

	dist, minSnake := g.headMinDist(pos, player)

	avoidSnake := len(g.snakes[minSnake].Body) >= len(g.snakes[g.you].Body)

	if dist == 1 && avoidSnake {
		eval -= 2
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

	tight := g.clostrophobia(pos)
	if tight < 15 {
		eval -= 0.05 * (15 - tight)
	}

	return eval
}

func (g Grid) clostrophobia(pos *Coord) float64 {
	total := 0.0

	for _, dir := range directions {
		total += g.raycast(pos, dir)
	}

	return total
}

func (g Grid) raycast(pos *Coord, dir Direction) float64 {
	check := pos.Copy()
	dist := 0.0

	for {
		if g.Get(&check) != Empty {
			return dist
		}

		dist += 1

		check.AddDir(&check, dir)
	}
}

func (g Grid) headMinDist(pos *Coord, ignorePlayer int) (dist float64, minSnake int) {
	dist = math.Inf(1)

	for i, snake := range g.snakes {
		if i == ignorePlayer {
			break
		}

		d := snake.Head.Dist(pos)

		if dist > d {
			dist = d
			minSnake = i
		}
	}

	return
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
