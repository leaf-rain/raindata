package parser

import (
	"go.uber.org/zap"
	"sync"
	"time"
)

// Parse is the Parser interface
type Parser interface {
	Parse(bs []byte) (metric Metric, err error)
}

type ParserConfig struct {
	Ty        string
	CsvFormat []string
	Delimiter string
	TimeUnit  float64
	Local     *time.Location
	Logger    *zap.Logger
}

func NewParse(cf ParserConfig) (Parser, error) {
	switch cf.Ty {
	case "gjson":
		return NewGjsonParser(cf.TimeUnit, cf.Local, cf.Logger), nil
	case "csv":
		return NewCsvParser(cf.CsvFormat, cf.Delimiter, cf.TimeUnit, cf.Local, cf.Logger), nil
	case "fastjson":
		fallthrough
	default:
		return NewFastjsonParser(cf.TimeUnit, cf.Local, cf.Logger)
	}
}

// Assuming that all values of a field of rkafka message has the same layout, and layouts of each field are unrelated.
// Automatically detect the layout from till the first successful detection and reuse that layout forever.
// Return time in UTC.
func ParseDateTime(key string, val string, knownLayouts *sync.Map, timeLocation *time.Location) (t time.Time, err error) {
	var layout string
	var lay interface{}
	var ok bool
	var t2 time.Time
	if val == "" {
		err = ErrParseDateTime
		return
	}
	if lay, ok = knownLayouts.Load(key); !ok {
		t2, layout = parseInLocation(val, timeLocation)
		if layout == "" {
			err = ErrParseDateTime
			return
		}
		t = t2
		knownLayouts.Store(key, layout)
		return
	}
	if layout, ok = lay.(string); !ok {
		err = ErrParseDateTime
		return
	}
	if t2, err = time.ParseInLocation(layout, val, timeLocation); err != nil {
		err = ErrParseDateTime
		return
	}
	t = t2.UTC()
	return
}
