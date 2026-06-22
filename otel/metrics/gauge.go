package otelmetrics

import (
	"context"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type Gauge struct {
	name        string
	description string
	value       int64
}

type instrument interface {
	register(meter metric.Meter) error
}

var registry []instrument

func NewGauge(name, description string) *Gauge {
	g := &Gauge{name: name, description: description}
	registry = append(registry, g)
	return g
}

// Set updates the value reported by the gauge.
func (g *Gauge) Set(v int64) {
	atomic.StoreInt64(&g.value, v)
}

// SetBool updates the gauge to 1 (true) or 0 (false).
func (g *Gauge) SetBool(v bool) {
	if v {
		g.Set(1)
	} else {
		g.Set(0)
	}
}

func (g *Gauge) register(meter metric.Meter) error {
	_, err := meter.Int64ObservableGauge(g.name,
		metric.WithDescription(g.description),
		metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
			o.Observe(atomic.LoadInt64(&g.value))
			return nil
		}),
	)
	return err
}

type LabeledGauge struct {
	name        string
	description string
	attrKey     string

	mu     sync.Mutex
	values map[string]int64
}

func NewLabeledGauge(name, description, attrKey string) *LabeledGauge {
	g := &LabeledGauge{name: name, description: description, attrKey: attrKey, values: map[string]int64{}}
	registry = append(registry, g)
	return g
}

func (g *LabeledGauge) Set(label string, v int64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.values[label] = v
}

func (g *LabeledGauge) register(meter metric.Meter) error {
	_, err := meter.Int64ObservableGauge(g.name,
		metric.WithDescription(g.description),
		metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
			g.mu.Lock()
			defer g.mu.Unlock()
			for label, v := range g.values {
				o.Observe(v, metric.WithAttributes(attribute.String(g.attrKey, label)))
			}
			return nil
		}),
	)
	return err
}

func Init(meter metric.Meter) error {
	for _, g := range registry {
		if err := g.register(meter); err != nil {
			return err
		}
	}
	return nil
}
