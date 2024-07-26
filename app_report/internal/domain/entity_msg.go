package domain

import (
	"github.com/leaf-rain/fastjson"
	"github.com/leaf-rain/raindata/app_report/internal/infrastructure/consts"
	"github.com/leaf-rain/raindata/common/rsql"
	"strconv"
	"sync"
)

type msgPool struct {
	sync.Pool
}

var defaultMsgPool = &msgPool{
	sync.Pool{
		New: func() interface{} {
			return &EntityMsg{}
		},
	},
}

func (p *msgPool) GetItem() *EntityMsg {
	m := defaultMsgPool.Pool.Get().(*EntityMsg)
	m.id = 0
	m.base = ""
	m.hook = nil
	m.value = nil
	m.appid = 0
	m.event = ""
	m.pool = p
	return m
}

type EntityMsg struct {
	id    int64
	base  string
	hook  *hook
	value *fastjson.Value
	appid int64
	event string
	pool  *msgPool
}

func NewEntityMsg(id, appid int64, msg string) (*EntityMsg, error) {
	result := defaultMsgPool.GetItem()
	result.id = id
	result.base = msg
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

func (msg *EntityMsg) getKeys() []rsql.FieldInfo {
	var fields []rsql.FieldInfo
	msg.value.Range(func(key []byte, v *fastjson.Value) {
		fields = append(fields, rsql.FieldInfo{
			Key:  string(key),
			Type: v.Type(),
		})
	})
	return fields
}

func (msg *EntityMsg) getEvent() string {
	return msg.event
}

func (msg *EntityMsg) getString() string {
	return msg.value.String()
}

func (msg *EntityMsg) putFields(fields map[string]string) bool {
	newValue := fastjson.MustParse("{}")
	var ok bool
	result := true
	msg.value.Range(func(key []byte, v *fastjson.Value) {
		if !result {
			return
		}
		_, ok = fields[string(key)]
		if !ok {
			result = false
		}
		newValue.Set(fields[string(key)], v)
	})
	msg.value = newValue
	return result
}

func (msg *EntityMsg) put() {
	msg.pool.Put(msg)
}
