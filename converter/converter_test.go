package converter

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	"github.com/dsrvlabs/prometheus-proxy/jsonselector"
	"github.com/dsrvlabs/prometheus-proxy/types"
)

var fixture = `
{
  "jsonrpc": "2.0",
  "id": -1,
  "result": {
    "node_info": {
      "protocol_version": {
        "p2p": "8",
        "block": "11",
        "app": "0"
      },
      "id": "1b38f69801c45bc30fdf15677c72af0296d56954",
      "listen_addr": "tcp://0.0.0.0:26656",
      "network": "archway-1",
      "version": "0.38.12",
      "channels": "40202122233038606100",
      "moniker": "alpha-bravo",
      "other": {
        "tx_index": "on",
        "rpc_address": "tcp://127.0.0.1:26657"
      }
    },
    "sync_info": {
      "latest_block_hash": "28579CACDCFE665CC1181D8888ED12F4E98563FCB966E28DAB7881E538629B10",
      "latest_app_hash": "E0DA1C67C7B2E055DDEA600EBBB90DE0E52D23C5FCD63CD77656C35E845A7AB1",
      "latest_block_height": "7090750",
      "latest_block_time": "2024-10-31T12:11:53.68977979Z",
      "earliest_block_hash": "D6AF49B540F7D22631482159412BC1239525344FD8BBC2EA8014746820BAF695",
      "earliest_app_hash": "49D90121E0964D3FC34CB8A8F2F0092EFA6E8CC034A354962779D0F341B9DB88",
      "earliest_block_height": "5904000",
      "earliest_block_time": "2024-08-12T07:39:10.582595419Z",
      "catching_up": false
    },
    "validator_info": {
      "address": "5E026F83F8DDC51308008A80518530E9C03C7771",
      "pub_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "hDarO/5z+0XjbNfJMUsPf40LZN7v8a1Y+5BG93oAN7U="
      },
      "voting_power": "27293355"
    }
  }
}
`

func Test_Converter(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	selector := jsonselector.NewSelector()
	converter := NewConverter(selector)

	rpcFetchConfig := types.RPCFetchConfig{
		Method: "GET",
		URL:    "http://localhost:26657/metrics",
		Body:   "",
		Fields: []types.Field{
			{
				Selector:   ".result.node_info.network",
				MetricName: "network",
			},
			{
				Selector:   ".result.sync_info.latest_block_height",
				MetricName: "latest_block_height",
			},
		},
	}

	httpmock.RegisterResponder(
		http.MethodGet,
		"http://localhost:26657/metrics",
		httpmock.NewStringResponder(200, fixture))

	ret, err := converter.Fetch(rpcFetchConfig)
	assert.Nil(t, err)

	assert.Equal(t, len(rpcFetchConfig.Fields), len(ret))
	assert.Equal(t, "archway-1", ret[0].Value)
	assert.Equal(t, "7090750", ret[1].Value)
}
