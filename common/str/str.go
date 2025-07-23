package str

import (
	"path/filepath"
	"strconv"
	"strings"
)

func Str2Int32(data string) int32 {
	result, _ := strconv.ParseInt(data, 10, 64)
	return int32(result)
}

func Str2Int64(data string) int64 {
	result, _ := strconv.ParseInt(data, 10, 64)
	return result
}

func GetFileName(path string) string {
	// 使用filepath.Base获取路径的基本文件名（包含后缀）
	baseName := filepath.Base(path)
	// 使用path.Ext获取文件的扩展名，并将其从baseName中移除以得到不带后缀的文件名
	return strings.TrimSuffix(baseName, filepath.Ext(baseName))
}

func RemoveRootPrefix(old, data string) string {
	data = strings.TrimLeft(data, "./")
	index := strings.Index(old, data)
	ext := filepath.Ext(old)
	if index != -1 {
		return RemoveSuffix(old[index+len(data):], ext)
	}
	return RemoveSuffix(old, ext)
}

func RemovePrefix(old, data string) string {
	if len(data) == 0 {
		return old
	}
	index := strings.Index(old, data)
	if index != -1 {
		return old[index+len(data):]
	}
	return old
}

func RemoveSuffix(old, data string) string {
	if len(data) == 0 {
		return old
	}
	index := strings.Index(old, data)
	if index != -1 {
		return old[:index]
	}
	return old
}

func Add(str ...string) string {
	var builder strings.Builder
	for _, s := range str {
		builder.WriteString(s)
	}
	return builder.String()
}

func GetDataByInterface(field string, data map[string]interface{}) string {
	i, ok := data[field]
	if !ok {
		return ""
	}
	str, _ := i.(string)
	return str
}
