package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/buildkite/go-buildkite/v3/buildkite"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
)

type TargetMode string

var IPAddress TargetMode = "ip-address"
var HostName TargetMode = "host-name"

type AppConfig struct {
	BuildkiteToken string     `split_words:"true" required:"true"`
	BuildkiteOrg   string     `split_words:"true" required:"true"`
	TargetMode     TargetMode `split_words:"true" default:"ip-address"`
	TargetPorts    string     `split_words:"true" required:"true"`
}

func main() {
	log.Println("Starting prometheus-buildkite-sd")

	var appConfig AppConfig
	envconfig.MustProcess("", &appConfig)

	// log useful information on server start
	log.Println("config -> org", appConfig.BuildkiteOrg)
	log.Println("config -> target-mode", appConfig.TargetMode)

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
					PerPage: 100,
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
				Labels: map[string]string{
					"buildkite_agent_host_name": *agent.Hostname,
					// TODO: see if you could filter by connected state in buildkite api itself
				},
			}

			if appConfig.TargetMode == IPAddress {
				entry.Targets = []string{
					*agent.IPAddress,
				}
			} else if appConfig.TargetMode == HostName {
				entry.Targets = []string{
					*agent.Hostname,
				}
			} else {
				c.Error(errors.New("unknown TARGET_MODE, please set that env to ip-address or host-name"))
			}

			for _, port := range strings.Split(appConfig.TargetPorts, ",") {
				for idx := range entry.Targets {
					entry.Targets[idx] = fmt.Sprintf("%s:%s", entry.Targets[idx], port)
				}
			}

			if agent.Job != nil {
				// A job is actively running on the agent
				// So add metadata related to it
				if agent.Job.ID != nil {
					entry.Labels["buildkite_job_id"] = *agent.Job.ID
				}

				if agent.Job.StepKey != nil {
					entry.Labels["buildkite_job_key"] = *agent.Job.StepKey
				}
			}

			// TODO: agent.Metadata is an array of string. It could contain two queues with different vaules
			// at that point, the agent is supposed to be serving both of the queues at the same time
			// so have separate SD entries for them
			sdEntries = append(sdEntries, entry)
		}

		// TODO: handle c.Error

		c.JSON(http.StatusOK, sdEntries)

		log.Println("end of root request")
	}
}
