package main

import "log"

func main() {
	RunServer()
}

func move(state GameState) BattlesnakeMoveResponse {
	return BattlesnakeMoveResponse{Move: "left"}
}

func start(state GameState) {
	log.Println("Start game: ", state)
}

func end(state GameState) {
	log.Println("End game: ", state)
}
