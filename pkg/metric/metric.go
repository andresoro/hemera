package metric

import (
	"errors"
	"regexp"
	"strconv"
)

// ErrPacketMismatch ...
var ErrPacketMismatch = errors.New("packet does not match desired format")
var ErrInvalidType = errors.New("metric type is not supported")

// metric type
type mtype int

const (
	Counter mtype = iota
	Gauge   mtype = iota
	Set     mtype = iota
	Timer   mtype = iota
)

// Metric is a published client event
type Metric struct {
	Name  string
	Value float64
	Type  mtype
	Rate  float64
}

// Parse a byte array (from UDP packet) into a metric or return an error
func Parse(packet []byte) (*Metric, error) {

	re := regexp.MustCompile(`\A(?P<name>[\w\.]+):((?P<sign>\-|\+))?(?P<val>([0-9]*[.])?[0-9]+)\|(?P<type>\w+)(\|@(?P<rate>\d+\.\d+))?`)

	if !re.Match(packet) {
		return nil, ErrPacketMismatch
	}

	newMetric := &Metric{}
	var sign float64

	matches := re.FindStringSubmatch(string(packet))
	names := re.SubexpNames()
	for i, match := range matches {
		if i != 0 {
			switch names[i] {
			case "name":
				newMetric.Name = match
			case "sign":
				if match == "" {
					sign = 1
				}
				if match == "-" {
					sign = -1
				}
			case "val":
				val, err := strconv.ParseFloat(match, 64)
				if err != nil {
					return nil, err
				}
				newMetric.Value = val
			case "type":
				var err error
				newMetric.Type, err = stringToMtype(match)
				if err != nil {
					return nil, err
				}
			case "rate":
				if match == "" {
					continue
				} else {
					rate, err := strconv.ParseFloat(match, 64)
					if err != nil {
						return nil, err
					}
					newMetric.Rate = rate
				}
			}
		}
	}

	if sign == -1 && newMetric.Type != Counter {
		newMetric.Value = newMetric.Value * -1
	}

	return newMetric, nil
}

func stringToMtype(s string) (mtype, error) {
	switch s {
	case "c":
		return Counter, nil
	case "g":
		return Gauge, nil
	case "ms":
		return Timer, nil
	case "s":
		return Set, nil
	default:
		return 0, ErrInvalidType
	}
}
