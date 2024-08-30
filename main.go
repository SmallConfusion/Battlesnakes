package main

import "log"

func main() {
	RunServer()
}

func move(state GameState) BattlesnakeMoveResponse {
	grid := Grid{}
	grid.SetupFromState(state)

	head := &state.You.Head

	check := Coord{}

	dirPrefs := [4]float64{}

	for _, dir := range directions {
		val := grid.Get(check.AddDir(head, dir))

		if val > 0 {
			dirPrefs[dir] -= 10
		}

		if val == OutOfBounds {
			dirPrefs[dir] -= 100
		}

		if val == Hazard {
			dirPrefs[dir] -= 0.5
		}

		dirPrefs[dir] += float64(distFromHeads(&check, state)) * 0.001

		if state.You.Health > 10 {
			dirPrefs[dir] += -float64(distFromFood(&check, state)) * 0.001
		} else {
			dirPrefs[dir] += -float64(distFromFood(&check, state)) * 0.02
		}
	}

	dir := DirNull
	max := -99999.9

	for i, pref := range dirPrefs {
		if pref > max {
			max = pref
			dir = directions[i]
		}
	}

	return BattlesnakeMoveResponse{Move: dir.String()}
}

func distFromHeads(check *Coord, state GameState) int {
	min := 99999

	for _, snake := range state.Board.Snakes {
		if snake.ID == state.You.ID {
			continue
		}

		d := snake.Head.Dist(check)

		if d < min {
			min = d
		}
	}

	return min
}

func distFromFood(check *Coord, state GameState) int {
	min := 99999

	for _, food := range state.Board.Food {
		d := food.Dist(check)

		if d < min {
			min = d
		}
	}

	return min
}

func start(state GameState) {
	log.Println("Start game:", state.Game.ID, state.Game.Map, state.Game.Ruleset)
}

func end(state GameState) {
	log.Println("End game", state.Game.ID)
}
