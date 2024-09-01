package main

import (
	"fmt"
	"math"
)

type Grid struct {
	sizeX   int
	sizeY   int
	you     int
	snakes  []Battlesnake
	hazards []Coord
	food    []Coord
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

	g.snakes = state.Board.Snakes
	g.hazards = state.Board.Hazards
	g.food = state.Board.Food

	for i, snake := range g.snakes {
		if snake.ID == state.You.ID {
			g.you = i
		}
	}
}

func (g Grid) Print() {
	board := make([]string, g.sizeX*g.sizeY)

	for _, hazard := range g.hazards {
		board[hazard.X+hazard.Y*g.sizeX] = "/"
	}

	for i, snake := range g.snakes {
		for _, segment := range snake.Body {
			board[segment.X+segment.Y*g.sizeX] = fmt.Sprint(i + 1)
		}
	}

	for _, food := range g.food {
		board[food.X+food.Y*g.sizeX] = "@"
	}

	for y := 0; y < g.sizeY; y++ {
		for x := 0; x < g.sizeX; x++ {
			if board[x+y*g.sizeX] == "" {
				print(". ")
			} else {
				print(board[x+y*g.sizeX] + " ")
			}
		}
		println()
	}
}

func (g Grid) IsCoordSafe(pos *Coord) bool {
	if pos.X < 0 || pos.Y < 0 || pos.X >= g.sizeX || pos.Y >= g.sizeY {
		return false
	} else {
		return g.Get(pos) == Empty
	}
}

func (g Grid) Get(pos *Coord) int {
	if pos.X < 0 || pos.Y < 0 || pos.X >= g.sizeX || pos.Y >= g.sizeY {
		return OutOfBounds
	} else {
		for i, snake := range g.snakes {
			for _, segment := range snake.Body {
				if segment.Equals(pos) {
					return Player + i
				}
			}
		}

		for _, hazard := range g.hazards {
			if hazard.Equals(pos) {
				return Hazard
			}
		}

		return Empty
	}
}

func (g *Grid) Move() Direction {
	return g.getBestMove(&g.snakes[g.you].Body[0], g.you, g.quickEval)
}

func (g *Grid) getBestMove(pos *Coord, player int, eval func(*Coord, int) float64) Direction {
	evals := make([]float64, 4)
	check := Coord{}

	for _, dir := range directions {
		evals[dir] = eval(check.AddDir(pos, dir), player)
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

	if val >= Player && !g.checkSafeTail(pos) {
		eval -= 990
	}

	if val == OutOfBounds {
		eval -= 999
	}

	if val == Hazard {
		eval -= 0.5
	}

	dist, minSnake := g.headMinDist(pos, player)

	avoidSnake := len(g.snakes[minSnake].Body) >= len(g.snakes[g.you].Body)

	if dist == 1 && avoidSnake {
		eval -= 20

	} else if dist == 1 && !avoidSnake {
		eval += 1

	} else if dist != math.Inf(1) {
		eval += dist * 0.001
	}

	foodMultiplier := 0.001
	if g.snakes[player].Health < 10 {
		foodMultiplier = 0.02
	}

	foodDist := g.foodMinDist(pos)
	if foodDist != math.Inf(1) {
		if foodDist == 0 {
			eval += 0.2
		} else {
			eval += -foodDist * foodMultiplier
		}
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
		if g.Get(check) != Empty {
			return dist
		}

		dist += 1

		check.AddDir(check, dir)
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

func (g Grid) checkSafeTail(pos *Coord) bool {
	for _, snake := range g.snakes {
		if snake.Body[snake.Length-1].Equals(pos) && !snake.Body[snake.Length-2].Equals(pos) {
			return true
		}
	}

	return false
}

func (g *Grid) simulate(selfDir Direction, self int, otherEval func(*Coord, int) float64) func() {
	undoList := []func(){}

	unsetSnakeBody := func(i int, b []Coord, v Coord) {
		undoList = append(undoList, func() {
			b[i] = v
		})
	}

	for i, snake := range g.snakes {
		pos := &g.snakes[i].Body[0]

		var dir Direction

		if i == self {
			dir = selfDir
		} else {
			dir = g.getBestMove(pos, i, otherEval)
		}

		for i := 1; i < snake.Length; i++ {
			unsetSnakeBody(i, snake.Body, snake.Body[i])
			snake.Body[i] = snake.Body[i-1]
		}

		unsetSnakeBody(i, snake.Body, snake.Body[i])
		snake.Body[i].AddDir(&snake.Body[i], dir)
	}

	return func() {
		for _, f := range undoList {
			f()
		}
	}
}
