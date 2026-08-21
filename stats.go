package geoalert

// Stats 返回运行统计。
func (e *Engine) Stats() EngineStats {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.statsLocked()
}

func (e *Engine) statsLocked() EngineStats {
	return EngineStats{
		Fences:       len(e.fences),
		Objects:      len(e.tracks),
		Ingests:      e.ingests.Load(),
		Enters:       e.metrics.Enters(),
		Exits:        e.metrics.Exits(),
		Dwells:       e.metrics.Dwells(),
		AlertsQueued: len(e.alerts),
		Closed:       e.closed.Load(),
	}
}
