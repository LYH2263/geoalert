package geoalert

import "fmt"

// Close 关闭引擎：先 Flush 告警日志与持久化，再清轨迹与围栏边沿状态。
func (e *Engine) Close() error {
	if !e.closed.CompareAndSwap(false, true) {
		return ErrClosed
	}
	e.mu.Lock()
	defer e.mu.Unlock()

	// 先记录关闭时刻的轨迹数并 Flush 告警日志，再清轨迹；
	// 否则 tracks 已被清空，engine-close 行恒为 0。
	if e.alertLog != nil {
		_ = e.alertLog.Write("engine-close", "tracks", fmt.Sprintf("%d", len(e.tracks)))
		_ = e.alertLog.Flush()
	}
	e.clearTracksLocked()
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
