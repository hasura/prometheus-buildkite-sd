package main

import (
	"log"
	"os"
)

func main() {
	log.Println("Starting prometheus-buildkite-sd")

	buildkiteToken := os.Getenv("BUILDKITE_TOKEN")
	if len(buildkiteToken) == 0 {
		log.Fatalln("BUILDKITE_TOKEN environment variable is empty. Please set it to a non-empty value.")
	}

}
