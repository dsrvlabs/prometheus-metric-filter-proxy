package types

// Config is a struct that holds the configuration for the application.
type Config struct {
	RPCFetch []RPCFetchConfig `yaml:"rpc_fetch,omitempty"`
}

// RPCFetchConfig is a struct that holds conversion configuration for RPC response into Prometheus.
type RPCFetchConfig struct {
	Method string  `yaml:"method,omitempty"`
	URL    string  `yaml:"url,omitempty"`
	Body   string  `yaml:"body,omitempty"`
	Fields []Field `yaml:"fields,omitempty"`
	Labels []Label `yaml:"labels,omitempty"`
}

// Field is a struct that holds the configuration for a field in the response.
type Field struct {
	Selector   string `yaml:"selector,omitempty"`
	MetricName string `yaml:"metric_name,omitempty"`
}

// Label is a struct that holds the configuration for a label in the response.
type Label struct {
	Key   string `yaml:"key,omitempty"`
	Value string `yaml:"value,omitempty"`
}
