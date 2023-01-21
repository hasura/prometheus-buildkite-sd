package buildkite

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const graphqlEndpoint = "https://graphql.buildkite.com/v1"

type graphqlRequestBody struct {
	Query     string            `json:"query"`
	Variables map[string]string `json:"variables,omitempty"`
}

// TODO use buildkite Agent instead
type Agent struct {
}

func GetActiveAgents(apiToken, org string) (*Agent, error) {
	gqlQuery := `
	query allAgents{
		organization(slug: "hasura") {
			agents(first: 100) {
				count
				pageInfo {
					startCursor
					endCursor
					hasPreviousPage
					hasNextPage
				}
				edges {
					node {
						id
						createdAt
						connectedAt
						isRunningJob
						connectionState
						ipAddress
						hostname
						metaData
						job {
							... on JobTypeCommand {
								id
								label
								build {
									id
									pipeline {
										slug
									}
									branch
								}
							}
						}
					}
				}
			}
		}
	}	
	`
	reqBody := graphqlRequestBody{
		Query:     gqlQuery,
		Variables: map[string]string{},
	}
	jsonReqBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	log.Printf("jsonReqBody %s\n", string(jsonReqBody))
	req, err := http.NewRequest("POST", graphqlEndpoint, bytes.NewReader(jsonReqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	fmt.Printf("body: %s", string(body))

	return nil, nil
}
