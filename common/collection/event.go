package collection

import (
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/logx"
	"go.uber.org/zap"
	"sync"
	"time"
)

type DelayEvent interface {
	// RegisterFunc 注册方法
	RegisterFunc(fName string, f EventFunc)
	// HandleEvent 执行事件
	HandleEvent(key, value interface{})
	// StorgeEvent 存储事件
	StorgeEvent(key interface{}, value EventArg, interval int64)
	// DelEvent 删除事件
	DelEvent(key interface{})
	// FindOneEvent 获取事件
	FindOneEvent(key interface{}) *EventArg
}

var defaultDivisor int64 = 1

//var defaultDivisor int64 = 2

type EventFunc func(key, value interface{})

type EventArg struct {
	FName     string      // 方法名
	Surplus   int64       // 剩余执行次数, == -1时，一直循环调用
	RunTimes  int64       // 已经执行次数
	Interval  int64       // 间隔
	CreteTime int64       // 创建时间
	LastTime  int64       // 最后一次执行事件
	Args      interface{} // 函数参数
}

type delayEvent struct {
	eventTimingWheel *collection.TimingWheel
	eventMap         *sync.Map // 事件ma ，key:事件标识(string), value: EventArg
	funcMap          *sync.Map // 处理方法map key: 方法名(string)， value: EventFunc
	logger           *zap.Logger
}

func NewDelayEvent(logger *zap.Logger) *delayEvent {
	var e = &delayEvent{
		eventTimingWheel: nil,
		eventMap:         new(sync.Map),
		funcMap:          new(sync.Map),
		logger:           logger,
	}
	// 创建公共时间轮
	eventTimingWheel, err := collection.NewTimingWheel(time.Second/10, 1000, e.HandleEvent)
	if err != nil {
		logx.Errorf("[NewGameServer]collection.NewTimingWheel failed, err:%v", err)
		return nil
	}
	e.eventTimingWheel = eventTimingWheel
	return e
}

// RegisterFunc 注册方法
func (ev *delayEvent) RegisterFunc(fName string, f EventFunc) {
	logx.Infof("[RegisterFunc] RegisterFunc  fName:%v", fName)
	if f != nil {
		ev.funcMap.Store(fName, f)
	}
}

// HandleEvent 执行事件
func (ev *delayEvent) HandleEvent(key, value interface{}) {
	arg, ok := value.(EventArg)
	if !ok {
		logx.Errorf("[delayEvent] handleEvent arg type failed, key:%v, value:%v", key, value)
		return
	}
	var fI interface{}
	fI, ok = ev.funcMap.Load(arg.FName)
	if !ok {
		logx.Errorf("[delayEvent] handle func not found, key:%v, value:%v, arg:%+v", key, value, arg)
		return
	}
	var fc EventFunc
	fc, ok = fI.(EventFunc)
	if !ok {
		logx.Errorf("[delayEvent] EventFunc type failed, key:%v, value:%v, arg:%+v", key, value, arg)
		return
	}
	arg.RunTimes += 1
	arg.LastTime = time.Now().UnixNano()
	if arg.Surplus-1 > 0 || arg.Surplus == -1 {
		if arg.Surplus != -1 {
			arg.Surplus = arg.Surplus - 1
		}
		ev.StorgeEvent(key, arg, arg.Interval)
	}
	fc(key, arg.Args)
}

// StorgeEvent 存储事件
func (ev *delayEvent) StorgeEvent(key interface{}, value EventArg, interval int64) {
	logx.Debugf("[StorgeEvent] StorgeEvent  value:%+v, interval:%d", value, interval)
	if key == "" || value.Interval == 0 || value.Surplus == 0 || value.Surplus < -1 {
		logx.Errorf("[delayEvent] AddEvent Interval or Count is zero， key:%s, value:%+v", key, value)
		return
	}
	value.CreteTime = time.Now().Unix()
	_, ok := ev.funcMap.Load(value.FName)
	if !ok { // 添加未注册方法事件
		logx.Errorf("[delayEvent] AddEvent func name Not registered， key:%s, value:%+v", key, value)
		return
	}
	if interval == 0 { // 这里是为了可以支持首次可以设置不同的事件间隔
		interval = value.Interval
	}
	err := ev.eventTimingWheel.SetTimer(key, value, time.Second/time.Duration(defaultDivisor)*time.Duration(interval))
	if err != nil {
		logx.Errorf("[delayEvent] SetTimer failed，err:%v, key:%s, value:%+v", err, key, value)
		return
	}
	ev.eventMap.Store(key, value)
}

// DelEvent 删除事件
func (ev *delayEvent) DelEvent(key interface{}) {
	logx.Debugf("[DelEvent] DelEvent  fName:%+v", key)
	err := ev.eventTimingWheel.RemoveTimer(key)
	if err != nil {
		logx.Errorf("[DelEvent] RemoveTimer failed，err:%v, key:%s", err, key)
	}
	ev.eventMap.Delete(key)
}

// FindOneEvent 获取事件
func (ev *delayEvent) FindOneEvent(key interface{}) *EventArg {
	data, ok := ev.eventMap.Load(key)
	if !ok {
		logx.Errorf("[FindOneEvent] key not found， key:%s", key)
		return nil
	}
	var ret EventArg
	ret, ok = data.(EventArg)
	if !ok {
		logx.Errorf("[FindOneEvent] EventArg type failed， key:%s, value:%v", key, ret)
		return nil
	}
	return &ret
}

func (arg *EventArg) GetNextTime() int64 {
	if arg == nil {
		return 0
	}
	return maxInt64(arg.CreteTime, arg.LastTime) + (arg.Interval / defaultDivisor) - time.Now().Unix()
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
