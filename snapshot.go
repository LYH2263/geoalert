package geoalert

// Snapshot 引擎只读快照。
type Snapshot struct {
	Fences []FenceView
	Tracks []TrackView
	Alerts []Alert
	Stats  EngineStats
}

// Snapshot 生成拷贝快照。
func (e *Engine) Snapshot() Snapshot {
	e.mu.RLock()
	defer e.mu.RUnlock()
	s := Snapshot{
		Fences: make([]FenceView, 0, len(e.fences)),
		Tracks: e.listTracksLocked(),
		Alerts: cloneAlerts(e.alerts),
		Stats:  e.statsLocked(),
	}
	for _, f := range e.fences {
		s.Fences = append(s.Fences, fenceView(f, e.active[f.ID]))
	}
	return s
}
