package converter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/dsrvlabs/prometheus-proxy/jsonselector"
	"github.com/dsrvlabs/prometheus-proxy/types"
)

// ResultConvert represents the result of a conversion.
type ResultConvert struct {
	Selector   string
	MetricName string
	Value      string
	Err        error
}

// Converter is the interface for retrieving the contents of a URL.
type Converter interface {
	Fetch(config types.RPCFetchConfig) ([]ResultConvert, error)
}

type converter struct {
	selector jsonselector.Selector
}

func (c *converter) Fetch(config types.RPCFetchConfig) ([]ResultConvert, error) {
	var reqBodyReader io.Reader = nil
	if config.Body != "" {
		reqBodyReader = strings.NewReader(config.Body)
	}

	req, err := http.NewRequest(config.Method, config.URL, reqBodyReader)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	d, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	mapData := make(map[string]interface{})
	err = json.Unmarshal(d, &mapData)
	if err != nil {
		return nil, err
	}

	ret := make([]ResultConvert, len(config.Fields))
	for i, field := range config.Fields {
		value, err := c.selector.Find(mapData, field.Selector)

		ret[i] = ResultConvert{
			Selector:   field.Selector,
			MetricName: field.MetricName,
			Value:      value,
			Err:        err,
		}
	}

	return ret, nil
}

// NewConverter creates a new Converter.
func NewConverter(selector jsonselector.Selector) Converter {
	return &converter{
		selector: selector,
	}
}
