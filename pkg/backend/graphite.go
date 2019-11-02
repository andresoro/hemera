package backend

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/andresoro/hemera/pkg/cache"
)

const (
	Prefix  = "hemera"
	Counter = "counters"
	Gauge   = "gauges"
	Set     = "sets"
	Timer   = "timers"
)

// Graphite implementation of backend interface
type Graphite struct {
	Addr   string
	Logger *log.Logger
}

// Purge implements backend interface
//
// FIXME: Function 'Purge' has too many statements (48 > 40) (funlen)
func (g *Graphite) Purge(c *cache.Cache) (err error) {
	if c.Seen == 0 {
		return nil
	}

	// concatenated buffer to hold all metrics
	var buffer strings.Builder

	now := time.Now().Unix()

	for name, value := range c.Counters {
		fullName := fmt.Sprintf("%s.%s.%s", Prefix, Counter, name)
		metric := metricString(fullName, value, now)
		_, err := buffer.WriteString(metric)
		if err != nil {
			return err
		}
	}

	for name, value := range c.Gauges {
		fullName := fmt.Sprintf("%s.%s.%s", Prefix, Gauge, name)
		metric := metricString(fullName, value, now)
		_, err := buffer.WriteString(metric)
		if err != nil {
			return err
		}
	}

	for name, value := range c.Sets {
		fullName := fmt.Sprintf("%s.%s.%s", Prefix, Counter, name)
		key := fullName + ".count"
		n := float64(len(value))

		metric := metricString(key, n, now)

		_, err := buffer.WriteString(metric)
		if err != nil {
			return err
		}
	}

	// need to call TimerStats to aggregate timer metrics
	stats := c.TimerStats()
	for name, value := range stats {
		fullName := fmt.Sprintf("%s.%s.%s", Prefix, Timer, name)
		metric := metricString(fullName, value, now)
		_, err := buffer.WriteString(metric)
		if err != nil {
			return err
		}
	}

	seen := fmt.Sprintf("%s.seen", Prefix)
	seenMetric := metricString(seen, float64(c.Seen), now)

	_, err = buffer.WriteString(seenMetric)
	if err != nil {
		return err
	}

	// dial graphite server
	conn, err := net.Dial("tcp", g.Addr)
	if err != nil {
		return err
	}
	defer func() {
		if errClose := conn.Close(); err == nil {
			err = errClose
		}
	}()

	// write string buffer
	if g.Logger != nil {
		g.Logger.Printf("pushing %s", buffer.String())
	}
	_, err = conn.Write([]byte(buffer.String()))
	if err != nil {
		return err
	}

	return nil
}

// convert metric to a string that graphite can understand
func metricString(name string, value float64, date int64) string {
	return fmt.Sprintf("%s %f %d \n", name, value, date)
}
