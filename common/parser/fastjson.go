package parser

import (
	"fmt"
	"github.com/leaf-rain/fastjson"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"golang.org/x/exp/constraints"
	"math"
	"net"
	"strconv"
	"sync"
	"time"
)

var _ Parser = (*FastjsonParser)(nil)
var EmpytObject = make(map[string]interface{})

// FastjsonParser, parser for get data in json format
type FastjsonParser struct {
	fjp          fastjson.Parser
	fields       *fastjson.Object
	value        *fastjson.Value
	knownLayouts *sync.Map
	pool         sync.Pool
	fieldsStr    string
	timeUnit     float64
	local        *time.Location
	logger       *zap.Logger
}

func NewFastjsonParser(fieldsStr string, timeUnit float64, local *time.Location, logger *zap.Logger) (*FastjsonParser, error) {
	buf := new(FastjsonParser)
	buf.fieldsStr = fieldsStr
	buf.knownLayouts = &sync.Map{}
	buf.logger = logger
	if fieldsStr != "" {
		value, err := fastjson.Parse(fieldsStr)
		if err != nil {
			err = errors.Wrapf(err, "failed to parse fields as a valid json object")
			return nil, err
		}
		buf.fields, err = value.Object()
		if err != nil {
			err = errors.Wrapf(err, "failed to retrive fields member")
			return nil, err
		}
	}
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
			return &FastjsonMetric{
				parser: buf,
			}
		},
	}
	return buf, nil
}

type FastjsonMetric struct {
	parser *FastjsonParser
	value  *fastjson.Value
}

func (p *FastjsonParser) Parse(bs []byte) (metric Metric, err error) {
	var value *fastjson.Value
	if value, err = p.fjp.ParseBytes(bs); err != nil {
		err = errors.Wrapf(err, "")
		return
	}

	if p.fields != nil {
		p.fields.Visit(func(key []byte, v *fastjson.Value) {
			value.Set(string(key), v)
		})
	}
	buf, _ := p.pool.Get().(*FastjsonMetric)
	if buf == nil {
		buf = new(FastjsonMetric)
	}
	buf.parser = p
	buf.value = value
	return
}

func (p *FastjsonMetric) GetString(key string, nullable bool) (val interface{}) {
	return getString(p.value.Get(key), nullable)
}

func (p *FastjsonMetric) GetBool(key string, nullable bool) interface{} {
	return getBool(p.value.Get(key), nullable)
}

func (p *FastjsonMetric) GetDecimal(key string, nullable bool) (val interface{}) {
	return getDecimal(p.value.Get(key), nullable)
}

func (p *FastjsonMetric) GetInt8(key string, nullable bool) (val interface{}) {
	return FastjsonGetInt[int8](p.value.Get(key), nullable, math.MinInt8, math.MaxInt8)
}

func (p *FastjsonMetric) GetInt16(key string, nullable bool) (val interface{}) {
	return FastjsonGetInt[int16](p.value.Get(key), nullable, math.MinInt16, math.MaxInt16)
}

func (p *FastjsonMetric) GetInt32(key string, nullable bool) (val interface{}) {
	return FastjsonGetInt[int32](p.value.Get(key), nullable, math.MinInt32, math.MaxInt32)
}

func (p *FastjsonMetric) GetInt64(key string, nullable bool) (val interface{}) {
	return FastjsonGetInt[int64](p.value.Get(key), nullable, math.MinInt64, math.MaxInt64)
}

func (p *FastjsonMetric) GetUint8(key string, nullable bool) (val interface{}) {
	return FastjsonGetUint[uint8](p.value.Get(key), nullable, math.MaxUint8)
}

func (p *FastjsonMetric) GetUint16(key string, nullable bool) (val interface{}) {
	return FastjsonGetUint[uint16](p.value.Get(key), nullable, math.MaxUint16)
}

func (p *FastjsonMetric) GetUint32(key string, nullable bool) (val interface{}) {
	return FastjsonGetUint[uint32](p.value.Get(key), nullable, math.MaxUint32)
}

func (p *FastjsonMetric) GetUint64(key string, nullable bool) (val interface{}) {
	return FastjsonGetUint[uint64](p.value.Get(key), nullable, math.MaxUint64)
}

