package main

import "log"

func main() {
	RunServer()
}

func move(state GameState) BattlesnakeMoveResponse {
	log.Println("Move requested")

	grid := Grid{}
	grid.SetupFromState(state)

	head := &state.You.Head

	move := Left

	check := Coord{}
	for _, dir := range directions {
		if grid.IsCoordSafe(check.AddDir(head, dir)) {
			move = dir
			break
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
