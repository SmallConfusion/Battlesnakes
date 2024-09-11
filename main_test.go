package hatsunesnaku

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(t *testing.T) {
	s := gin.Default()
	RunServer(s)
	s.Run(":8000")
}
