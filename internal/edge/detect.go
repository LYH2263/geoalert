package edge

import (
	"sync"
	"time"
)

// Detector 对象×围栏边沿检测器。
type Detector struct {
	mu    sync.Mutex
	state map[key]cell
}

func New() *Detector {
	return &Detector{state: make(map[key]cell)}
}

// Observe 记录一次在/外观察并返回边沿。
// 首次见到且在内 → Enter；首次在外 → None；之后同态 Stay，异态 Enter/Exit。
func (d *Detector) Observe(objectID, fenceID string, inside Presence, at time.Time) Transition {
	d.mu.Lock()
	defer d.mu.Unlock()
	k := key{Object: objectID, Fence: fenceID}
	cur := bool(inside)
	prev, ok := d.state[k]
	tr := Transition{ObjectID: objectID, FenceID: fenceID, At: at, Inside: cur}
	if !ok || !prev.seen {
		d.state[k] = cell{inside: cur, seen: true, at: at}
		if cur {
			tr.Kind = Enter
		} else {
			tr.Kind = None
		}
		return tr
	}
	if prev.inside == cur {
		tr.Kind = Stay
		d.state[k] = cell{inside: cur, seen: true, at: at}
		return tr
	}
	if cur {
		tr.Kind = Enter
	} else {
		tr.Kind = Exit
	}
	d.state[k] = cell{inside: cur, seen: true, at: at}
	return tr
}

// ClearFence 清除某围栏全部状态。
func (d *Detector) ClearFence(fenceID string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for k := range d.state {
		if k.Fence == fenceID {
			delete(d.state, k)
		}
	}
}

// ClearObject 清除某对象。
func (d *Detector) ClearObject(objectID string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for k := range d.state {
		if k.Object == objectID {
			delete(d.state, k)
		}
	}
}

// Clear 清空。
func (d *Detector) Clear() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.state = nil
}

// Inside 查询缓存状态。
func (d *Detector) Inside(objectID, fenceID string) (bool, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	c, ok := d.state[key{Object: objectID, Fence: fenceID}]
	if !ok || !c.seen {
		return false, false
	}
	return c.inside, true
}

// Snapshot 导出全部状态（测试用）。
func (d *Detector) Snapshot() map[string]bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	out := make(map[string]bool, len(d.state))
	for k, c := range d.state {
		if c.seen {
			out[k.Object+"|"+k.Fence] = c.inside
		}
	}
	return out
}
