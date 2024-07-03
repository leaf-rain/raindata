package parser

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"golang.org/x/exp/constraints"
	"math"
	"net"
	"strconv"
	"sync"
	"time"
)

var _ Parser = (*GjsonParser)(nil)

type GjsonParser struct {
	knownLayouts *sync.Map
	pool         sync.Pool
	timeUnit     float64
	local        *time.Location
	logger       *zap.Logger
}

func NewGjsonParser(timeUnit float64, local *time.Location, logger *zap.Logger) *GjsonParser {
	buf := new(GjsonParser)
	buf.knownLayouts = &sync.Map{}
	buf.logger = logger
	buf.timeUnit = 1
	if timeUnit != 0 {
		buf.timeUnit = timeUnit
	}
	if local != nil {
		buf.local = local
	} else {
		buf.local = time.Local
	}
	buf.pool = sync.Pool{
		New: func() interface{} {
			return &GjsonMetric{
				parser: buf,
			}
		},
	}
	return buf
}

func (p *GjsonParser) Parse(bs []byte) (metric Metric, err error) {
	metric = &GjsonMetric{p, string(bs)}
	return
}

type GjsonMetric struct {
	parser *GjsonParser
	raw    string
}

func (c *GjsonMetric) getField(key string) gjson.Result {
	ret := gjson.Get(c.raw, key)
	if !ret.Exists() {
		ret = gjson.Get(c.raw, key)
	}
	return ret
}

func (c *GjsonMetric) GetString(key string, nullable bool) (val interface{}) {
	return getGJsonString(c.getField(key), nullable)
}

func (c *GjsonMetric) GetBool(key string, nullable bool) (val interface{}) {
	return getGJsonBool(c.getField(key), nullable)
}

func (c *GjsonMetric) GetDecimal(key string, nullable bool) (val interface{}) {
	return getGJsonDecimal(c.getField(key), nullable)
}

func (c *GjsonMetric) GetInt8(key string, nullable bool) (val interface{}) {
	return GjsonGetInt[int8](c, c.getField(key), nullable, math.MinInt8, math.MaxInt8)
}

func (c *GjsonMetric) GetInt16(key string, nullable bool) (val interface{}) {
	return GjsonGetInt[int16](c, c.getField(key), nullable, math.MinInt16, math.MaxInt16)
}

func (c *GjsonMetric) GetInt32(key string, nullable bool) (val interface{}) {
	return GjsonGetInt[int32](c, c.getField(key), nullable, math.MinInt32, math.MaxInt32)
}

func (c *GjsonMetric) GetInt64(key string, nullable bool) (val interface{}) {
	return GjsonGetInt[int64](c, c.getField(key), nullable, math.MinInt64, math.MaxInt64)
}

func (c *GjsonMetric) GetUint8(key string, nullable bool) (val interface{}) {
	return GjsonGetUint[uint8](c, c.getField(key), nullable, math.MaxUint8)
}

func (c *GjsonMetric) GetUint16(key string, nullable bool) (val interface{}) {
	return GjsonGetUint[uint16](c, c.getField(key), nullable, math.MaxUint16)
}

func (c *GjsonMetric) GetUint32(key string, nullable bool) (val interface{}) {
	return GjsonGetUint[uint32](c, c.getField(key), nullable, math.MaxUint32)
}

func (c *GjsonMetric) GetUint64(key string, nullable bool) (val interface{}) {
	return GjsonGetUint[uint64](c, c.getField(key), nullable, math.MaxUint64)
}

func (c *GjsonMetric) GetFloat32(key string, nullable bool) (val interface{}) {
	return GjsonGetFloat[float32](c, c.getField(key), nullable, math.MaxFloat32)
}

func (c *GjsonMetric) GetFloat64(key string, nullable bool) (val interface{}) {
	return GjsonGetFloat[float64](c, c.getField(key), nullable, math.MaxFloat64)
}

func (c *GjsonMetric) GetIPv4(key string, nullable bool) (val interface{}) {
	return getGJsonIPv4(c, c.getField(key), nullable)
}

func (c *GjsonMetric) GetIPv6(key string, nullable bool) (val interface{}) {
	return getGJsonIPv6(c, c.getField(key), nullable)
}

func (c *GjsonMetric) GetNewKeys(knownKeys *sync.Map) map[string]string {
	var result = make(map[string]string)
	ite := func(k, v gjson.Result) bool {
		strKey := k.Str
		if _, loaded := knownKeys.Load(strKey); !loaded {
			if typ, array := gjDetectType(v, 0); typ != Unknown && typ != Object && !array {
				result[strKey] = GetTypeName(typ)
			}
		}
		return true
	}
	gjson.Parse(c.raw).ForEach(ite)
	return result
}

