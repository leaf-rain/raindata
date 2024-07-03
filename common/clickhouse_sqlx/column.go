package clickhouse_sqlx

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
)

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

// ColumnWithType
type ColumnWithType struct {
	Name string    `json:"name,omitempty"`
	Type *TypeInfo `json:"type,omitempty"`
}

type TypeInfo struct {
	Type     int       `json:"type,omitempty"`
	Nullable bool      `json:"nullable,omitempty"`
	Array    bool      `json:"array,omitempty"`
	MapKey   *TypeInfo `json:"map_key,omitempty"`
	MapValue *TypeInfo `json:"map_value,omitempty"`
}

func (ty TypeInfo) ToString() string {
	if ty.Array {
		return fmt.Sprintf("Array(%s)", ty.MapKey.ToString())
	}
	if ty.MapKey != nil && ty.MapValue != nil {
		return fmt.Sprintf("Map(%s,%s)", ty.MapKey.ToString(), ty.MapValue.ToString())
	}
	return GetTypeName(ty.Type)
}

type Columns struct {
	rwLock sync.RWMutex
	keys   []string
	values map[string]*ColumnWithType
}

func (om *Columns) Get(key string) *ColumnWithType {
	om.rwLock.RLock()
	defer om.rwLock.RUnlock()
	if value, present := om.values[key]; present {
		return value
	}
	return nil
}

func (om *Columns) Put(key string, value *ColumnWithType) {
	om.rwLock.Lock()
	defer om.rwLock.Unlock()
	if _, present := om.values[key]; present {
		om.values[key] = value
		return
	}
	om.keys = append(om.keys, key)
	om.values[key] = value
}

func (om *Columns) Keys() []string {
	om.rwLock.RLock()
	defer om.rwLock.RUnlock()
	result := make([]string, len(om.keys))
	for i := range om.keys {
		result[i] = om.keys[i]
	}
	return result
}

func (om *Columns) GetValues() map[string]*ColumnWithType {
	om.rwLock.RLock()
	defer om.rwLock.RUnlock()
	result := make(map[string]*ColumnWithType, len(om.keys))
	for i := range om.keys {
		result[om.keys[i]] = om.values[om.keys[i]]
	}
	return result
}

func NewColumns() *Columns {
	om := Columns{}
	om.keys = []string{}
	om.values = map[string]*ColumnWithType{}
	return &om
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
	GetNewKeys(knownKeys *sync.Map) map[string]string
}

func GetValueByType(metric Metric, cwt *ColumnWithType) (val interface{}) {
	name := cwt.Name
	if cwt.Type.Array {
		val = metric.GetArray(name, cwt.Type.Type)
	} else {
		switch cwt.Type.Type {
		case Bool:
			val = metric.GetBool(name, cwt.Type.Nullable)
		case Int8:
			val = metric.GetInt8(name, cwt.Type.Nullable)
		case Int16:
			val = metric.GetInt16(name, cwt.Type.Nullable)
		case Int32:
			val = metric.GetInt32(name, cwt.Type.Nullable)
		case Int64:
			val = metric.GetInt64(name, cwt.Type.Nullable)
		case UInt8:
			val = metric.GetUint8(name, cwt.Type.Nullable)
		case UInt16:
			val = metric.GetUint16(name, cwt.Type.Nullable)
		case UInt32:
			val = metric.GetUint32(name, cwt.Type.Nullable)
		case UInt64:
			val = metric.GetUint64(name, cwt.Type.Nullable)
		case Float32:
			val = metric.GetFloat32(name, cwt.Type.Nullable)
		case Float64:
			val = metric.GetFloat64(name, cwt.Type.Nullable)
		case Decimal:
			val = metric.GetDecimal(name, cwt.Type.Nullable)
		case DateTime:
			val = metric.GetDateTime(name, cwt.Type.Nullable)
		case String:
			val = metric.GetString(name, cwt.Type.Nullable)
		case Map:
			val = metric.GetMap(name, cwt.Type)
		case Object:
			val = metric.GetObject(name, cwt.Type.Nullable)
		case IPv4:
			val = metric.GetIPv4(name, cwt.Type.Nullable)
		case IPv6:
			val = metric.GetIPv6(name, cwt.Type.Nullable)
		default:
			log.Fatal("LOGIC ERROR: reached switch default condition")
		}
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
		log.Fatal(fmt.Sprintf("ClickHouse column type %v is not inside supported ones(case-sensitive): %v", origTyp, typeInfo))
	}
	ti = &TypeInfo{Type: dataType, Nullable: nullable, Array: array}
	typeInfo[origTyp] = ti
	return ti
}

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
