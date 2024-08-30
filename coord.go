package main

type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (c *Coord) Add(a, b *Coord) *Coord {
	c.X = a.X + b.X
	c.Y = a.Y + b.Y
	return c
}

func (c *Coord) AddDir(a *Coord, d Direction) *Coord {
	switch d {
	case Up:
		c.Add(a, &Coord{0, 1})
	case Right:
		c.Add(a, &Coord{1, 0})
	case Down:
		c.Add(a, &Coord{0, -1})
	case Left:
		c.Add(a, &Coord{-1, 0})
	}

	return c
}
