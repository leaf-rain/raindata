package rsql

import (
	"encoding/json"
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

// FastjsonParser parser for get data in json format
type FastjsonParser struct {
	pp           fastjson.ParserPool
	knownLayouts *sync.Map
	timeUnit     float64
	local        *time.Location
	logger       *zap.Logger
	metricPool   *sync.Pool
}

func (p *FastjsonParser) Parse(bs []byte) (metric Metric, err error) {
	var value *fastjson.Value
	var ps = p.pp.Get()
	value, err = ps.ParseBytes(bs)
	if err != nil {
		err = errors.Wrapf(err, "")
		return
	}
	var result = p.metricPool.Get().(*FastjsonMetric)
	result.parser = p
	result.value = value
	result.ps = ps
	return result, nil
}

func NewFastjsonParser(timeUnit float64, local *time.Location, logger *zap.Logger) (Parser, error) {
	buf := new(FastjsonParser)
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
	buf.metricPool = &sync.Pool{
		New: func() interface{} {
			return &FastjsonMetric{
				parser: buf,
			}
		},
	}
	return buf, nil
}

var _ Metric = (*FastjsonMetric)(nil)

type FastjsonMetric struct {
	parser *FastjsonParser
	value  *fastjson.Value
	ps     *fastjson.Parser
}

func (p *FastjsonMetric) Get(key string, nullable bool, ty int) interface{} {
	switch ty {
	case TINYINT:
		return p.GetTINYINT(key, nullable)
	case SMALLINT:
		return p.GetSMALLINT(key, nullable)
	case INT:
		return p.GetINT(key, nullable)
	case BIGINT:
		return p.GetBIGINT(key, nullable)
	case LARGEINT:
		return p.GetLARGEINT(key, nullable)
	case DECIMAL:
		return p.GetDECIMAL(key, nullable)
	case DOUBLE:
		return p.GetDOUBLE(key, nullable)
	case BOOLEAN:
		return p.GetBOOLEAN(key, nullable)
	case FLOAT:
		return p.GetFLOAT(key, nullable)
	case CHAR:
		return p.GetCHAR(key, nullable)
	case VARCHAR:
		return p.GetVARCHAR(key, nullable)
	case BINARY:
		return p.GetBINARY(key, nullable)
	case DATE:
		return p.GetDATE(key, nullable)
	case DATETIME:
		return p.GetDATETIME(key, nullable)
	case ARRAY:
		return p.GetARRAY(key, ty)
	case JSON:
		return p.GetJSON(key, ty)
	case MAP:
		return p.GetMAP(key, &TypeInfo{
			Type: ty,
		})
	case STRUCT:
		return p.GetSTRUCT(key, true)
	case HLL:
		return p.GetHLL(key, true)
	case BITMAP:
		return p.GetBITMAP(key, true)
	default:
		return nil
	}
}

func (p *FastjsonMetric) Set(key string, val interface{}) {
	js, _ := json.Marshal(val)
	obj, _ := fastjson.ParseBytes(js)
	p.value.Set(key, obj)
}

func (p *FastjsonMetric) Value() string {
	return p.value.String()
}

func (p *FastjsonMetric) GetTINYINT(key string, nullable bool) (val interface{}) {
	return FastjsonGetInt[int8](p.value.Get(key), nullable, math.MinInt8, math.MaxInt8)
}

func (p *FastjsonMetric) GetSMALLINT(key string, nullable bool) (val interface{}) {
	return FastjsonGetInt[int16](p.value.Get(key), nullable, math.MinInt16, math.MaxInt16)
}

func (p *FastjsonMetric) GetINT(key string, nullable bool) (val interface{}) {
	return FastjsonGetInt[int32](p.value.Get(key), nullable, math.MinInt32, math.MaxInt32)
}

func (p *FastjsonMetric) GetBIGINT(key string, nullable bool) (val interface{}) {
	return FastjsonGetInt[int64](p.value.Get(key), nullable, math.MinInt64, math.MaxInt64)
}

func (p *FastjsonMetric) GetLARGEINT(key string, nullable bool) (val interface{}) {
	return FastjsonGetInt[int64](p.value.Get(key), nullable, math.MinInt64, math.MaxInt64)
}

func (p *FastjsonMetric) GetDECIMAL(key string, nullable bool) (val interface{}) {
	return FastjsonGetFloat[float64](p.value.Get(key), nullable, math.MaxFloat64)
}

func (p *FastjsonMetric) GetDOUBLE(key string, nullable bool) (val interface{}) {
	return FastjsonGetFloat[float64](p.value.Get(key), nullable, math.MaxFloat64)
}

func (p *FastjsonMetric) GetFLOAT(key string, nullable bool) (val interface{}) {
	return FastjsonGetFloat[float32](p.value.Get(key), nullable, math.MaxFloat32)
}

func (p *FastjsonMetric) GetBOOLEAN(key string, nullable bool) (val interface{}) {
	return getString(p.value.Get(key), nullable)
}

func (p *FastjsonMetric) GetCHAR(key string, nullable bool) (val interface{}) {
	return getString(p.value.Get(key), nullable)
}

func (p *FastjsonMetric) GetSTRING(key string, nullable bool) (val interface{}) {
	return getString(p.value.Get(key), nullable)
}

func (p *FastjsonMetric) GetVARCHAR(key string, nullable bool) (val interface{}) {
	return getString(p.value.Get(key), nullable)
}

func (p *FastjsonMetric) GetBINARY(key string, nullable bool) (val interface{}) {
	return []byte(getString(p.value.Get(key), nullable).(string))
}

func (p *FastjsonMetric) GetDATE(key string, nullable bool) (val interface{}) {
	return getDateTime(p, key, p.value.Get(key), nullable)
}

func (p *FastjsonMetric) GetDATETIME(key string, nullable bool) (val interface{}) {
	return getDateTime(p, key, p.value.Get(key), nullable)
}

func (p *FastjsonMetric) GetARRAY(key string, ty int) (val interface{}) {
	return getArray(p, key, p.value.Get(key), ty)
}

func (p *FastjsonMetric) GetJSON(key string, t int) (val interface{}) {
	v := p.value.Get(key)
	val = val2map(v)
	return
}

func (p *FastjsonMetric) GetMAP(key string, typeinfo *TypeInfo) (val interface{}) {
	return getMap(p, p.value.Get(key), typeinfo)
}

func (p *FastjsonMetric) GetBITMAP(key string, nullable bool) (val interface{}) {
	//TODO implement me
	panic("implement me")
}

func (p *FastjsonMetric) GetHLL(key string, nullable bool) (val interface{}) {
	//TODO implement me
	panic("implement me")
}

func (p *FastjsonMetric) GetSTRUCT(key string, nullable bool) (val interface{}) {
	//TODO implement me
	panic("implement me")
}

func (p *FastjsonMetric) Close() {
	p.parser.pp.Put(p.ps)
	p.parser.metricPool.Put(p)
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
			if typ, arr := fjDetectType(v, 0); typ != 0 && !arr {
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

// Assuming that all values of a field of rkafka message has the same layout, and layouts of each field are unrelated.
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
	case BOOLEAN:
		arr := make([]bool, 0)
		for _, e := range array {
			v := e != nil && e.Type() == fastjson.TypeTrue
			arr = append(arr, v)
		}
		val = arr
	case TINYINT:
		val = FastjsonIntArray[int8](array, math.MinInt8, math.MaxInt8)
	case SMALLINT:
		val = FastjsonIntArray[int16](array, math.MinInt16, math.MaxInt16)
	case INT:
		val = FastjsonIntArray[int32](array, math.MinInt32, math.MaxInt32)
	case BIGINT:
		val = FastjsonIntArray[int64](array, math.MinInt64, math.MaxInt64)
	case LARGEINT:
		val = FastjsonIntArray[int64](array, math.MinInt64, math.MaxInt64)
	case DOUBLE:
		val = FastjsonFloatArray[float64](array, math.MaxFloat64)
	case FLOAT:
		val = FastjsonFloatArray[float64](array, math.MaxFloat32)
	case DECIMAL:
		arr := make([]decimal.Decimal, 0)
		for _, e := range array {
			v, _ := e.Float64()
			arr = append(arr, decimal.NewFromFloat(v))
		}
		val = arr
	case CHAR, STRING, VARCHAR, BINARY:
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
	case DATE, DATETIME:
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
	case STRUCT, MAP, JSON:
		arr := make([]map[string]interface{}, 0)
		for _, e := range array {
			m := val2map(e)
			if m != nil {
				arr = append(arr, m)
			}
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
	case TINYINT:
		if res, err := strconv.ParseInt(string(key), 10, 8); err == nil {
			return int8(res)
		} else {
			p.parser.logger.Error("failed to parse map key", zap.Error(err))
		}
	case SMALLINT:
		if res, err := strconv.ParseInt(string(key), 10, 16); err == nil {
			return int16(res)
		} else {
			p.parser.logger.Error("failed to parse map key", zap.Error(err))
		}
	case INT:
		if res, err := strconv.ParseInt(string(key), 10, 32); err == nil {
			return int32(res)
		} else {
			p.parser.logger.Error("failed to parse map key", zap.Error(err))
		}
	case BIGINT, LARGEINT:
		if res, err := strconv.ParseInt(string(key), 10, 64); err == nil {
			return int64(res)
		} else {
			p.parser.logger.Error("failed to parse map key", zap.Error(err))
		}

	case DATE, DATETIME:
		if res, err := ParseDateTime(string(key), string(key), p.parser.knownLayouts, p.parser.local); err == nil {
			return res
		} else {
			p.parser.logger.Error("failed to parse map key", zap.Error(err))
		}
		val = getDefaultDateTime(typeinfo.Nullable)
	case CHAR, VARCHAR, STRING:
		return string(key)
	case DECIMAL, DOUBLE, FLOAT:
		fallthrough
	case BOOLEAN:
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
		case BOOLEAN:
			val = getBool(value, typeinfo.Nullable)
		case TINYINT:
			val = FastjsonGetInt[int8](value, typeinfo.Nullable, math.MinInt8, math.MaxInt8)
		case SMALLINT:
			val = FastjsonGetInt[int16](value, typeinfo.Nullable, math.MinInt16, math.MaxInt16)
		case INT:
			val = FastjsonGetInt[int32](value, typeinfo.Nullable, math.MinInt32, math.MaxInt32)
		case BIGINT, LARGEINT:
			val = FastjsonGetInt[int64](value, typeinfo.Nullable, math.MinInt64, math.MaxInt64)
		case FLOAT:
			val = FastjsonGetFloat[float32](value, typeinfo.Nullable, math.MaxFloat32)
		case DOUBLE:
			val = FastjsonGetFloat[float64](value, typeinfo.Nullable, math.MaxFloat64)
		case DECIMAL:
			val = getDecimal(value, typeinfo.Nullable)
		case DATE, DATETIME:
			val = getDateTime(p, sourcename, value, typeinfo.Nullable)
		case CHAR, VARCHAR, STRING, BINARY:
			val = getString(value, typeinfo.Nullable)
		case MAP:
			val = getMap(p, value, typeinfo)
		case STRUCT, JSON:
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
		typ = BOOLEAN
	case fastjson.TypeNumber:
		typ = DOUBLE
		if _, err := v.Int64(); err == nil {
			typ = BIGINT
		}
	case fastjson.TypeString:
		typ = STRING
		if val, err := v.StringBytes(); err == nil {
			if _, layout := parseInLocation(string(val), time.Local); layout != "" {
				typ = DATETIME
			}
		}
	case fastjson.TypeObject:
		typ = STRUCT
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
