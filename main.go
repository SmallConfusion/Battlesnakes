package battlesnake

import (
	"log"
	"time"
)

func move(state GameState) BattlesnakeMoveResponse {
	start := time.Now()

	grid := Grid{}
	grid.SetupFromState(state)

	move := grid.Move()

	log.Println("Got move", move.String(), "in time", time.Since(start))

	return BattlesnakeMoveResponse{Move: move.String()}
}

func start(state GameState) {
	log.Println("Start game:", state.Game.ID, state.Game.Map, state.Game.Ruleset)
}

func end(state GameState) {
	log.Println("End game", state.Game.ID)
}
