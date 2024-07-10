package domain

import (
	"github.com/leaf-rain/fastjson"
	"github.com/leaf-rain/raindata/app_report/internal/infrastructure/consts"
	"strconv"
)

type EntityMsg struct {
	id    int64
	base  string
	hook  *hook
	value *fastjson.Value
	appid int64
	event string
}

func NewEntityMsg(id, appid int64, msg string) (*EntityMsg, error) {
	result := &EntityMsg{
		id:   id,
		base: msg,
	}
	var err error
	result.value, err = fastjson.Parse(msg)
	result.value.Del(consts.KeyAppidForMsg)
	// 添加功能默认字段
	result.appid = appid
	idValue := fastjson.MustParse(strconv.FormatInt(id, 10))
	result.value.Set(consts.KeyIdForMsg, idValue)
	result.event = string(result.value.GetStringBytes(consts.KeyEventForMsg))
	if result.event == "" {
		result.value.Set(consts.KeyEventForMsg, fastjson.MustParse("\"log\""))
		result.event = "log"
	}
	return result, err
}

func (msg *EntityMsg) getKeys() []string {
	var keys []string
	msg.value.Range(func(key []byte, v *fastjson.Value) {
		keys = append(keys, string(key))
	})
	return keys
}

func (msg *EntityMsg) getEvent() []string {
	var keys []string
	msg.value.Range(func(key []byte, v *fastjson.Value) {
		keys = append(keys, string(key))
	})
	return keys
}

func (msg *EntityMsg) getString() string {
	return msg.value.String()
}
