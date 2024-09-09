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
	return Left
}
