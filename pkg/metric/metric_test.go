package metric

import "testing"

func TestMetricMatch(t *testing.T) {
	goodMetricString := []byte("performance.os.disk:1099511627776|g|@0.2")
	badMetricString := []byte("adadad:hajkdhakjd")

	metric, err := Parse(goodMetricString)
	if err != nil {
		t.Error("parse should pass on a valid metric string")
	}

	_, err = Parse(badMetricString)
	if err == nil {
		t.Error("parse should fail on bad string")
	}

	if metric.Name != "performance.os.disk" {
		t.Errorf("metric name not correct \n wanted: %s \n got: %s", "performance.os.disk", metric.Name)
	}

	if metric.Value != float64(1099511627776) {
		t.Errorf("metric value not correct \n wanted %f \n got %f", float64(1099511627776), metric.Value)
	}

	if metric.Type != Gauge {
		t.Error("incorrect type")
	}

	if metric.Rate != float64(0.2) {
		t.Error("incorrect rate")
	}
}