func GjsonGetInt[T constraints.Signed](c *GjsonMetric, r gjson.Result, nullable bool, min, max int64) (val interface{}) {
	if !gjCompatibleInt(r) {
		val = getDefaultInt[T](nullable)
		return
	}
	switch r.Type {
	case gjson.True:
		val = T(1)
	case gjson.False:
		val = T(0)
	case gjson.Number:
		if val2 := r.Int(); float64(val2) != r.Num {
			val = getDefaultInt[T](nullable)
		} else if val2 < min {
			val = T(min)
		} else if val2 > max {
			val = T(max)
		} else {
			val = T(val2)
		}
	default:
		val = getDefaultInt[T](nullable)
	}
	return
}

func GjsonGetUint[T constraints.Unsigned](c *GjsonMetric, r gjson.Result, nullable bool, max uint64) (val interface{}) {
	if !gjCompatibleInt(r) {
		val = getDefaultInt[T](nullable)
		return
	}
	switch r.Type {
	case gjson.True:
		val = T(1)
	case gjson.False:
		val = T(0)
	case gjson.Number:
		if val2 := r.Uint(); float64(val2) != r.Num {
			val = getDefaultInt[T](nullable)
		} else if val2 > max {
			val = T(max)
		} else {
			val = T(val2)
		}
	default:
		val = getDefaultInt[T](nullable)
	}
	return
}

func GjsonGetFloat[T constraints.Float](c *GjsonMetric, r gjson.Result, nullable bool, max float64) (val interface{}) {
	if !gjCompatibleFloat(r) {
		val = getDefaultFloat[T](nullable)
		return
	}
	switch r.Type {
	case gjson.Number:
		if r.Num > max {
			val = T(max)
		} else {
			val = T(r.Num)
		}
	default:
		val = getDefaultFloat[T](nullable)
	}
	return
}

func (c *GjsonMetric) GetDateTime(key string, nullable bool) (val interface{}) {
	return getGJsonDateTime(c, key, c.getField(key), nullable)
}

func (c *GjsonMetric) GetObject(key string, nullable bool) (val interface{}) {
	return
}

func (c *GjsonMetric) GetArray(key string, typ int) (val interface{}) {
	return getGJsonArray(c, key, c.getField(key), typ)
}

func (c *GjsonMetric) GetMap(key string, typeinfo *TypeInfo) (val interface{}) {
	return getGJsonMap(c, c.getField(key), typeinfo)
}

func (c *GjsonMetric) val2OrderedMap(v gjson.Result, typeinfo *TypeInfo) (m *OrderedMap) {
	m = NewOrderedMap()
	v.ForEach(func(k, v gjson.Result) bool {
		str := k.String()
		rawKey := c.castResultByType(str, k, typeinfo.MapKey)
		m.Put(rawKey, c.castResultByType(str, v, typeinfo.MapValue))
		return true
	})
	return
}

func GjsonIntArray[T constraints.Signed](a []gjson.Result, min, max int64) (arr []T) {
	arr = make([]T, 0, len(a))
	var val T
	for _, e := range a {
		switch e.Type {
		case gjson.True:
			val = T(1)
		case gjson.Number:
			if val2 := e.Int(); float64(val2) != e.Num {
				val = T(0)
			} else if val2 < min {
				val = T(min)
			} else if val2 > max {
				val = T(max)
			} else {
				val = T(val2)
			}
		default:
			val = T(0)
		}
		arr = append(arr, val)
	}
	return
}

func GjsonUintArray[T constraints.Unsigned](a []gjson.Result, max uint64) (arr []T) {
	arr = make([]T, 0, len(a))
	var val T
	for _, e := range a {
		switch e.Type {
		case gjson.True:
			val = T(1)
		case gjson.Number:
			if val2 := e.Uint(); float64(val2) != e.Num {
				val = T(0)
			} else if val2 > max {
				val = T(max)
			} else {
				val = T(val2)
			}
		default:
			val = T(0)
		}
		arr = append(arr, val)
	}
	return
}

func GjsonFloatArray[T constraints.Float](a []gjson.Result, max float64) (arr []T) {
	arr = make([]T, 0, len(a))
	var val T
	for _, e := range a {
		switch e.Type {
		case gjson.Number:
			if e.Num > max {
				val = T(max)
			} else {
				val = T(e.Num)
			}
		default:
			val = T(0.0)
		}
		arr = append(arr, val)
	}
	return
}

func gjCompatibleBool(r gjson.Result) (ok bool) {
	if !r.Exists() {
		return
	}
	switch r.Type {
	case gjson.True, gjson.False:
		ok = true
	default:
	}
	return
}

func gjCompatibleInt(r gjson.Result) (ok bool) {
	if !r.Exists() {
		return
	}
	switch r.Type {
	case gjson.True, gjson.False, gjson.Number:
		ok = true
	default:
	}
	return
}

