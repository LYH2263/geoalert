package geoalert

func (e *Engine) emitLocked(a Alert) {
	a.ID = e.nextAlertID()
	a.Seq = e.seq
	e.alerts = append(e.alerts, a)
	if max := e.opts.MaxAlerts; max > 0 && len(e.alerts) > max {
		overflow := len(e.alerts) - max
		copy(e.alerts, e.alerts[overflow:])
		e.alerts = e.alerts[:max]
	}
	if e.alertLog != nil {
		_ = e.alertLog.Write(a.Kind.String(), a.ObjectID, a.FenceID)
	}
	if e.auditor != nil {
		_ = e.auditor.Write("alert:"+a.Kind.String(), a.ObjectID, a.FenceID)
	}
}

// cloneAlerts 深拷贝告警切片。Alert 为纯值类型结构体，make+copy 即完整拷贝，
// 返回独立底层数组，调用方修改不影响引擎内部队列。
func cloneAlerts(src []Alert) []Alert {
	if src == nil {
		return nil
	}
	dst := make([]Alert, len(src))
	copy(dst, src)
	return dst
}

// Alerts 返回告警队列的拷贝，调用方修改不影响内部。
func (e *Engine) Alerts() []Alert {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return cloneAlerts(e.alerts)
}

// DrainAlerts 取出并清空告警队列（返回拷贝）。
func (e *Engine) DrainAlerts() []Alert {
	e.mu.Lock()
	defer e.mu.Unlock()
	out := cloneAlerts(e.alerts)
	e.alerts = e.alerts[:0]
	return out
}

// ClearAlerts 清空告警队列。
func (e *Engine) ClearAlerts() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.alerts = e.alerts[:0]
}
