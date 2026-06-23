package otelmetrics

import (
	"context"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type Counter struct {
	name        string
	description string
	value       int64
}

func NewCounter(name, description string) *Counter {
	c := &Counter{name: name, description: description}
	registry = append(registry, c)
	return c
}

func (c *Counter) Set(v int64) {
	atomic.StoreInt64(&c.value, v)
}

func (c *Counter) register(meter metric.Meter) error {
	_, err := meter.Int64ObservableCounter(c.name,
		metric.WithDescription(c.description),
		metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
			o.Observe(atomic.LoadInt64(&c.value))
			return nil
		}),
	)
	return err
}

type LabeledCounter struct {
	name        string
	description string
	attrKey     string

	mu     sync.Mutex
	values map[string]int64
}

func NewLabeledCounter(name, description, attrKey string) *LabeledCounter {
	c := &LabeledCounter{name: name, description: description, attrKey: attrKey, values: map[string]int64{}}
	registry = append(registry, c)
	return c
}

func (c *LabeledCounter) Set(label string, v int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.values[label] = v
}

func (c *LabeledCounter) register(meter metric.Meter) error {
	_, err := meter.Int64ObservableCounter(c.name,
		metric.WithDescription(c.description),
		metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
			c.mu.Lock()
			defer c.mu.Unlock()
			for label, v := range c.values {
				o.Observe(v, metric.WithAttributes(attribute.String(c.attrKey, label)))
			}
			return nil
		}),
	)
	return err
}
