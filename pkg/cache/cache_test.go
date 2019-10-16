package cache

import (
	"math"
	"testing"

	"github.com/andresoro/hemera/pkg/metric"
)

var EPSILON = 0.00000001

func TestCache(t *testing.T) {

	t.Run("test cache reset", func(t *testing.T) {

		c := New()

		for i := 0; i < 100; i++ {
			m := &metric.Metric{
				Name:  "test",
				Value: 1.0,
				Type:  metric.Counter,
				Rate:  1.0,
			}

			err := c.Add(m)
			if err != nil {
				t.Fatal(err)
			}
		}

		c.Clear()

		if len(c.counters) != 0 {
			t.Logf("counter has %d elements", len(c.counters))
			t.Log(c.counters)
			t.Fatal("Counters should have no elements after clear")
		}

	})

	t.Run("test timer stats", func(t *testing.T) {
		c := New()

		// add some timer metrics
		for i := 1; i <= 5; i++ {
			m := &metric.Metric{
				Name:  "test",
				Value: float64(i),
				Type:  metric.Timer,
				Rate:  0,
			}

			c.Add(m)
		}

		// compute timer metric stats
		c.TimerStats()

		// get differences between computed and actual values
		minDiff := math.Abs(c.timerData["test.min"] - 1.0)
		maxDiff := math.Abs(c.timerData["test.max"] - 5.0)
		countDiff := math.Abs(c.timerData["test.count"] - 5.0)
		avgDiff := math.Abs(c.timerData["test.average"] - 3.0)
		stdDiff := math.Abs(c.timerData["test.std_dev"] - math.Sqrt(2.0))
		medDiff := math.Abs(c.timerData["test.median"] - 3.0)
		upperDiff := math.Abs(c.timerData["test.upper_95"] - 5.0)

		// ensure error is within a small bound
		if minDiff > EPSILON {
			t.Error("min not within error bound")
		}

		if maxDiff > EPSILON {
			t.Error("max not within error bound")
		}

		if countDiff > EPSILON {
			t.Error("count not within error bound")
		}

		if avgDiff > EPSILON {
			t.Error("avg not within error bound")
		}

		if stdDiff > EPSILON {
			t.Error("std dev not within error bound")
		}

		if medDiff > EPSILON {
			t.Error("median not within error bound")
		}

		if upperDiff > EPSILON {
			t.Error("upper 95th not within error bound")
		}
	})
}
