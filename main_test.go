package hatsunesnaku

import (
	"flag"
	"testing"

	"github.com/gin-gonic/gin"
)

var runServer = flag.Bool("server", false, "Runs the server in addition.")

func TestMain(t *testing.T) {
	if !*runServer {
		t.SkipNow()
	}

	s := gin.Default()
	RunServer(s)
	s.Run(":8000")
}
