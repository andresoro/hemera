package cache

import (
	"math"
	"testing"

	"github.com/andresoro/hemera/pkg/metric"
)

func TestCacheReset(t *testing.T) {
	c := New()

	for i := 0; i < 100; i++ {
		m := &metric.Metric{
			Name:  "test",
			Value: 1.0,
			Type:  metric.Counter,
			Rate:  1.0,
		}

		c.Add(m)
	}

	c.Clear()

	if len(c.Counters) != 0 {
		t.Logf("counter has %d elements", len(c.Counters))
		t.Log(c.Counters)
		t.Fatal("Counters should have no elements after clear")
	}
}

func TestTimerStats(t *testing.T) {
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
	minDiff := math.Abs(c.TimerData["test.min"] - 1.0)
	maxDiff := math.Abs(c.TimerData["test.max"] - 5.0)
	countDiff := math.Abs(c.TimerData["test.count"] - 5.0)
	avgDiff := math.Abs(c.TimerData["test.average"] - 3.0)
	stdDiff := math.Abs(c.TimerData["test.std_dev"] - math.Sqrt(2.0))
	medDiff := math.Abs(c.TimerData["test.median"] - 3.0)
	upperDiff := math.Abs(c.TimerData["test.upper_95"] - 5.0)

	// ensure error is within a small bound
	const epsilon = 0.00000001

	if minDiff > epsilon {
		t.Error("min not within error bound")
	}

	if maxDiff > epsilon {
		t.Error("max not within error bound")
	}

	if countDiff > epsilon {
		t.Error("count not within error bound")
	}

	if avgDiff > epsilon {
		t.Error("avg not within error bound")
	}

	if stdDiff > epsilon {
		t.Error("std dev not within error bound")
	}

	if medDiff > epsilon {
		t.Error("median not within error bound")
	}

	if upperDiff > epsilon {
		t.Error("upper 95th not within error bound")
	}
}
