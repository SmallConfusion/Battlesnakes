package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/tiendc/go-deepcopy"
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

const moveTime = 300 * time.Millisecond

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

		g.snakes[i].Dead = false
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

func (g Grid) isPosDeadly(pos *Coord) bool {
	if pos.X < 0 || pos.Y < 0 || pos.X >= g.sizeX || pos.Y >= g.sizeY {
		return true
	} else {
		for _, snake := range g.snakes {
			if snake.Dead {
				continue
			}

			for j, segment := range snake.Body {
				if j == 0 {
					continue
				}

				if segment.Equals(pos) {
					return true
				}
			}
		}

		return false
	}
}

func (g *Grid) isFoodAt(pos *Coord) bool {
	for _, food := range g.food {
		if pos.Equals(&food) {
			return true
		}
	}

	return false
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

func (g *Grid) GetMove() Direction {
	start := time.Now()

	var dir Direction
	var eval float64

	for depth := 1; depth < 5; depth++ {
		thisEval, thisDir := g.eval(depth)

		if eval != thisEval {
			dir = thisDir
			eval = thisEval
		}

		log.Println("Depth", depth, "Move", dir, "Eval", eval)

		if time.Since(start) > moveTime {
			break
		}
	}

	return dir
}

func (g *Grid) simulate(dirs []Direction) {
	for i, dir := range dirs {
		g.moveSnake(i, dir)
	}

	toBeDead := make([]int, 0)

Snakes:
	for i := 0; i < len(g.snakes); i++ {
		head := g.snakes[i].Head

		if g.snakes[i].Health <= 0 {
			toBeDead = append(toBeDead, i)
			continue Snakes
		}

		if g.isPosDeadly(&head) {
			toBeDead = append(toBeDead, i)
			continue Snakes
		}

		for j := 0; j < len(g.snakes); j++ {
			if j == i {
				continue Snakes
			}

			if head.Equals(&g.snakes[j].Head) && g.snakes[i].Length <= g.snakes[j].Length {
				toBeDead = append(toBeDead, i)
				continue Snakes
			}
		}
	}

	for _, index := range toBeDead {
		g.snakes[index].Dead = true
	}
}

func (g *Grid) moveSnake(snakeIndex int, dir Direction) {
	snake := &g.snakes[snakeIndex]

	if snake.Dead {
		return
	}

	snake.Head.AddDir(&snake.Head, dir)

	for i := len(snake.Body) - 1; i > 0; i -= 1 {
		snake.Body[i] = snake.Body[i-1]
	}

	snake.Body[0] = snake.Head

	snake.Health -= 1

	if g.isFoodAt(&snake.Head) {
		snake.Health = 100
		snake.Body = append(snake.Body, snake.Body[len(snake.Body)-1])
		snake.Length++
	}
}

func (g Grid) eval(depth int) (float64, Direction) {
	eval := g.evalBase()

	if eval < -1000 {
		return -10000, Left
	}

	if depth == 0 {
		return eval, Left
	} else {
		max := math.Inf(-1)
		maxDir := Left

		for _, dir := range directions {
			min := math.Inf(1)

			totalMoves := int(math.Pow(4, float64(len(g.snakes)-1)))

			for i := 0; i < totalMoves; i++ {
				moves := make([]Direction, len(g.snakes))

				for j := 0; j < len(g.snakes); j++ {
					if j == g.you {
						moves[j] = dir
						continue
					}

					index := j

					if j > g.you {
						j -= 1
					}

					base := math.Pow(4, float64(j))
					next := math.Pow(4, float64(j+1))
					dirIndex := int(math.Mod(float64(i), next) / base)

					moves[index] = directions[dirIndex]

				}

				var newGrid Grid
				deepcopy.Copy(&newGrid, &g)

				newGrid.simulate(moves)
				eval, _ = newGrid.eval(depth - 1)

				if eval < min {
					min = eval
				}
			}

			if min >= max {
				maxDir = dir
				max = min
			}

		}

		return max, maxDir
	}
}

func (g Grid) evalBase() float64 {
	if g.snakes[g.you].Dead {
		return -10000
	}

	deadCount := 0
	otherLength := 0

	for i, snake := range g.snakes {
		if snake.Dead {
			deadCount++
		} else if i != g.you {
			otherLength += snake.Length
		}
	}

	return float64(deadCount)*1000 + float64(g.snakes[g.you].Length) - float64(otherLength)*0.1
}
