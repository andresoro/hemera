package metric

import (
	"regexp"
	"strconv"
)

// ErrPacketMismatch ...

type _err string

func (e _err) Error() string { return string(e) }

const (
	ErrPacketMismatch _err = "packet does not match desired format"
	ErrInvalidType    _err = "metric type is not supported"
)

// nolint:gochecknoglobals
var (
	_re = regexp.MustCompile(`\A` +
		`(?P<name>[\w.]+)` +
		`:` +
		`(?P<sign>[-+])?` +
		`(?P<val>([0-9]*[.])?[0-9]+)` +
		`\|` +
		`(?P<type>\w+)` +
		`(\|@(?P<rate>\d+\.\d+))?`,
	)
	_reSubExpNames = _re.SubexpNames()
)

// metric type
type mType int

const (
	Counter mType = iota + 1
	Gauge
	Set
	Timer
)

// Metric is a published client event
type Metric struct {
	Name  string
	Value float64
	Type  mType
	Rate  float64
}

// Parse a byte array (from UDP packet) into a metric or return an error
func Parse(packet []byte) (*Metric, error) {
	matches := _re.FindStringSubmatch(string(packet))
	if matches == nil {
		return nil, ErrPacketMismatch
	}

	var newMetric Metric
	var sign float64

ForMatches:
	for i, match := range matches {
		if i == 0 {
			continue ForMatches
		}
		switch _reSubExpNames[i] {
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
			newMetric.Type, err = stringToMType(match)
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

	if sign == -1 && newMetric.Type != Counter {
		newMetric.Value *= -1
	}

	return &newMetric, nil
}

func stringToMType(s string) (mType, error) {
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