func (p *FastjsonMetric) GetFloat32(key string, nullable bool) (val interface{}) {
	return FastjsonGetFloat[float32](p.value.Get(key), nullable, math.MaxFloat32)
}

func (p *FastjsonMetric) GetFloat64(key string, nullable bool) (val interface{}) {
	return FastjsonGetFloat[float64](p.value.Get(key), nullable, math.MaxFloat64)
}

func (p *FastjsonMetric) GetIPv4(key string, nullable bool) (val interface{}) {
	return getIPv4(p.value.Get(key), nullable)
}

func (p *FastjsonMetric) GetIPv6(key string, nullable bool) (val interface{}) {
	return getIPv6(p.value.Get(key), nullable)
}

func (p *FastjsonMetric) GetNewKeys(knownKeys *sync.Map) map[string]string {
	var obj *fastjson.Object
	var err error
	var result = make(map[string]string)
	if obj, err = p.value.Object(); err != nil {
		return result
	}
	obj.Visit(func(key []byte, v *fastjson.Value) {
		strKey := string(key)
		if _, loaded := knownKeys.Load(strKey); !loaded {
			if typ, arr := fjDetectType(v, 0); typ != Unknown && typ != Object && !arr {
				result[strKey] = GetTypeName(typ)
			}
		}
	})
	return result
}

func FastjsonGetInt[T constraints.Signed](v *fastjson.Value, nullable bool, min, max int64) (val interface{}) {
	if !fjCompatibleInt(v) {
		val = getDefaultInt[T](nullable)
		return
	}
	switch v.Type() {
	case fastjson.TypeTrue:
		val = T(1)
	case fastjson.TypeFalse:
		val = T(0)
	default:
		if val2, err := v.Int64(); err != nil {
			val = getDefaultInt[T](nullable)
		} else if val2 < min {
			val = T(min)
		} else if val2 > max {
			val = T(max)
		} else {
			val = T(val2)
		}
	}
	return
}

func FastjsonGetUint[T constraints.Unsigned](v *fastjson.Value, nullable bool, max uint64) (val interface{}) {
	if !fjCompatibleInt(v) {
		val = getDefaultInt[T](nullable)
		return
	}
	switch v.Type() {
	case fastjson.TypeTrue:
		val = T(1)
	case fastjson.TypeFalse:
		val = T(0)
	default:
		if val2, err := v.Uint64(); err != nil {
			val = getDefaultInt[T](nullable)
		} else if val2 > max {
			val = T(max)
		} else {
			val = T(val2)
		}
	}
	return
}

func FastjsonGetFloat[T constraints.Float](v *fastjson.Value, nullable bool, max float64) (val interface{}) {
	if !fjCompatibleFloat(v) {
		val = getDefaultFloat[T](nullable)
		return
	}
	if val2, err := v.Float64(); err != nil {
		val = getDefaultFloat[T](nullable)
	} else if val2 > max {
		val = T(max)
	} else {
		val = T(val2)
	}
	return
}

func (p *FastjsonMetric) GetDateTime(key string, nullable bool) (val interface{}) {
	return getDateTime(p, key, p.value.Get(key), nullable)
}

func (p *FastjsonMetric) GetObject(key string, nullable bool) (val interface{}) {
	v := p.value.Get(key)
	val = val2map(v)
	return
}

func (p *FastjsonMetric) GetArray(key string, typ int) (val interface{}) {
	return getArray(p, key, p.value.Get(key), typ)
}

func (p *FastjsonMetric) GetMap(key string, typeinfo *TypeInfo) (val interface{}) {
	return getMap(p, p.value.Get(key), typeinfo)
}

func (p *FastjsonMetric) val2OrderedMap(v *fastjson.Value, typeinfo *TypeInfo) (m *OrderedMap) {
	var err error
	var obj *fastjson.Object
	m = NewOrderedMap()
	if v == nil {
		return
	}
	if obj, err = v.Object(); err != nil {
		return
	}
	obj.Visit(func(key []byte, v *fastjson.Value) {
		rawKey := p.castMapKeyByType(key, typeinfo.MapKey)
		m.Put(rawKey, p.castMapValueByType(string(key), v, typeinfo.MapValue))
	})
	return
}

