package parser

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"
)

type TypeInfo struct {
	Type     int
	Nullable bool
	Array    bool
	MapKey   *TypeInfo
	MapValue *TypeInfo
}

var (
	typeInfo             map[string]*TypeInfo
	lowCardinalityRegexp = regexp.MustCompile(`^LowCardinality\((.+)\)`)
)

func init() {
	typeInfo = make(map[string]*TypeInfo)
	for _, t := range []int{Bool, Int8, Int16, Int32, Int64, UInt8, UInt16, UInt32, UInt64, Float32, Float64, DateTime, String, Object, IPv4, IPv6} {
		tn := GetTypeName(t)
		typeInfo[tn] = &TypeInfo{Type: t}
		nullTn := fmt.Sprintf("Nullable(%s)", tn)
		typeInfo[nullTn] = &TypeInfo{Type: t, Nullable: true}
		arrTn := fmt.Sprintf("Array(%s)", tn)
		typeInfo[arrTn] = &TypeInfo{Type: t, Array: true}
	}
	typeInfo["UUID"] = &TypeInfo{Type: String}
	typeInfo["Nullable(UUID)"] = &TypeInfo{Type: String, Nullable: true}
	typeInfo["Array(UUID)"] = &TypeInfo{Type: String, Array: true}
	typeInfo["Date"] = &TypeInfo{Type: DateTime}
	typeInfo["Nullable(Date)"] = &TypeInfo{Type: DateTime, Nullable: true}
	typeInfo["Array(Date)"] = &TypeInfo{Type: DateTime, Array: true}
}

const (
	Unknown = iota
	Bool
	Int8
	Int16
	Int32
	Int64
	UInt8
	UInt16
	UInt32
	UInt64
	Float32
	Float64
	Decimal
	DateTime
	String
	Object
	Map
	IPv4
	IPv6
)

// GetTypeName returns the column type in ClickHouse
func GetTypeName(typ int) (name string) {
	switch typ {
	case Bool:
		name = "Bool"
	case Int8:
		name = "Int8"
	case Int16:
		name = "Int16"
	case Int32:
		name = "Int32"
	case Int64:
		name = "Int64"
	case UInt8:
		name = "UInt8"
	case UInt16:
		name = "UInt16"
	case UInt32:
		name = "UInt32"
	case UInt64:
		name = "UInt64"
	case Float32:
		name = "Float32"
	case Float64:
		name = "Float64"
	case Decimal:
		name = "Decimal"
	case DateTime:
		name = "DateTime"
	case String:
		name = "String"
	case Object:
		name = "Object('json')"
	case Map:
		name = "Map"
	case IPv4:
		name = "IPv4"
	case IPv6:
		name = "IPv6"
	default:
		name = "Unknown"
	}
	return
}

func WhichType(typ string) (ti *TypeInfo) {
	typ = lowCardinalityRegexp.ReplaceAllString(typ, "$1")

	ti, ok := typeInfo[typ]
	if ok {
		return ti
	}
	origTyp := typ
	nullable := strings.HasPrefix(typ, "Nullable(")
	array := strings.HasPrefix(typ, "Array(")
	var dataType int
	if nullable {
		typ = typ[len("Nullable(") : len(typ)-1]
	} else if array {
		typ = typ[len("Array(") : len(typ)-1]
	}
	if strings.HasPrefix(typ, "DateTime64") {
		dataType = DateTime
	} else if strings.HasPrefix(typ, "Decimal") {
		dataType = Decimal
	} else if strings.HasPrefix(typ, "FixedString") {
		dataType = String
	} else if strings.HasPrefix(typ, "Enum8(") {
		dataType = String
	} else if strings.HasPrefix(typ, "Enum16(") {
		dataType = String
	} else if strings.HasPrefix(typ, "Map") {
		dataType = Map
		idx := strings.Index(typ, ", ")
		ti = &TypeInfo{
			Type:     dataType,
			Nullable: nullable,
			Array:    array,
			MapKey:   WhichType(typ[len("Map("):idx]),
			MapValue: WhichType(typ[idx+2 : len(typ)-1]),
		}
		typeInfo[origTyp] = ti
		return ti
	} else {
		log.Fatalf(fmt.Sprintf("ClickHouse column type %v is not inside supported ones(case-sensitive): %v", origTyp, typeInfo))
	}
	ti = &TypeInfo{Type: dataType, Nullable: nullable, Array: array}
	typeInfo[origTyp] = ti
	return ti
}

// Metric interface for metric collection
type Metric interface {
	GetBool(key string, nullable bool) (val interface{})
	GetInt8(key string, nullable bool) (val interface{})
	GetInt16(key string, nullable bool) (val interface{})
	GetInt32(key string, nullable bool) (val interface{})
	GetInt64(key string, nullable bool) (val interface{})
	GetUint8(key string, nullable bool) (val interface{})
	GetUint16(key string, nullable bool) (val interface{})
	GetUint32(key string, nullable bool) (val interface{})
	GetUint64(key string, nullable bool) (val interface{})
	GetFloat32(key string, nullable bool) (val interface{})
	GetFloat64(key string, nullable bool) (val interface{})
	GetDecimal(key string, nullable bool) (val interface{})
	GetDateTime(key string, nullable bool) (val interface{})
	GetString(key string, nullable bool) (val interface{})
	GetObject(key string, nullable bool) (val interface{})
	GetMap(key string, typeinfo *TypeInfo) (val interface{})
	GetArray(key string, t int) (val interface{})
	GetIPv4(key string, nullable bool) (val interface{})
	GetIPv6(key string, nullable bool) (val interface{})
	GetNewKeys(knownKeys, newKeys, warnKeys *sync.Map, white, black *regexp.Regexp, partition int, offset int64) bool
}

// DimMetrics
type DimMetrics struct {
	Dims   []*ColumnWithType
	Fields []*ColumnWithType
}

// ColumnWithType
type ColumnWithType struct {
	Name       string
	Type       *TypeInfo
	SourceName string
}

// struct for ingesting a clickhouse Map type value
type OrderedMap struct {
	keys   []interface{}
	values map[interface{}]interface{}
}

func (om *OrderedMap) Get(key interface{}) (interface{}, bool) {
	if value, present := om.values[key]; present {
		return value, present
	}
	return nil, false
}

func (om *OrderedMap) Put(key interface{}, value interface{}) {
	if _, present := om.values[key]; present {
		om.values[key] = value
		return
	}
	om.keys = append(om.keys, key)
	om.values[key] = value
}

func (om *OrderedMap) Keys() <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		defer close(ch)
		for _, key := range om.keys {
			ch <- key
		}
	}()
	return ch
}

func (om *OrderedMap) GetValues() map[interface{}]interface{} {
	return om.values
}

func NewOrderedMap() *OrderedMap {
	om := OrderedMap{}
	om.keys = []interface{}{}
	om.values = map[interface{}]interface{}{}
	return &om
}

type SeriesQuota struct {
	sync.Mutex     `json:"-"`
	NextResetQuota time.Time
	BmSeries       map[int64]int64 // sid:mid
	WrSeries       int
	Birth          time.Time
}
