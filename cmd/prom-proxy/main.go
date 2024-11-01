package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v2"

	"github.com/dsrvlabs/prometheus-proxy/config"
	"github.com/dsrvlabs/prometheus-proxy/converter"
	"github.com/dsrvlabs/prometheus-proxy/jsonselector"
	"github.com/dsrvlabs/prometheus-proxy/types"
)

const (
	defaultHTTPPort   = 9200
	defaultMetricPath = "/metrics"
)

// App represents the application.
type App struct {
	config    *types.Config
	converter converter.Converter
	gauges    map[string]prometheus.Gauge
}

// Prepare prepares prometheus metrics.
func (a *App) Prepare() (*http.ServeMux, error) {
	a.gauges = make(map[string]prometheus.Gauge)

	rpcs := a.config.RPCFetch
	for _, rpc := range rpcs {
		fields := rpc.Fields
		for _, field := range fields {
			gauge := prometheus.NewGauge(
				prometheus.GaugeOpts{
					Name: field.MetricName,
				},
			)

			a.gauges[field.MetricName] = gauge

			prometheus.MustRegister(gauge)
		}
	}

	mux := http.NewServeMux()
	mux.Handle(defaultMetricPath, promhttp.Handler())

	return mux, nil
}

// Run runs the application.
func (a *App) Run(router *http.ServeMux) error {
	go func() {
		for {
			a.updateMetrics()
			time.Sleep(5 * time.Second)
		}
	}()

	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", defaultHTTPPort), router)
}

func (a *App) updateMetrics() {
	for _, rpc := range a.config.RPCFetch {
		results, err := a.converter.Fetch(rpc)
		if err != nil {
			continue
		}

		for _, result := range results {
			gauge, ok := a.gauges[result.MetricName]
			if !ok {
				continue
			}

			gauge.Set(result.Value)
		}
	}
}

func main() {
	app := &cli.App{
		Name:  "prom-proxy",
		Usage: "A Prometheus proxy application",
		Commands: []*cli.Command{
			{
				Name:    "start",
				Aliases: []string{"s"},
				Usage:   "start proxy server",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "config",
						Value:    "config.yaml",
						Usage:    "config file",
						Required: true,
					},
				},
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
	selector := jsonselector.NewSelector()
	converter := converter.NewConverter(selector)
	// TODO: Read config
	return &App{
		config:    config,
		converter: converter,
	}
}