func gjCompatibleFloat(r gjson.Result) (ok bool) {
	if !r.Exists() {
		return
	}
	switch r.Type {
	case gjson.Number:
		ok = true
	default:
	}
	return
}

func gjCompatibleDateTime(r gjson.Result) (ok bool) {
	if !r.Exists() {
		return
	}
	switch r.Type {
	case gjson.Number, gjson.String:
		ok = true
	default:
	}
	return
}

func gjDetectType(v gjson.Result, depth int) (typ int, array bool) {
	typ = Unknown
	if depth > 1 {
		return
	}
	switch v.Type {
	case gjson.True, gjson.False:
		typ = Bool
	case gjson.Number:
		if _, err := strconv.ParseInt(v.Raw, 10, 64); err == nil {
			typ = Int64
		} else {
			typ = Float64
		}
	case gjson.String:
		typ = String
		if _, layout := parseInLocation(v.Str, time.Local); layout != "" {
			typ = DateTime
		}
	case gjson.JSON:
		if v.IsArray() {
			if depth >= 1 {
				return
			}
			array = true
			if array := v.Array(); len(array) != 0 {
				typ, _ = gjDetectType(array[0], depth+1)
			}
		} else if v.IsObject() {
			typ = Object
		}
	default:
	}
	return
}

func (c *GjsonMetric) castResultByType(sourcename string, value gjson.Result, typeinfo *TypeInfo) (val interface{}) {
	if typeinfo.Array {
		val = getGJsonArray(c, sourcename, value, typeinfo.Type)
		return
	} else {
		switch typeinfo.Type {
		case Bool:
			val = getGJsonBool(value, typeinfo.Nullable)
		case Int8:
			val = GjsonGetInt[int8](c, value, typeinfo.Nullable, math.MinInt8, math.MaxInt8)
		case Int16:
			val = GjsonGetInt[int16](c, value, typeinfo.Nullable, math.MinInt16, math.MaxInt16)
		case Int32:
			val = GjsonGetInt[int32](c, value, typeinfo.Nullable, math.MinInt32, math.MaxInt32)
		case Int64:
			val = GjsonGetInt[int64](c, value, typeinfo.Nullable, math.MinInt64, math.MaxInt64)
		case UInt8:
			val = GjsonGetUint[uint8](c, value, typeinfo.Nullable, math.MaxUint8)
		case UInt16:
			val = GjsonGetUint[uint16](c, value, typeinfo.Nullable, math.MaxUint16)
		case UInt32:
			val = GjsonGetUint[uint32](c, value, typeinfo.Nullable, math.MaxUint32)
		case UInt64:
			val = GjsonGetUint[uint64](c, value, typeinfo.Nullable, math.MaxUint64)
		case IPv6:
			val = getGJsonIPv6(c, value, typeinfo.Nullable)
		case IPv4:
			val = getGJsonIPv4(c, value, typeinfo.Nullable)
		case Float32:
			val = GjsonGetFloat[float32](c, value, typeinfo.Nullable, math.MaxFloat32)
		case Float64:
			val = GjsonGetFloat[float64](c, value, typeinfo.Nullable, math.MaxFloat64)
		case Decimal:
			val = getGJsonDecimal(value, typeinfo.Nullable)
		case DateTime:
			val = getGJsonDateTime(c, sourcename, value, typeinfo.Nullable)
		case String:
			val = getGJsonString(value, typeinfo.Nullable)
		case Map:
			val = getGJsonMap(c, value, typeinfo)
		default:
			c.parser.logger.Fatal("LOGIC ERROR: reached switch default condition")
		}
	}
	return
}

func getGJsonDecimal(r gjson.Result, nullable bool) (val interface{}) {
	if !gjCompatibleFloat(r) {
		val = getDefaultDecimal(nullable)
		return
	}
	switch r.Type {
	case gjson.Number:
		val = decimal.NewFromFloat(r.Num)
	default:
		val = getDefaultDecimal(nullable)
	}
	return
}

func getGJsonBool(r gjson.Result, nullable bool) (val interface{}) {
	if !gjCompatibleBool(r) {
		val = getDefaultBool(nullable)
		return
	}
	val = (r.Type == gjson.True)
	return
}

func getGJsonString(r gjson.Result, nullable bool) (val interface{}) {
	if !r.Exists() || r.Type == gjson.Null {
		if nullable {
			return
		}
		val = ""
		return
	}
	switch r.Type {
	case gjson.Null:
		val = ""
	case gjson.String:
		val = r.Str
	default:
		val = r.Raw
	}
	return
}

