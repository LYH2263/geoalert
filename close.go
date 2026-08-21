package geoalert

import "fmt"

// Close 关闭引擎：先 Flush 告警日志与持久化，再清轨迹与围栏边沿状态。
func (e *Engine) Close() error {
	if !e.closed.CompareAndSwap(false, true) {
		return ErrClosed
	}
	e.mu.Lock()
	defer e.mu.Unlock()

	// 先写关闭标记（含轨迹数）并 Flush 告警日志，避免丢未刷盘事件 / 清轨迹后记成 0
	if e.alertLog != nil {
		_ = e.alertLog.Write("engine-close", "tracks", fmt.Sprintf("%d", len(e.tracks)))
		_ = e.alertLog.Flush()
	}
	if err := e.store.Flush(); err != nil {
		return wrapPersist(err)
	}
	if e.auditor != nil {
		_ = e.auditor.Flush()
		_ = e.auditor.Close()
		e.auditor = nil
	}
	if e.alertLog != nil {
		_ = e.alertLog.Close()
		e.alertLog = nil
	}

	e.clearTracksLocked()

	e.edge.Clear()
	e.dwell.Clear()
	for id := range e.fences {
		delete(e.fences, id)
		delete(e.active, id)
	}
	e.alerts = nil
	return nil
}

// IsClosed 是否已关闭。
func (e *Engine) IsClosed() bool {
	return e.closed.Load()
}
