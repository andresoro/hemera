package cache

import (
	"fmt"
	"math"
	"sort"
	"sync"

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

	c.seen = 0
	c.badMetrics = 0

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		for k := range c.counters {
			delete(c.counters, k)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for k := range c.gauges {
			delete(c.gauges, k)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for k := range c.sets {
			delete(c.sets, k)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for k := range c.timers {
			delete(c.timers, k)
		}
		wg.Done()
	}()

	wg.Wait()
}

// TimerStats will aggregate all the timers and compute individual statistics
func (c *Cache) TimerStats() {
	timerData := make(map[string]float64)
	var sum float64

	for metric, times := range c.timers {

		sort.Float64s(times)

		count := float64(len(times))
		sum = 0
		for _, i := range times {
			sum += i
		}
		average := sum / count
		stdDev := dev(times, average, count)
		median := percentile(times, count, 0.5)
		upper95 := percentile(times, count, 0.95)

		timerData[fmt.Sprintf("%s.min", metric)] = times[0]
		timerData[fmt.Sprintf("%s.max", metric)] = times[len(times)-1]
		timerData[fmt.Sprintf("%s.count", metric)] = count
		timerData[fmt.Sprintf("%s.average", metric)] = average
		timerData[fmt.Sprintf("%s.std_dev", metric)] = stdDev
		timerData[fmt.Sprintf("%s.median", metric)] = median
		timerData[fmt.Sprintf("%s.upper_95", metric)] = upper95
	}

	c.timerData = timerData
}

// standard deviation
func dev(times []float64, avg, count float64) float64 {
	var sd float64
	for _, time := range times {
		sd += math.Pow(time-avg, 2)
	}

	return math.Sqrt(sd / count)
}

// nth percentile
func percentile(times []float64, count, percent float64) float64 {
	index := int64(count * percent)

	// if even number of values, return average of those values
	if len(times)%2 == 0 {
		return (times[index-1] + times[index]) / 2
	}

	return times[index]
}
