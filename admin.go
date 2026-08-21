package geoalert

import (
	"time"

	"github.com/LYH2263/go-geoalert/internal/clone"
)

// UpdateFence 更新已有围栏（ID 不变）；持久化失败则不改内存。
func (e *Engine) UpdateFence(f Fence) error {
	if e.closed.Load() {
		return ErrClosed
	}
	norm, err := normalizeFence(f)
	if err != nil {
		return err
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, ok := e.fences[norm.ID]; !ok {
		return ErrNotFound
	}
	if err := e.store.SaveFence(toPersistFence(norm)); err != nil {
		return wrapPersist(err)
	}
	e.fences[norm.ID] = norm
	e.active[norm.ID] = true
	if e.auditor != nil {
		_ = e.auditor.Write("update", norm.ID, norm.Kind.String())
	}
	return nil
}

// SetFenceActive 激活/停用围栏（停用后不再产生告警）。
func (e *Engine) SetFenceActive(id string, active bool) error {
	if e.closed.Load() {
		return ErrClosed
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, ok := e.fences[id]; !ok {
		return ErrNotFound
	}
	e.active[id] = active
	return nil
}

// ExportFences 导出全部围栏定义拷贝。
func (e *Engine) ExportFences() []Fence {
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := make([]Fence, 0, len(e.fences))
	for _, f := range e.fences {
		cp := f
		cp.Vertices = cloneLatLngs(f.Vertices)
		cp.Tags = clone.Strings(f.Tags)
		out = append(out, cp)
	}
	return out
}

// PurgeObject 删除对象轨迹与边沿/驻留状态。
func (e *Engine) PurgeObject(objectID string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.tracks, objectID)
	e.edge.ClearObject(objectID)
	for id := range e.fences {
		e.dwell.Exit(objectID, id)
	}
}

// Now 当前引擎时间；无 Clock 时返回 zero。
func (e *Engine) Now() time.Time {
	if e.opts.Clock == nil {
		return time.Time{}
	}
	return e.opts.Clock.Now()
}
