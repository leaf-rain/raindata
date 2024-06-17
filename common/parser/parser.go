package parser

import (
	"github.com/leaf-rain/fastjson"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"sync"
	"time"
)

// Parse is the Parser interface
type Parser interface {
	Parse(bs []byte) (metric Metric, err error)
}

// Pool may be used for pooling Parsers for similarly typed JSONs.
type Pool struct {
	name         string
	csvFormat    map[string]int
	delimiter    string
	timeZone     *time.Location
	timeUnit     float64
	knownLayouts sync.Map
	pool         sync.Pool
	once         sync.Once // only need to detect new keys from fields once
	fields       string
	logger       *zap.Logger
}

// NewParserPool creates a parser pool
func NewParserPool(name string, csvFormat []string, delimiter string, timezone string, timeunit float64, fields string, logger *zap.Logger) (pp *Pool, err error) {
	var tz *time.Location
	if timezone == "" {
		tz = time.Local
	} else if tz, err = time.LoadLocation(timezone); err != nil {
		err = errors.Wrapf(err, "")
		return
	}
	pp = &Pool{
		name:      name,
		delimiter: delimiter,
		timeZone:  tz,
		timeUnit:  timeunit,
		fields:    fields,
		logger:    logger,
	}
	if csvFormat != nil {
		pp.csvFormat = make(map[string]int, len(csvFormat))
		for i, title := range csvFormat {
			pp.csvFormat[title] = i
		}
	}
	return
}

// Get returns a Parser from pp.
//
// The Parser must be Put to pp after use.
func (pp *Pool) Get() (Parser, error) {
	v := pp.pool.Get()
	if v == nil {
		switch pp.name {
		case "gjson":
			return &GjsonParser{pp: pp}, nil
		case "csv":
			if pp.fields != "" {
				pp.logger.Warn("extra fields for csv parser is not supported, fields ignored")
			}
			return &CsvParser{pp: pp}, nil
		case "fastjson":
			fallthrough
		default:
			var obj *fastjson.Object
			if pp.fields != "" {
				value, err := fastjson.Parse(pp.fields)
				if err != nil {
					err = errors.Wrapf(err, "failed to parse fields as a valid json object")
					return nil, err
				}
				obj, err = value.Object()
				if err != nil {
					err = errors.Wrapf(err, "failed to retrive fields member")
					return nil, err
				}
			}
			return &FastjsonParser{pp: pp, fields: obj}, nil
		}
	}
	return v.(Parser), nil
}

// Put returns p to pp.
//
// p and objects recursively returned from p cannot be used after p
// is put into pp.
func (pp *Pool) Put(p Parser) {
	pp.pool.Put(p)
}

// Assuming that all values of a field of kafka message has the same layout, and layouts of each field are unrelated.
// Automatically detect the layout from till the first successful detection and reuse that layout forever.
// Return time in UTC.
func (pp *Pool) ParseDateTime(key string, val string) (t time.Time, err error) {
	var layout string
	var lay interface{}
	var ok bool
	var t2 time.Time
	if val == "" {
		err = ErrParseDateTime
		return
	}
	if lay, ok = pp.knownLayouts.Load(key); !ok {
		t2, layout = parseInLocation(val, pp.timeZone)
		if layout == "" {
			err = ErrParseDateTime
			return
		}
		t = t2
		pp.knownLayouts.Store(key, layout)
		return
	}
	if layout, ok = lay.(string); !ok {
		err = ErrParseDateTime
		return
	}
	if t2, err = time.ParseInLocation(layout, val, pp.timeZone); err != nil {
		err = ErrParseDateTime
		return
	}
	t = t2.UTC()
	return
}
