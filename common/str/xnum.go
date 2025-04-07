package str

import (
	"math"
	"strconv"
)

// IsNumericAndConvertToInt64 判断是否为数值类型，并转换为 int64
func IsNumericAndConvertToInt64(value interface{}) (int64, bool) {
	switch v := value.(type) {
	case int:
		return int64(v), true
	case int8:
		return int64(v), true
	case int16:
		return int64(v), true
	case int32:
		return int64(v), true
	case int64:
		return v, true
	case uint:
		return int64(v), true
	case uint8:
		return int64(v), true
	case uint16:
		return int64(v), true
	case uint32:
		return int64(v), true
	case uint64:
		if v > math.MaxUint64 {
			return 0, false // 超出 int64 范围
		}
		return int64(v), true
	default:
		return 0, false
	}
}

func IsNumericOrStringToInt64(value interface{}) (int64, bool) {
	v, ok := value.(string)
	if ok {
		tmp, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			return tmp, true
		} else {
			return 0, false
		}
	} else {
		return IsNumericAndConvertToInt64(value)
	}
}
