# prometheus-buildkite-sd

:construction: Work In Progress :construction: 

Prometheus service discovery for buildkite agents, builds, and jobs.

## Configuration

#### `BUILDKITE_TOKEN`

(Environment Variable. Required.)

Buildkite API token used to fetch data from Buildkite. It could be generated from this [buildkite page](https://buildkite.com/user/api-access-tokens).

The API token will need the following REST API scopes:
- Read Agents (`read_agents`)

#### `BUILDKITE_ORG`

(Environment Variable. Required.)

Buildkite organisation slug. Example: `https://buildkite.com/ORG_SLUG` is where you can visit the buildkite dashboard of your organisation.

#### `TARGET_MODE`

(Environment Variable. Optional)

Target mode denotes the way in which Prometheus target ip addresses could be constructured.

Available options for `TARGET_MODE` are
- `ip-address`
- `host-name`


#### `TARGET_PORTS`

(Environment Variable. Required.)

Target ports is a comma-separated value of ports that will be used to construct `<target_ip_address>:<target_port>` combination of the prometheus targets.

## Develop

```bash
git clone git@github.com:hasura/prometheus-buildkite-sd.git
cd prometheus-buildkite-sd
BUILDKITE_TOKEN="XXXX" BUILDKITE_ORG="YYYY" go run main.go
```
