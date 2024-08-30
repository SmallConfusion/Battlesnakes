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
	for _, dir := range directions {
		check := grid.Get(check.AddDir(head, dir))

		if check == 0 {
			move = dir
			break
		} else if check == Hazard {
			move = dir
		}
	}

	return BattlesnakeMoveResponse{Move: move.String()}
}

func start(state GameState) {
	log.Println("Start game")
}

func end(state GameState) {
	log.Println("End game")
}
