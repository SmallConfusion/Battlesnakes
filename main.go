package main

import "log"

func main() {
	RunServer()
}

func move(state GameState) BattlesnakeMoveResponse {
	grid := Grid{}
	grid.SetupFromState(state)

	head := &state.You.Head

	move := DirNull

	check := Coord{}

	dirPrefs := [4]float64{}

	for _, dir := range directions {
		val := grid.Get(check.AddDir(head, dir))

		if val == 0 {
			dirPrefs[dir] += 1
		} else if val == Hazard {
			dirPrefs[dir] += 0.5
		}

		dirPrefs[dir] += float64(distFromHeads(&check, state)) * 0.01
	}

	return BattlesnakeMoveResponse{Move: move.String()}
}

func distFromHeads(check *Coord, state GameState) int {
	min := 99999

	for _, snake := range state.Board.Snakes {
		d := snake.Head.Dist(check)

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
