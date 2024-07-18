package go_sql_driver

import (
	"fmt"
	"go.uber.org/zap"
	"regexp"
	"strings"
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

type TypeInfo struct {
	Type     int       `json:"type,omitempty"`
	Nullable bool      `json:"nullable,omitempty"`
	Array    bool      `json:"array,omitempty"`
	MapKey   *TypeInfo `json:"map_key,omitempty"`
	MapValue *TypeInfo `json:"map_value,omitempty"`
}

func (ty TypeInfo) ToString() string {
	if ty.Array {
		return fmt.Sprintf("Array<%s>", ty.MapKey.ToString())
	}
	if ty.MapKey != nil && ty.MapValue != nil {
		return fmt.Sprintf("Map<%s,%s>", ty.MapKey.ToString(), ty.MapValue.ToString())
	}
	return GetTypeName(ty.Type)
}

var (
	typeInfo             map[string]*TypeInfo
	lowCardinalityRegexp = regexp.MustCompile(`^LowCardinality\((.+)\)`)
)

func init() {
	typeInfo = make(map[string]*TypeInfo)
	for _, t := range []int{TINYINT, SMALLINT, INT, BIGINT, LARGEINT, DECIMAL, DOUBLE, FLOAT, BOOLEAN, CHAR, STRING, VARCHAR, BINARY, DATE, DATETIME} {
		tn := GetTypeName(t)
		typeInfo[tn] = &TypeInfo{Type: t}
		nullTn := fmt.Sprintf("%s null", tn)
		typeInfo[nullTn] = &TypeInfo{Type: t, Nullable: true}
		arrTn := fmt.Sprintf("Array<%s>", tn)
		typeInfo[arrTn] = &TypeInfo{Type: t, Array: true}
	}
}

const (
	Unknown = iota

	TINYINT
	SMALLINT
	INT
	BIGINT
	LARGEINT
	DECIMAL
	DOUBLE
	FLOAT
	BOOLEAN

	CHAR
	STRING
	VARCHAR
	BINARY

	DATE
	DATETIME

	ARRAY
	JSON
	MAP
	STRUCT

	BITMAP
	HLL

	Tinyint  = "tinyint"
	Smallint = "smallint"
	Int      = "int"
	Bigint   = "bigint"
	Largeint = "largeint"
	Decimal  = "decimal"
	Double   = "double"
	Float    = "float"
	Boolean  = "boolean"

	Char    = "char"
	String  = "string"
	Varchar = "varchar"
	Binary  = "binary"

	Date     = "date"
	Datetime = "datetime"

	Array  = "array"
	Json   = "json"
	Map    = "map"
	Struct = "struct"

	Bitmap = "bitmap"
	Hll    = "hll"
)

// GetTypeName returns the column type in ClickHouse
func GetTypeName(typ int) (name string) {
	switch typ {
	case TINYINT:
		name = Tinyint
	case SMALLINT:
		name = Smallint
	case INT:
		name = Int
	case BIGINT:
		name = Bigint
	case LARGEINT:
		name = Largeint
	case DECIMAL:
		name = Decimal
	case DOUBLE:
		name = Double
	case FLOAT:
		name = Float
	case BOOLEAN:
		name = Boolean

	case CHAR:
		name = Char
	case STRING:
		name = String
	case VARCHAR:
		name = Varchar
	case BINARY:
		name = Binary

	case DATE:
		name = Date
	case DATETIME:
		name = Datetime

	case ARRAY:
		name = Array
	case JSON:
		name = Json
	case MAP:
		name = Map
	case STRUCT:
		name = Struct

	case BITMAP:
		name = Bitmap
	case HLL:
		name = Hll

	default:
		name = "unknown"
	}
	return
}

func WhichType(typ string, nullable bool) (ti *TypeInfo) {
	typ = lowCardinalityRegexp.ReplaceAllString(typ, "$1")

	ti, ok := typeInfo[typ]
	if ok {
		return ti
	}
	origTyp := typ
	array := strings.HasPrefix(typ, "array<")
	var dataType int
	if array {
		typ = typ[len("array<") : len(typ)-1]
	}
	if strings.HasPrefix(typ, Tinyint) {
		dataType = TINYINT
	} else if strings.HasPrefix(typ, "smallint") {
		dataType = SMALLINT
	} else if strings.HasPrefix(typ, "int") {
		dataType = INT
	} else if strings.HasPrefix(typ, Bigint) {
		dataType = BIGINT
	} else if strings.HasPrefix(typ, Largeint) {
		dataType = LARGEINT
	} else if strings.HasPrefix(typ, Decimal) {
		dataType = DECIMAL
	} else if strings.HasPrefix(typ, Double) {
		dataType = DOUBLE
	} else if strings.HasPrefix(typ, Float) {
		dataType = FLOAT
	} else if strings.HasPrefix(typ, Bigint) {
		dataType = BIGINT
	} else if strings.HasPrefix(typ, Boolean) {
		dataType = BOOLEAN
	} else if strings.HasPrefix(typ, Char) {
		dataType = CHAR
	} else if strings.HasPrefix(typ, String) {
		dataType = STRING
	} else if strings.HasPrefix(typ, Varchar) {
		dataType = VARCHAR
	} else if strings.HasPrefix(typ, Binary) {
		dataType = BINARY
	} else if strings.HasPrefix(typ, Date) {
		dataType = DATE
	} else if strings.HasPrefix(typ, Datetime) {
		dataType = DATETIME
	} else if strings.HasPrefix(typ, Bitmap) {
		dataType = BITMAP
	} else if strings.HasPrefix(typ, Hll) {
		dataType = HLL
	} else if strings.HasPrefix(typ, Map) {
		dataType = MAP
		idx := strings.Index(typ, ", ")
		ti = &TypeInfo{
			Type:     dataType,
			Nullable: nullable,
			Array:    array,
			MapKey:   WhichType(typ[len("map<"):idx], false),
			MapValue: WhichType(typ[idx+2:len(typ)-1], false),
		}
		typeInfo[origTyp] = ti
		return ti
	} else {
		return nil
	}
	ti = &TypeInfo{Type: dataType, Nullable: nullable, Array: array}
	typeInfo[origTyp] = ti
	return ti
}

// Metric interface for metric collection
type Metric interface {
	GetTINYINT(key string, nullable bool) (val interface{})
	GetSMALLINT(key string, nullable bool) (val interface{})
	GetINT(key string, nullable bool) (val interface{})
	GetBIGINT(key string, nullable bool) (val interface{})
	GetLARGEINT(key string, nullable bool) (val interface{})
	GetDECIMAL(key string, nullable bool) (val interface{})
	GetDOUBLE(key string, nullable bool) (val interface{})
	GetFLOAT(key string, nullable bool) (val interface{})
	GetBOOLEAN(key string, nullable bool) (val interface{})
	GetCHAR(key string, nullable bool) (val interface{})
	GetSTRING(key string, nullable bool) (val interface{})
	GetVARCHAR(key string, nullable bool) (val interface{})
	GetBINARY(key string, nullable bool) (val interface{})
	GetDATE(key string, nullable bool) (val interface{})
	GetDATETIME(key string, nullable bool) (val interface{})
	GetARRAY(key string, typeinfo *TypeInfo) (val interface{})
	GetJSON(key string, t int) (val interface{})
	GetMAP(key string, nullable bool) (val interface{})
	GetBITMAP(key string, nullable bool) (val interface{})
	GetHLL(key string, nullable bool) (val interface{})
	GetSTRUCT(key string, nullable bool) (val interface{})
	GetNewKeys(knownKeys *sync.Map) map[string]string
	Close()
}

// DimMetrics
type DimMetrics struct {
	Dims   []*ColumnWithType
	Fields []*ColumnWithType
}

// ColumnWithType
type ColumnWithType struct {
	Name string
	Type *TypeInfo
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
