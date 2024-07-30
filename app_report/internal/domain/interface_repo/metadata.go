package interface_repo

import (
	"github.com/leaf-rain/raindata/common/consts"
	"github.com/leaf-rain/raindata/common/rsql"
	"strconv"
)

type InterfaceMetadataRepo interface {
	MetadataPut(appid int64, event string, keys []rsql.FieldInfo) (map[string]string, error)
}

type DefaultMetadata struct{}

func NewMetadata() *DefaultMetadata {
	return &DefaultMetadata{}
}

func (d DefaultMetadata) MetadataPut(appid int64, event string, keys []rsql.FieldInfo) (map[string]string, error) {
	var result = make(map[string]string)
	for i := range keys {
		switch keys[i].Key {
		case consts.KeyAppidForMsg, consts.KeyEventForMsg, consts.KeyIdForMsg, consts.KeyVersionForMsg, consts.KeyCreateTimeForMsg:
			result[keys[i].Key] = keys[i].Key
		default:
			result[keys[i].Key] = keys[i].Key + "_" + strconv.FormatInt(int64(keys[i].Type), 10)
		}
	}
	return result, nil
}
