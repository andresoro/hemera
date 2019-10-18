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
	PREFIX  = "hemera"
	COUNTER = "counters"
	GAUGE   = "gauges"
	SET     = "sets"
	TIMER   = "timers"
)

// Graphite implementation of backend interface
type Graphite struct {
	Addr string
}

// Purge implements backend interface
func (g *Graphite) Purge(c *cache.Cache) error {
	now := time.Now().Unix()

	// concatenated buffer to hold all metrics
	var buffer strings.Builder

	for name, value := range c.Counters {
		fullName := fmt.Sprintf("%s.%s.%s", PREFIX, COUNTER, name)
		metric := metricString(fullName, value, now)
		_, err := buffer.WriteString(metric)
		if err != nil {
			return err
		}
	}

	for name, value := range c.Gauges {
		fullName := fmt.Sprintf("%s.%s.%s", PREFIX, GAUGE, name)
		metric := metricString(fullName, value, now)
		_, err := buffer.WriteString(metric)
		if err != nil {
			return err
		}
	}

	for name, value := range c.Sets {
		fullName := fmt.Sprintf("%s.%s.%s", PREFIX, COUNTER, name)
		key := fullName + ".count"
		n := float64(len(value))

		metric := metricString(key, n, now)

		_, err := buffer.WriteString(metric)
		if err != nil {
			return err
		}
	}

	for name, value := range c.TimerData {
		fullName := fmt.Sprintf("%s.%s.%s", PREFIX, TIMER, name)
		metric := metricString(fullName, value, now)
		_, err := buffer.WriteString(metric)
		if err != nil {
			return err
		}
	}

	seen := fmt.Sprintf("%s.seen", PREFIX)
	seenMetric := metricString(seen, float64(c.Seen), now)

	_, err := buffer.WriteString(seenMetric)
	if err != nil {
		return err
	}

	// dial graphite server
	conn, err := net.Dial("tcp", g.Addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// write string buffer
	log.Printf("pushing %s", buffer.String())
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
