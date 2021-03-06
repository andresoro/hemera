package cache

import (
	"encoding/json"
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
	Counters map[string]float64
	Gauges   map[string]float64
	Timers   map[string][]float64
	Sets     map[string]map[int64]struct{}

	// meta data on cache and Timers
	TimerData  map[string]float64
	Seen       int64
	badMetrics int64
	mu         *sync.Mutex
}

// New - return a fresh cache
func New() *Cache {
	c := &Cache{
		Counters: make(map[string]float64),
		Gauges:   make(map[string]float64),
		Timers:   make(map[string][]float64),
		Sets:     make(map[string]map[int64]struct{}),

		Seen:       0,
		badMetrics: 0,
		mu:         &sync.Mutex{},
	}

	return c
}

// Add handles a metric and increments or adds to the bucket
func (c *Cache) Add(m *metric.Metric) {
	c.Seen++

	// handle each metric type
	switch m.Type {
	case metric.Counter:
		name := m.Name
		value := m.Value

		if m.Rate > 0 {
			value *= m.Rate
		}

		if _, ok := c.Counters[name]; !ok {
			c.Counters[name] = 0
		}
		c.Counters[name] += value

	case metric.Gauge:
		name := m.Name
		value := m.Value

		c.Gauges[name] = value

	case metric.Timer:
		name := m.Name
		value := m.Value

		// check existence for this metric name and append or create a new array
		if _, ok := c.Timers[name]; !ok {
			c.Timers[name] = make([]float64, 0)
		}
		c.Timers[name] = append(c.Timers[name], value)

	case metric.Set:
		name := m.Name
		value := int64(math.Round(m.Value))

		if set, ok := c.Sets[name]; ok {
			set[value] = struct{}{}
		}
	}
}

// Clear - Set this cache to be a fresh cache with no entries
func (c *Cache) Clear() {

	c.mu.Lock()

	c.Seen = 0
	c.badMetrics = 0

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		for k := range c.Counters {
			delete(c.Counters, k)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for k := range c.Gauges {
			delete(c.Gauges, k)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for k := range c.Sets {
			delete(c.Sets, k)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for k := range c.Timers {
			delete(c.Timers, k)
		}
		wg.Done()
	}()

	wg.Wait()

	c.mu.Unlock()
}

// TimerStats will aggregate all the Timers and compute individual statistics
func (c *Cache) TimerStats() map[string]float64 {
	timerData := make(map[string]float64)
	var sum float64

	for value, times := range c.Timers {
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

		timerData[fmt.Sprintf("%s.min", value)] = times[0]
		timerData[fmt.Sprintf("%s.max", value)] = times[len(times)-1]
		timerData[fmt.Sprintf("%s.count", value)] = count
		timerData[fmt.Sprintf("%s.average", value)] = average
		timerData[fmt.Sprintf("%s.std_dev", value)] = stdDev
		timerData[fmt.Sprintf("%s.median", value)] = median
		timerData[fmt.Sprintf("%s.upper_95", value)] = upper95
	}

	c.TimerData = timerData

	return timerData
}

// CountersJSON returns json encoded bytes of counters map
func (c *Cache) CountersJSON() ([]byte, error) {
	b, err := json.Marshal(c.Counters)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// TimersJSON returns json encoded bytes of computed timers map
func (c *Cache) TimersJSON() ([]byte, error) {
	c.TimerStats()

	b, err := json.Marshal(c.TimerData)
	if err != nil {
		return nil, err
	}

	return b, err
}

// GaugesJSON returns json encoded bytes of gauges map
func (c *Cache) GaugesJSON() ([]byte, error) {
	b, err := json.Marshal(c.Gauges)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// SetsJSON returns json encoded bytes of gauges map
func (c *Cache) SetsJSON() ([]byte, error) {
	b, err := json.Marshal(c.Sets)
	if err != nil {
		return nil, err
	}

	return b, nil
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
