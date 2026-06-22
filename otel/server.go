package otelexporter

import (
	"context"
	"time"

	"transmission-otel/config"
	otelmetrics "transmission-otel/otel/metrics"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

var log = logrus.WithFields(logrus.Fields{
	"prefix": "otel",
})

func Start() {
	ctx := context.Background()

	exporter, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpointURL(config.C.OtelEndpoint),
	)
	if err != nil {
		log.WithError(err).Fatal("Failed to create OTEL exporter")
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(exporter,
				sdkmetric.WithInterval(time.Duration(config.C.FetchInterval)*time.Second),
			),
		),
	)

	otel.SetMeterProvider(provider)

	meter := provider.Meter("transmission-otel")
	if err := otelmetrics.Init(meter); err != nil {
		log.WithError(err).Fatal("Failed to initialize OTEL metrics")
	}

	log.WithField("endpoint", config.C.OtelEndpoint).Info("Starting OTEL metrics exporter.")

}
