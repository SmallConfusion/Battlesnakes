package main

import (
	"fmt"
	"math"
	"sync"
)

type Grid struct {
	sizeX   int
	sizeY   int
	you     int
	snakes  []Battlesnake
	hazards []Coord
	food    []Coord
	state   GameState
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
	g.state = state

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
	evals := make([]float64, 4)

	wg := sync.WaitGroup{}
	wg.Add(4)

	for _, dir := range directions {
		go func() {
			evals[dir] = g.eval(dir, 4)
			wg.Done()
		}()
	}

	wg.Wait()

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

func (g Grid) evalGrid() float64 {
	min := math.Inf(1)

	for _, dir := range directions {
		e := g.posEval((&Coord{}).AddDir(&g.snakes[g.you].Head, dir))

		if e < min {
			min = e
		}
	}

	return min
}

func (g Grid) posEval(pos *Coord) float64 {
	player := g.you
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

	avoidSnake := g.snakes[minSnake].Length >= g.snakes[player].Length

	if dist <= 1 && avoidSnake {
		eval -= 20

	} else if dist < 1 && !avoidSnake {
		eval += 1

	} else if dist != math.Inf(1) {
		eval += dist * 0.001
	}

	foodMultiplier := 0.001
	if g.snakes[player].Health < 10 {
		foodMultiplier = 0.2
	} else if g.snakes[player].Health < 20 {
		foodMultiplier = 0.08
	}

	foodDist := g.foodMinDist(pos)
	if foodDist != math.Inf(1) {
		if foodDist == 0 {
			eval += 0.2 + foodMultiplier
		} else {
			eval += -foodDist * foodMultiplier
		}
	}

	tight := g.clostrophobia(pos)
	if tight < 5 {
		eval -= 1.2 + (5-tight)*0.1
	} else if tight < 15 {
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
			continue
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

func (g *Grid) eval(dir Direction, depth int) float64 {
	if depth == 0 {
		return g.evalGrid()
	} else {
		h := g.snakes[g.you].Head

		for _, seg := range g.snakes[g.you].Body[1:] {
			if h.Equals(&seg) {
				return -100000000
			}
		}

		if h.X < 0 || h.X >= g.sizeX || h.Y < 0 || h.Y >= g.sizeY {
			return -11000000
		}

		return g.evalMoves(0, dir, depth)
	}
}

func (g *Grid) evalMoves(snakeIndex int, evalDir Direction, depth int) float64 {
	min := math.Inf(1)

	for _, dir := range directions {
		e := 0.0

		undo := g.simulate(snakeIndex, dir)

		if snakeIndex == len(g.snakes)-1 {
			e = g.eval(dir, depth-1)
		} else {
			e = g.evalMoves(snakeIndex+1, evalDir, depth)
		}

		if e < min {
			min = e
		}

		undo()
	}

	return min
}

func (g *Grid) simulate(snake int, dir Direction) func() {
	prev_snake_body := make([]Coord, len(g.snakes[snake].Body))
	copy(prev_snake_body, g.snakes[snake].Body)

	prev_snake_head := g.snakes[snake].Head

	for i := len(g.snakes[snake].Body) - 1; i > 0; i-- {
		g.snakes[snake].Body[i] = g.snakes[snake].Body[i-1]
	}

	var head Coord
	head.AddDir(&prev_snake_head, dir)

	g.snakes[snake].Body[0] = head
	g.snakes[snake].Head = head

	return func(prev_snake_body []Coord, prev_snake_head Coord) func() {
		return func() {
			g.snakes[snake].Body = prev_snake_body
			g.snakes[snake].Head = prev_snake_head
		}
	}(prev_snake_body, prev_snake_head)
}
