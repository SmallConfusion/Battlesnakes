package hatsunesnaku

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
	DirNull
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

	return "up"
}
