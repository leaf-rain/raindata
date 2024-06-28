package collection

import (
	"go.uber.org/zap"
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
	eventTimingWheel *TimingWheel
	eventMap         *SafeMap // 事件ma ，key:事件标识(string), value: EventArg
	funcMap          *SafeMap // 处理方法map key: 方法名(string)， value: EventFunc
	logger           *zap.Logger
}

func NewDelayEvent(logger *zap.Logger) *delayEvent {
	var e = &delayEvent{
		eventTimingWheel: nil,
		eventMap:         NewSafeMap(),
		funcMap:          NewSafeMap(),
		logger:           logger,
	}
	// 创建公共时间轮
	eventTimingWheel, err := NewTimingWheel(time.Second/10, 1000, e.HandleEvent)
	if err != nil {
		logger.Error("[NewGameServer]collection.NewTimingWheel failed", zap.Error(err))
		return nil
	}
	e.eventTimingWheel = eventTimingWheel
	return e
}

// RegisterFunc 注册方法
func (ev *delayEvent) RegisterFunc(fName string, f EventFunc) {
	ev.logger.Info("[RegisterFunc] RegisterFunc,", zap.String("funcName", fName))
	if f != nil {
		ev.funcMap.Set(fName, f)
	}
}

// HandleEvent 执行事件
func (ev *delayEvent) HandleEvent(key, value interface{}) {
	arg, ok := value.(EventArg)
	if !ok {
		ev.logger.Error("[delayEvent] handleEvent arg type failed.", zap.Any("key", key))
		return
	}
	var fI interface{}
	fI, ok = ev.funcMap.Get(arg.FName)
	if !ok {
		ev.logger.Error("[delayEvent] handle func not found.", zap.Any("key", key), zap.String("fName", arg.FName))
		return
	}
	var fc EventFunc
	fc, ok = fI.(EventFunc)
	if !ok {
		ev.logger.Error("[delayEvent] EventFunc type failed.", zap.Any("key", key), zap.String("fName", arg.FName))
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
	ev.logger.Debug("[StorgeEvent] StorgeEvent,", zap.Any("key", key), zap.Int64("interval", interval))
	if key == "" || value.Interval == 0 || value.Surplus == 0 || value.Surplus < -1 {
		ev.logger.Error("[delayEvent] AddEvent Interval or Count is zero.", zap.Any("key", key), zap.Int64("interval", interval))
		return
	}
	value.CreteTime = time.Now().Unix()
	_, ok := ev.funcMap.Get(value.FName)
	if !ok { // 添加未注册方法事件
		ev.logger.Error("[delayEvent] AddEvent func name Not registerer.", zap.Any("key", key), zap.Int64("interval", interval))
		return
	}
	if interval == 0 { // 这里是为了可以支持首次可以设置不同的事件间隔
		interval = value.Interval
	}
	err := ev.eventTimingWheel.SetTimer(key, value, time.Second/time.Duration(defaultDivisor)*time.Duration(interval))
	if err != nil {
		ev.logger.Error("[delayEvent] SetTimer failed.", zap.Error(err), zap.Any("key", key), zap.Int64("interval", interval))
		return
	}
	ev.eventMap.Set(key, value)
}

// DelEvent 删除事件
func (ev *delayEvent) DelEvent(key interface{}) {
	ev.logger.Debug("[DelEvent] DelEvent.", zap.Any("key", key))
	err := ev.eventTimingWheel.RemoveTimer(key)
	if err != nil {
		ev.logger.Error("[DelEvent] RemoveTimer failed.", zap.Error(err), zap.Any("key", key))
	}
	ev.eventMap.Del(key)
}

// FindOneEvent 获取事件
func (ev *delayEvent) FindOneEvent(key interface{}) *EventArg {
	data, ok := ev.eventMap.Get(key)
	if !ok {
		ev.logger.Error("[FindOneEvent] key not found.", zap.Any("key", key))
		return nil
	}
	var ret EventArg
	ret, ok = data.(EventArg)
	if !ok {
		ev.logger.Error("[FindOneEvent] EventArg type failed.", zap.Any("key", key))
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
