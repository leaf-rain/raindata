package domain

import "sync"

type hook struct {
	before     map[string]func(msg string) string
	syncBefore map[string]func(msg string)
	after      map[string]func(msg string)
	lock       sync.Mutex
}

// todo:实现类似pipline对数据进行过滤&修改等操作
func (h *hook) RegisterBefore(id string, f func(msg string) string) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.before == nil {
		h.before = make(map[string]func(msg string) string)
	}
	h.before[id] = f
}
func (h *hook) CancelBefore(id string) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if len(h.before) == 0 {
		return
	}
	delete(h.before, id)
}
func (h *hook) RegisterBeforeSync(id string, f func(msg string)) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.syncBefore == nil {
		h.syncBefore = make(map[string]func(msg string))
	}
	h.syncBefore[id] = f
}
func (h *hook) CancelBeforeSync(id string) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if len(h.before) == 0 {
		return
	}
	delete(h.before, id)
}

func (h *hook) RegisterAfter(id string, f func(msg string)) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.after == nil {
		h.after = make(map[string]func(msg string))
	}
	h.after[id] = f
}
func (h *hook) CancelAfter(id string) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if len(h.after) == 0 {
		return
	}
	delete(h.after, id)
}
