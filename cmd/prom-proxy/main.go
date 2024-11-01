package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v2"

	"github.com/dsrvlabs/prometheus-proxy/config"
	"github.com/dsrvlabs/prometheus-proxy/types"
)

const (
	defaultHTTPPort   = 9200
	defaultMetricPath = "/metrics"
)

// App represents the application.
type App struct {
	config *types.Config
}

// Prepare prepares prometheus metrics.
func (a *App) Prepare() (*http.ServeMux, error) {
	rpcs := a.config.RPCFetch
	for _, rpc := range rpcs {
		fields := rpc.Fields
		for _, field := range fields {
			prometheus.MustRegister(
				prometheus.NewGauge(
					prometheus.GaugeOpts{
						Name: field.MetricName,
					},
				),
			)
		}
	}

	mux := http.NewServeMux()
	mux.Handle(defaultMetricPath, promhttp.Handler())

	return mux, nil
}

// Run runs the application.
func (a *App) Run(router *http.ServeMux) error {
	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", defaultHTTPPort), router)
}

func main() {
	app := &cli.App{
		Name:  "prom-proxy",
		Usage: "A Prometheus proxy application",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Value:    "config.yaml",
				Usage:    "config file",
				Required: true,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "start",
				Aliases: []string{"s"},
				Usage:   "start proxy server",
				Action: func(cCtx *cli.Context) error {
					// TODO: Get port number
					// TODO: Get server address

					configFile := cCtx.String("config")
					config, err := config.Load(configFile)

					app := NewApp(config)
					router, err := app.Prepare()
					if err != nil {
						return err
					}

					return app.Run(router)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

// NewApp creates a new App instance.
func NewApp(config *types.Config) *App {
	// TODO: Read config
	return &App{
		config: config,
	}
}
