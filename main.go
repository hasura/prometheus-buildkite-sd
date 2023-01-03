package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting prometheus-buildkite-sd")

	buildkiteToken := os.Getenv("BUILDKITE_TOKEN")
	if len(buildkiteToken) == 0 {
		log.Fatalln("BUILDKITE_TOKEN environment variable is empty. Please set it to a non-empty value.")
	}

	r := gin.Default()
	r.GET("/", buildkiteServiceDiscoveryHandler(buildkiteToken))
	r.Run()
}

type Target string

type TargetResult struct {
	Targets []Target `json:"targets"`
	Labels  map[string]string
}

func buildkiteServiceDiscoveryHandler(buildkiteToken string) func(c *gin.Context) {
	// TODO: while testing this out,
	// a single host machine could host multiple agents
	// so test that case also
	return func(c *gin.Context) {
		log.Println("start of root request")
		c.JSON(http.StatusOK, gin.H{
			"message": "root",
		})
		log.Println("end of root request")
	}
}
