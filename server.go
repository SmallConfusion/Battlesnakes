package battlesnake

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	port     = "8000"
	path     = "snake"
	serverId = "smallconfusion/github/snake"
)

func handleIndex(c *gin.Context) {
	r := BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "SmallConfusion",
		Color:      "#ff92f7",
		Head:       "trans-rights-scarf",
		Tail:       "round-bum",
	}

	c.JSON(http.StatusOK, r)
}

func handleStart(c *gin.Context) {
	state := GameState{}
	err := c.BindJSON(&state)

	if err != nil {
		log.Println("Error starting game", err)
		log.Println(io.ReadAll(c.Request.Body))
	}

	start(state)
}

func handleEnd(c *gin.Context) {
	state := GameState{}
	err := c.BindJSON(&state)

	if err != nil {
		log.Println("Error ending game", err)
		log.Println(io.ReadAll(c.Request.Body))
	}

	end(state)
}

func handleMove(c *gin.Context) {
	state := GameState{}
	err := c.BindJSON(&state)

	if err != nil {
		log.Println("Error handling move", err)
		log.Println(io.ReadAll(c.Request.Body))
		return
	}

	r := move(state)

	c.JSON(http.StatusOK, r)
}

func withServerId(c *gin.Context) {
	c.Set("Server", serverId)
}

func RunServer(s *gin.Engine) {
	sh := s.Use(withServerId)

	basePath := "/" + path
	sh.GET(basePath+"/", handleIndex)
	sh.POST(basePath+"/start", handleStart)
	sh.POST(basePath+"/end", handleEnd)
	sh.POST(basePath+"/move", handleMove)
}