func val2map(v *fastjson.Value) (m map[string]interface{}) {
	var err error
	var obj *fastjson.Object
	m = EmpytObject
	if v == nil {
		return
	}
	if obj, err = v.Object(); err != nil {
		return
	}
	m = make(map[string]interface{}, obj.Len())
	obj.Visit(func(key []byte, v *fastjson.Value) {
		strKey := string(key)
		switch v.Type() {
		case fastjson.TypeString:
			var vb []byte
			if vb, err = v.StringBytes(); err != nil {
				return
			}
			m[strKey] = string(vb)
		case fastjson.TypeNumber:
			var f float64
			if f, err = v.Float64(); err != nil {
				return
			}
			m[strKey] = f
		}
	})
	return
}

func getString(v *fastjson.Value, nullable bool) (val interface{}) {
	if v == nil || v.Type() == fastjson.TypeNull {
		if nullable {
			return
		}
		val = ""
		return
	}
	switch v.Type() {
	case fastjson.TypeString:
		b, _ := v.StringBytes()
		val = string(b)
	default:
		val = v.String()
	}
	return
}

func getBool(v *fastjson.Value, nullable bool) (val interface{}) {
	if !fjCompatibleBool(v) {
		val = getDefaultBool(nullable)
		return
	}
	val = (v.Type() == fastjson.TypeTrue)
	return
}

func getIPv4(v *fastjson.Value, nullable bool) (val interface{}) {
	if v == nil || v.Type() == fastjson.TypeNull {
		if nullable {
			return
		}
		val = ""
		return
	}
	switch v.Type() {
	case fastjson.TypeString:
		b, _ := v.StringBytes()
		s := string(b)
		if net.ParseIP(s) != nil {
			val = s
		} else {
			val = net.IPv4zero.String()
		}
	case fastjson.TypeNumber:
		val = FastjsonGetUint[uint32](v, nullable, math.MaxUint32)
	default:
		val = net.IPv4zero.String()
	}
	return
}

