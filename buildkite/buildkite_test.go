package buildkite_test

import (
	"os"
	"testing"

	"github.com/hasura/prometheus-buildkite-sd/buildkite"
)

func TestGetActiveAgents(t *testing.T) {
	if _, err := buildkite.GetActiveAgents(os.Getenv("BUILDKITE_TOKEN"), "hasura"); err != nil {
		t.Fail()
	}
}
