package cache

import (
	"testing"

	"github.com/andresoro/hemera/pkg/metric"
)

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
}