func getIPv6(v *fastjson.Value, nullable bool) (val interface{}) {
	if v == nil || v.Type() == fastjson.TypeNull {
		if nullable {
			return
		}
		val = ""
		return
	}
	switch v.Type() {
	case fastjson.TypeString:
		b, _ := v.StringBytes()
		s := string(b)
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

func getDecimal(v *fastjson.Value, nullable bool) (val interface{}) {
	if !fjCompatibleFloat(v) {
		val = getDefaultDecimal(nullable)
		return
	}
	if val2, err := v.Float64(); err != nil {
		val = getDefaultDecimal(nullable)
	} else {
		val = decimal.NewFromFloat(val2)
	}
	return
}

func getDateTime(p *FastjsonMetric, sourcename string, v *fastjson.Value, nullable bool) (val interface{}) {
	if !fjCompatibleDateTime(v) {
		val = getDefaultDateTime(nullable)
		return
	}
	var err error
	switch v.Type() {
	case fastjson.TypeNumber:
		var f float64
		if f, err = v.Float64(); err != nil {
			val = getDefaultDateTime(nullable)
			return
		}
		timeStamp := int64(f)
		val = time.Unix(timeStamp, 0)
	case fastjson.TypeString:
		var b []byte
		if b, err = v.StringBytes(); err != nil || len(b) == 0 {
			val = getDefaultDateTime(nullable)
			return
		}
		if val, err = p.ParseDateTime(sourcename, string(b)); err != nil {
			val = getDefaultDateTime(nullable)
		}
	default:
		val = getDefaultDateTime(nullable)
	}
	return
}

// Assuming that all values of a field of kafka message has the same layout, and layouts of each field are unrelated.
// Automatically detect the layout from till the first successful detection and reuse that layout forever.
// Return time in UTC.
func (p *FastjsonMetric) ParseDateTime(key string, val string) (t time.Time, err error) {
	var layout string
	var lay interface{}
	var ok bool
	var t2 time.Time
	if val == "" {
		err = ErrParseDateTime
		return
	}
	if lay, ok = p.parser.knownLayouts.Load(key); !ok {
		t2, layout = parseInLocation(val, time.Local)
		if layout == "" {
			err = ErrParseDateTime
			return
		}
		t = t2
		p.parser.knownLayouts.Store(key, layout)
		return
	}
	if layout, ok = lay.(string); !ok {
		err = ErrParseDateTime
		return
	}
	if t2, err = time.ParseInLocation(layout, val, time.Local); err != nil {
		err = ErrParseDateTime
		return
	}
	t = t2.UTC()
	return
}

func parseInLocation(val string, loc *time.Location) (t time.Time, layout string) {
	var err error
	var lay string
	for _, lay = range Layouts {
		if t, err = time.ParseInLocation(lay, val, loc); err == nil {
			t = t.UTC()
			layout = lay
			return
		}
	}
	return
}

func getArray(p *FastjsonMetric, sourcename string, v *fastjson.Value, typ int) (val interface{}) {
	var array []*fastjson.Value
	if v != nil {
		array, _ = v.Array()
	}
	switch typ {
	case Bool:
		arr := make([]bool, 0)
		for _, e := range array {
			v := e != nil && e.Type() == fastjson.TypeTrue
			arr = append(arr, v)
		}
		val = arr
	case Int8:
		val = FastjsonIntArray[int8](array, math.MinInt8, math.MaxInt8)
	case Int16:
		val = FastjsonIntArray[int16](array, math.MinInt16, math.MaxInt16)
	case Int32:
		val = FastjsonIntArray[int32](array, math.MinInt32, math.MaxInt32)
	case Int64:
		val = FastjsonIntArray[int64](array, math.MinInt64, math.MaxInt64)
	case UInt8:
		val = FastjsonUintArray[uint8](array, math.MaxUint8)
	case UInt16:
		val = FastjsonUintArray[uint16](array, math.MaxUint16)
	case UInt32:
		val = FastjsonUintArray[uint32](array, math.MaxUint32)
	case UInt64:
		val = FastjsonUintArray[uint64](array, math.MaxUint64)
	case Float32:
		val = FastjsonFloatArray[float32](array, math.MaxFloat32)
	case Float64:
		val = FastjsonFloatArray[float64](array, math.MaxFloat64)
	case Decimal:
		arr := make([]decimal.Decimal, 0)
		for _, e := range array {
			v, _ := e.Float64()
			arr = append(arr, decimal.NewFromFloat(v))
		}
		val = arr
	case String:
		arr := make([]string, 0)
		var s string
		for _, e := range array {
			switch e.Type() {
			case fastjson.TypeNull:
				s = ""
			case fastjson.TypeString:
				b, _ := e.StringBytes()
				s = string(b)
			default:
				s = e.String()
			}
			arr = append(arr, s)
		}
		val = arr
	case DateTime:
		arr := make([]time.Time, 0)
		var t time.Time
		for _, e := range array {
			switch e.Type() {
			case fastjson.TypeNumber:
				if f, err := e.Float64(); err != nil {
					t = Epoch
				} else {
					t = UnixFloat(f, p.parser.timeUnit)
				}
			case fastjson.TypeString:
				if b, err := e.StringBytes(); err != nil || len(b) == 0 {
					t = Epoch
				} else {
					var err error
					if t, err = ParseDateTime(sourcename, string(b), p.parser.knownLayouts, p.parser.local); err != nil {
						t = Epoch
					}
				}
			default:
				t = Epoch
			}
			arr = append(arr, t)
		}
		val = arr
	case Object:
		arr := make([]map[string]interface{}, 0)
		for _, e := range array {
			m := val2map(e)
			if m != nil {
				arr = append(arr, m)
			}
		}
		val = arr
	case IPv4:
		arr := make([]interface{}, 0)
		for _, e := range array {
			v := getIPv4(e, false)
			arr = append(arr, v)
		}
		val = arr
	case IPv6:
		arr := make([]interface{}, 0)
		for _, e := range array {
			v := getIPv6(e, false)
			arr = append(arr, v)
		}
		val = arr
	default:
		p.parser.logger.Fatal(fmt.Sprintf("LOGIC ERROR: unsupported array type %v", typ))
	}
	return
}

func getMap(p *FastjsonMetric, v *fastjson.Value, typeinfo *TypeInfo) (val interface{}) {
	if v == nil || v.Type() == fastjson.TypeObject {
		val = p.val2OrderedMap(v, typeinfo)
	} else {
		p.parser.logger.Fatal(fmt.Sprintf("SOURCE ERROR: unsupported map type: %v", v.Type()))
	}
	return
}

func (p *FastjsonMetric) castMapKeyByType(key []byte, typeinfo *TypeInfo) (val interface{}) {
	switch typeinfo.Type {
	case Int8:
		if res, err := strconv.ParseInt(string(key), 10, 8); err == nil {
			return int8(res)
		} else {
			p.parser.logger.Error("failed to parse map key", zap.Error(err))
		}
	case Int16:
		if res, err := strconv.ParseInt(string(key), 10, 16); err == nil {
			return int16(res)
		} else {
			p.parser.logger.Error("failed to parse map key", zap.Error(err))
		}
	case Int32:
		if res, err := strconv.ParseInt(string(key), 10, 32); err == nil {
			return int32(res)
		} else {
			p.parser.logger.Error("failed to parse map key", zap.Error(err))
		}
	case Int64:
		if res, err := strconv.ParseInt(string(key), 10, 64); err == nil {
			return int64(res)
		} else {
			p.parser.logger.Error("failed to parse map key", zap.Error(err))
		}
	case UInt8:
		if res, err := strconv.ParseUint(string(key), 10, 8); err == nil {
			return uint8(res)
		} else {
			p.parser.logger.Error("failed to parse map key", zap.Error(err))
		}
	case UInt16:
		if res, err := strconv.ParseUint(string(key), 10, 16); err == nil {
			return uint16(res)
		} else {
			p.parser.logger.Error("failed to parse map key", zap.Error(err))
		}
	case UInt32:
		if res, err := strconv.ParseUint(string(key), 10, 32); err == nil {
			return uint32(res)
		} else {
			p.parser.logger.Error("failed to parse map key", zap.Error(err))
		}
	case UInt64:
		if res, err := strconv.ParseUint(string(key), 10, 64); err == nil {
			return uint64(res)
		} else {
			p.parser.logger.Error("failed to parse map key", zap.Error(err))
		}
	case DateTime:
		if res, err := ParseDateTime(string(key), string(key), p.parser.knownLayouts, p.parser.local); err == nil {
			return res
		} else {
			p.parser.logger.Error("failed to parse map key", zap.Error(err))
		}
		val = getDefaultDateTime(typeinfo.Nullable)
	case String:
		return string(key)
	case Float32:
		fallthrough
	case Float64:
		fallthrough
	case Bool:
		p.parser.logger.Fatal("unsupported map key type")
	default:
		p.parser.logger.Fatal("LOGIC ERROR: reached switch default condition")
	}

	return
}

func (p *FastjsonMetric) castMapValueByType(sourcename string, value *fastjson.Value, typeinfo *TypeInfo) (val interface{}) {
	if typeinfo.Array {
		val = getArray(p, sourcename, value, typeinfo.Type)
		return
	} else {
		switch typeinfo.Type {
		case Bool:
			val = getBool(value, typeinfo.Nullable)
		case Int8:
			val = FastjsonGetInt[int8](value, typeinfo.Nullable, math.MinInt8, math.MaxInt8)
		case Int16:
			val = FastjsonGetInt[int16](value, typeinfo.Nullable, math.MinInt16, math.MaxInt16)
		case Int32:
			val = FastjsonGetInt[int32](value, typeinfo.Nullable, math.MinInt32, math.MaxInt32)
		case Int64:
			val = FastjsonGetInt[int64](value, typeinfo.Nullable, math.MinInt64, math.MaxInt64)
		case UInt8:
			val = FastjsonGetUint[uint8](value, typeinfo.Nullable, math.MaxUint8)
		case UInt16:
			val = FastjsonGetUint[uint16](value, typeinfo.Nullable, math.MaxUint16)
		case UInt32:
			val = FastjsonGetUint[uint32](value, typeinfo.Nullable, math.MaxUint32)
		case UInt64:
			val = FastjsonGetUint[uint64](value, typeinfo.Nullable, math.MaxUint64)
		case IPv4:
			val = getIPv4(value, typeinfo.Nullable)
		case IPv6:
			val = getIPv6(value, typeinfo.Nullable)
		case Float32:
			val = FastjsonGetFloat[float32](value, typeinfo.Nullable, math.MaxFloat32)
		case Float64:
			val = FastjsonGetFloat[float64](value, typeinfo.Nullable, math.MaxFloat64)
		case Decimal:
			val = getDecimal(value, typeinfo.Nullable)
		case DateTime:
			val = getDateTime(p, sourcename, value, typeinfo.Nullable)
		case String:
			val = getString(value, typeinfo.Nullable)
		case Map:
			val = getMap(p, value, typeinfo)
		case Object:
			val = val2map(value)
		default:
			p.parser.logger.Fatal("LOGIC ERROR: reached switch default condition")
		}
	}
	return
}

func FastjsonIntArray[T constraints.Signed](a []*fastjson.Value, min, max int64) (arr []T) {
	arr = make([]T, 0)
	var val T
	for _, e := range a {
		if e.Type() == fastjson.TypeTrue {
			val = T(1)
		} else {
			val2, _ := e.Int64()
			if val2 < min {
				val = T(min)
			} else if val2 > max {
				val = T(max)
			} else {
				val = T(val2)
			}
		}
		arr = append(arr, val)
	}
	return
}

func FastjsonUintArray[T constraints.Unsigned](a []*fastjson.Value, max uint64) (arr []T) {
	arr = make([]T, 0)
	var val T
	for _, e := range a {
		if e.Type() == fastjson.TypeTrue {
			val = T(1)
		} else {
			val2, _ := e.Uint64()
			if val2 > max {
				val = T(max)
			} else {
				val = T(val2)
			}
		}
		arr = append(arr, val)
	}
	return
}

func FastjsonFloatArray[T constraints.Float](a []*fastjson.Value, max float64) (arr []T) {
	arr = make([]T, 0)
	var val T
	for _, e := range a {
		val2, _ := e.Float64()
		if val2 > max {
			val = T(max)
		} else {
			val = T(val2)
		}
		arr = append(arr, val)
	}
	return
}

func fjCompatibleBool(v *fastjson.Value) (ok bool) {
	if v == nil {
		return
	}
	switch v.Type() {
	case fastjson.TypeTrue, fastjson.TypeFalse:
		ok = true
	}
	return
}

func fjCompatibleInt(v *fastjson.Value) (ok bool) {
	if v == nil {
		return
	}
	switch v.Type() {
	case fastjson.TypeTrue, fastjson.TypeFalse, fastjson.TypeNumber:
		ok = true
	}
	return
}

func fjCompatibleFloat(v *fastjson.Value) (ok bool) {
	if v == nil {
		return
	}
	switch v.Type() {
	case fastjson.TypeNumber:
		ok = true
	}
	return
}

func fjCompatibleDateTime(v *fastjson.Value) (ok bool) {
	if v == nil {
		return
	}
	switch v.Type() {
	case fastjson.TypeNumber, fastjson.TypeString:
		ok = true
	}
	return
}

func getDefaultBool(nullable bool) (val interface{}) {
	if nullable {
		return
	}
	val = false
	return
}

func getDefaultInt[T constraints.Integer](nullable bool) (val interface{}) {
	if nullable {
		return
	}
	var zero T
	val = zero
	return
}

func getDefaultFloat[T constraints.Float](nullable bool) (val interface{}) {
	if nullable {
		return
	}
	val = T(0.0)
	return
}

func getDefaultDecimal(nullable bool) (val interface{}) {
	if nullable {
		return
	}
	val = decimal.NewFromInt(0)
	return
}

func getDefaultDateTime(nullable bool) (val interface{}) {
	if nullable {
		return
	}
	val = Epoch
	return
}

func fjDetectType(v *fastjson.Value, depth int) (typ int, array bool) {
	typ = Unknown
	if depth > 1 {
		return
	}
	switch v.Type() {
	case fastjson.TypeNull:
		typ = Unknown
	case fastjson.TypeTrue, fastjson.TypeFalse:
		typ = Bool
	case fastjson.TypeNumber:
		typ = Float64
		if _, err := v.Int64(); err == nil {
			typ = Int64
		}
	case fastjson.TypeString:
		typ = String
		if val, err := v.StringBytes(); err == nil {
			if _, layout := parseInLocation(string(val), time.Local); layout != "" {
				typ = DateTime
			}
		}
	case fastjson.TypeObject:
		typ = Object
	case fastjson.TypeArray:
		if depth >= 1 {
			return
		}
		array = true
		if arr, err := v.Array(); err == nil && len(arr) > 0 {
			typ, _ = fjDetectType(arr[0], depth+1)
		}
	default:
	}
	return
}
