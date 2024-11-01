package config

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dsrvlabs/prometheus-proxy/types"
)

func TestConfig_Load(t *testing.T) {
	tests := []struct {
		Name         string
		Filename     string
		ExpectConfig *types.Config
		ExpectErr    error
	}{
		{
			Name:     "Parse Config",
			Filename: "./fixtures/rpc_config.yaml",
			ExpectConfig: &types.Config{
				RPCFetch: []types.RPCFetchConfig{
					{
						Method: "GET",
						URL:    "http://localhost:26657/status",
						Body:   "",
						Fields: []types.Field{
							{
								Selector:   ".result.sync_info.latest_block_height",
								MetricName: "block_height",
							},
							{
								Selector:   ".result.sync_info.catch_up",
								MetricName: "catch_up",
							},
						},
					},
					{
						Method: "POST",
						URL:    "http://localhost:8545",
						Body:   `{"jsonrpc":"2.0","id":83,"result":false}`,
						Fields: []types.Field{
							{
								Selector:   ".result",
								MetricName: "is_syncing",
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			config, err := Load(tt.Filename)

			assert.Equal(t, tt.ExpectErr, err)

			for i, rpc := range config.RPCFetch {
				assert.Equal(t, tt.ExpectConfig.RPCFetch[i].Method, rpc.Method)
				assert.Equal(t, tt.ExpectConfig.RPCFetch[i].URL, rpc.URL)
				assert.Equal(t, tt.ExpectConfig.RPCFetch[i].Body, rpc.Body)

				for j, field := range rpc.Fields {
					assert.Equal(t, tt.ExpectConfig.RPCFetch[i].Fields[j].Selector, field.Selector)
					assert.Equal(t, tt.ExpectConfig.RPCFetch[i].Fields[j].MetricName, field.MetricName)
				}
			}

		})
	}
}
