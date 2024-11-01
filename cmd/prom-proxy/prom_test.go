package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dsrvlabs/prometheus-proxy/types"
)

func Test_Prom(t *testing.T) {
	config := &types.Config{
		RPCFetch: []types.RPCFetchConfig{
			{
				Method: "GET",
				URL:    "http://localhost:26657/status",
				Fields: []types.Field{
					{
						Selector:   ".result.sync_info.latest_block_height",
						MetricName: "latest_block_height",
					},
					{
						Selector:   ".result.validator_info.voting_power",
						MetricName: "voting_power",
					},
				},
			},
		},
	}

	app := NewApp(config)
	router, err := app.Prepare()
	assert.Nil(t, err)

	server := httptest.NewServer(router)
	defer server.Close()

	resp, err := http.Get(server.URL + defaultMetricPath)
	assert.Nil(t, err)

	d, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)

	for _, rpc := range config.RPCFetch {
		for _, field := range rpc.Fields {
			assert.Contains(t, string(d), field.MetricName)
		}
	}
}
