package cache

import (
	"math"

	"github.com/andresoro/hemera/pkg/metric"
)

// Cache is the metric buffer that holds collections of metrics and purges them at given intervals
// to the backend
type Cache struct {
	// main metrics being tracked
	counters map[string]float64
	gauges   map[string]float64
	timers   map[string][]float64
	sets     map[string]map[int64]struct{}

	// meta data on cache and timers
	timerData  map[string]float64
	seen       int64
	badMetrics int64
}

// New - return a fresh cache
func New() *Cache {
	c := &Cache{
		counters: make(map[string]float64, 0),
		gauges:   make(map[string]float64, 0),
		timers:   make(map[string][]float64, 0),
		sets:     make(map[string]map[int64]struct{}),

		timerData:  make(map[string]float64),
		seen:       int64(0),
		badMetrics: int64(0),
	}

	return c
}

// Add handles a metric and increments or adds to the bucket
func (c *Cache) Add(m *metric.Metric) error {

	c.seen++

	// handle each metric type
	switch m.Type {
	case metric.Counter:
		name := m.Name
		value := m.Value

		if m.Rate > 0 {
			value = value * m.Rate
		}

		// if name counter exists for this name exists
		if _, ok := c.counters[name]; ok {
			c.counters[name] += value
		} else {
			c.counters[name] = value
		}

	case metric.Gauge:
		name := m.Name
		value := m.Value

		c.gauges[name] = value

	case metric.Timer:
		name := m.Name
		value := m.Value

		// check existence for this metric name and append or create a new array
		if _, ok := c.timers[name]; ok {
			c.timers[name] = append(c.timers[name], value)
		} else {
			c.timers[name] = make([]float64, 0)
			c.timers[name] = append(c.timers[name], value)
		}

	case metric.Set:
		name := m.Name
		value := int64(math.Round(m.Value))

		if set, ok := c.sets[name]; ok {
			set[value] = struct{}{}
		}
	}

	return nil
}

// Clear - Set this cache to be a fresh cache with no entries
func (c *Cache) Clear() {
	c = New()
}
