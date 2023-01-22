package main

import (
	"log"
	"net/http"

	"github.com/buildkite/go-buildkite/v3/buildkite"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	BuildkiteToken string `envconfig:"BUILDKITE_TOKEN" required:"true"`
	BuildkiteOrg   string `envconfig:"BUILDKITE_ORG" required:"true"`
}

func main() {
	log.Println("Starting prometheus-buildkite-sd")

	var appConfig AppConfig
	envconfig.MustProcess("", &appConfig)

	r := gin.Default()
	r.GET("/", buildkiteServiceDiscoveryHandler(appConfig))
	r.Run()
}

type Target string

type TargetResult struct {
	Targets []Target `json:"targets"`
	Labels  map[string]string
}

func buildkiteServiceDiscoveryHandler(appConfig AppConfig) func(c *gin.Context) {
	// TODO: while testing this out,
	// a single host machine could host multiple agents
	// so test that case also
	return func(c *gin.Context) {
		log.Println("start of root request")

		config, err := buildkite.NewTokenConfig(appConfig.BuildkiteToken, false)
		if err != nil {
			c.Error(err)
		}
		buildkiteClient := buildkite.NewClient(config.Client())
		var allAgents []buildkite.Agent

		for currentPage, lastPage := 1, 0; (currentPage == 1) || (currentPage <= lastPage); {
			agents, resp, err := buildkiteClient.Agents.List(appConfig.BuildkiteOrg, &buildkite.AgentListOptions{
				ListOptions: buildkite.ListOptions{
					Page:    currentPage,
					PerPage: 4,
				},
			})
			if err != nil {
				c.Error(err)
			}
			lastPage = resp.LastPage
			currentPage++
			allAgents = append(allAgents, agents...)
		}

		type PromSDEntry struct {
			Targets []string          `json:"targets"`
			Labels  map[string]string `json:"labels"`
		}
		var sdEntries []PromSDEntry

		for _, agent := range allAgents {
			entry := PromSDEntry{
				Targets: []string{
					*agent.IPAddress,
				},
				Labels: map[string]string{
					"__meta__host_name":       *agent.Hostname,
					"__meta__connected_state": *agent.ConnectedState, // TODO: see if you could filter by connected state in buildkite api itself
				},
			}

			// TODO: agent.Metadata is an array of string. It could contain two queues with different vaules
			// at that point, the agent is supposed to be serving both of the queues at the same time
			// so have separate SD entries for them
			sdEntries = append(sdEntries, entry)
		}

		c.JSON(http.StatusOK, sdEntries)

		log.Println("end of root request")
	}
}