func getGJsonIPv4(c *GjsonMetric, r gjson.Result, nullable bool) (val interface{}) {
	if !r.Exists() || r.Type == gjson.Null {
		if nullable {
			return
		}
		val = ""
		return
	}
	switch r.Type {
	case gjson.Null:
		val = ""
	case gjson.String:
		s := r.Str
		if net.ParseIP(s) != nil {
			val = s
		} else {
			val = net.IPv4zero.String()
		}
	case gjson.Number:
		val = GjsonGetUint[uint32](c, r, nullable, math.MaxUint32)
	default:
		val = net.IPv4zero.String()
	}
	return
}

func getGJsonIPv6(c *GjsonMetric, r gjson.Result, nullable bool) (val interface{}) {
	if !r.Exists() || r.Type == gjson.Null {
		if nullable {
			return
		}
		val = ""
		return
	}
	switch r.Type {
	case gjson.Null:
		val = ""
	case gjson.String:
		s := r.Str
		if net.ParseIP(s) != nil {
			val = s
		} else {
			val = net.IPv6zero.String()
		}
	default:
		val = net.IPv6zero.String()
	}
	return
}

func getGJsonDateTime(c *GjsonMetric, key string, r gjson.Result, nullable bool) (val interface{}) {
	if !gjCompatibleDateTime(r) {
		val = getDefaultDateTime(nullable)
		return
	}
	switch r.Type {
	case gjson.Number:
		val = UnixFloat(r.Num, c.parser.timeUnit)
	case gjson.String:
		var err error
		if val, err = ParseDateTime(key, r.Str, c.parser.knownLayouts, c.parser.local); err != nil {
			val = getDefaultDateTime(nullable)
		}
	default:
		val = getDefaultDateTime(nullable)
	}
	return
}

func getGJsonArray(c *GjsonMetric, key string, r gjson.Result, typ int) (val interface{}) {
	var array []gjson.Result
	if r.IsArray() {
		array = r.Array()
	}
	switch typ {
	case Bool:
		results := make([]bool, 0, len(array))
		for _, e := range array {
			v := (e.Exists() && e.Type == gjson.True)
			results = append(results, v)
		}
		val = results
	case Int8:
		val = GjsonIntArray[int8](array, math.MinInt8, math.MaxInt8)
	case Int16:
		val = GjsonIntArray[int16](array, math.MinInt16, math.MaxInt16)
	case Int32:
		val = GjsonIntArray[int32](array, math.MinInt32, math.MaxInt32)
	case Int64:
		val = GjsonIntArray[int64](array, math.MinInt64, math.MaxInt64)
	case UInt8:
		val = GjsonUintArray[uint8](array, math.MaxUint8)
	case UInt16:
		val = GjsonUintArray[uint16](array, math.MaxUint16)
	case UInt32:
		val = GjsonUintArray[uint32](array, math.MaxUint32)
	case UInt64:
		val = GjsonUintArray[uint64](array, math.MaxUint64)
	case Float32:
		val = GjsonFloatArray[float32](array, math.MaxFloat32)
	case Float64:
		val = GjsonFloatArray[float64](array, math.MaxFloat64)
	case Decimal:
		results := make([]decimal.Decimal, 0, len(array))
		var f float64
		for _, e := range array {
			switch e.Type {
			case gjson.Number:
				f = e.Num
			default:
				f = float64(0.0)
			}
			results = append(results, decimal.NewFromFloat(f))
		}
		val = results
	case String:
		results := make([]string, 0, len(array))
		var s string
		for _, e := range array {
			switch e.Type {
			case gjson.Null:
				s = ""
			case gjson.String:
				s = e.Str
			default:
				s = e.Raw
			}
			results = append(results, s)
		}
		val = results
	case DateTime:
		results := make([]time.Time, 0, len(array))
		var t time.Time
		for _, e := range array {
			switch e.Type {
			case gjson.Number:
				t = UnixFloat(e.Num, c.parser.timeUnit)
			case gjson.String:
				var err error
				if t, err = ParseDateTime(key, e.Str, c.parser.knownLayouts, c.parser.local); err != nil {
					t = Epoch
				}
			default:
				t = Epoch
			}
			results = append(results, t)
		}
		val = results
	case IPv4:
		arr := make([]interface{}, 0)
		for _, e := range array {
			v := getGJsonIPv4(c, e, false)
			arr = append(arr, v)
		}
		val = arr
	case IPv6:
		arr := make([]interface{}, 0)
		for _, e := range array {
			v := getGJsonIPv6(c, e, false)
			arr = append(arr, v)
		}
		val = arr
	default:
		c.parser.logger.Fatal(fmt.Sprintf("LOGIC ERROR: unsupported array type %v", typ))
	}
	return
}

func getGJsonMap(c *GjsonMetric, r gjson.Result, typeinfo *TypeInfo) (val interface{}) {
	if r.Type == gjson.JSON {
		val = c.val2OrderedMap(r, typeinfo)
	}
	return
}
