package parser

import (
	"errors"
	"github.com/leaf-rain/raindata/common/go_sql_driver"
	"go.uber.org/zap"
	"math"
	"time"
)

var (
	Layouts = []string{
		//DateTime, RFC3339
		"2006-01-02T15:04:05Z07:00", //time.RFC3339, `date --iso-8601=s` on Ubuntu 20.04
		"2006-01-02T15:04:05Z0700",  //`date --iso-8601=s` on CentOS 7.6
		"2006-01-02T15:04:05",
		//DateTime, ISO8601
		"2006-01-02 15:04:05Z07:00", //`date --rfc-3339=s` output format
		"2006-01-02 15:04:05Z0700",
		"2006-01-02 15:04:05",
		//DateTime, other layouts supported by golang
		"Mon Jan _2 15:04:05 2006",        //time.ANSIC
		"Mon Jan _2 15:04:05 MST 2006",    //time.UnixDate
		"Mon Jan 02 15:04:05 -0700 2006",  //time.RubyDate
		"02 Jan 06 15:04 MST",             //time.RFC822
		"02 Jan 06 15:04 -0700",           //time.RFC822Z
		"Monday, 02-Jan-06 15:04:05 MST",  //time.RFC850
		"Mon, 02 Jan 2006 15:04:05 MST",   //time.RFC1123
		"Mon, 02 Jan 2006 15:04:05 -0700", //time.RFC1123Z
		//DateTime, linux utils
		"Mon Jan 02 15:04:05 MST 2006",    // `date` on CentOS 7.6 default output format
		"Mon 02 Jan 2006 03:04:05 PM MST", // `date` on Ubuntu 20.4 default output format
		//DateTime, home-brewed
		"Jan 02, 2006 15:04:05Z07:00",
		"Jan 02, 2006 15:04:05Z0700",
		"Jan 02, 2006 15:04:05",
		"02/Jan/2006 15:04:05 Z07:00",
		"02/Jan/2006 15:04:05 Z0700",
		"02/Jan/2006 15:04:05",
		//Date
		"2006-01-02",
		"02/01/2006",
		"02/Jan/2006",
		"Jan 02, 2006",
		"Mon Jan 02, 2006",
	}
	Epoch            = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	ErrParseDateTime = errors.New("value doesn't contain DateTime")
)

func UnixFloat(sec, unit float64) (t time.Time) {
	sec *= unit
	//2^32 seconds since epoch: 2106-02-07T06:28:16Z
	if sec < 0 || sec >= 4294967296.0 {
		return Epoch
	}
	i, f := math.Modf(sec)
	return time.Unix(int64(i), int64(f*1e9))
}

// Parse is the Parser interface
type Parser interface {
	Parse(bs []byte) (metric go_sql_driver.Metric, err error)
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
	//case "gjson":
	//	return NewGjsonParser(cf.TimeUnit, cf.Local, cf.Logger), nil
	//case "csv":
	//	return NewCsvParser(cf.CsvFormat, cf.Delimiter, cf.TimeUnit, cf.Local, cf.Logger), nil
	//case "fastjson":
	//	fallthrough
	//default:
	//	return NewFastjsonParser(cf.TimeUnit, cf.Local, cf.Logger)
	}
	return nil, nil
}
