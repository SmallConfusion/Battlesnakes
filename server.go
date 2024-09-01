package main

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	port     = "8000"
	serverId = "smallconfusion/github/snake"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	response := BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "SmallConfusion",
		Color:      "#ff92f7",
		Head:       "trans-rights-scarf",
		Tail:       "round-bum",
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		log.Println("Error with json encoder, message: ", err)
	}
}

func handleStart(w http.ResponseWriter, r *http.Request) {
	state := GameState{}
	err := json.NewDecoder(r.Body).Decode(&state)

	if err != nil {
		log.Println("Error with json decoder, message: ", err)
	}

	start(state)
}

func HandleEnd(w http.ResponseWriter, r *http.Request) {
	state := GameState{}
	err := json.NewDecoder(r.Body).Decode(&state)

	if err != nil {
		log.Println("Error with json decoder, message: ", err)
	}

	end(state)
}

func HandleMove(w http.ResponseWriter, r *http.Request) {
	state := GameState{}
	err := json.NewDecoder(r.Body).Decode(&state)

	if err != nil {
		log.Println("Error with json decoder, message: ", err)
	}

	response := move(state)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		log.Println("Error with json encoder, message: ", err)
	}
}

func withServerId(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", serverId)
		next(w, r)
	}
}

func RunServer() {
	http.HandleFunc("/", withServerId(handleIndex))
	http.HandleFunc("/start", withServerId(handleStart))
	http.HandleFunc("/move", withServerId(HandleMove))
	http.HandleFunc("/end", withServerId(HandleEnd))

	log.Println("Starting server")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
