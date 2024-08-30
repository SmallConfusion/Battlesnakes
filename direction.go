package main

import "log"

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

var directions = [4]Direction{Up, Right, Down, Left}

func (d Direction) String() string {
	switch d {
	case Up:
		return "up"
	case Right:
		return "right"
	case Down:
		return "down"
	case Left:
		return "left"
	}

	log.Println("Error with direction, this SHOULD NOT HAPPEN")
	return "up"
}
