package policy

import "time"

// Limits 引擎策略上限。
type Limits struct {
	MaxFences      int
	MaxObjects     int
	MaxTrackPoints int
	MaxAlerts      int
	MinDwell       time.Duration
	MaxDwell       time.Duration
}

// DefaultLimits 缺省策略。
func DefaultLimits() Limits {
	return Limits{
		MaxFences:      1024,
		MaxObjects:     10000,
		MaxTrackPoints: 256,
		MaxAlerts:      1024,
		MinDwell:       0,
		MaxDwell:       24 * time.Hour,
	}
}

// ClampDwell 将驻留阈值夹到合法区间。
func (l Limits) ClampDwell(d time.Duration) time.Duration {
	if d < l.MinDwell {
		return l.MinDwell
	}
	if l.MaxDwell > 0 && d > l.MaxDwell {
		return l.MaxDwell
	}
	return d
}

// AllowFenceCount 是否允许再注册。
func (l Limits) AllowFenceCount(n int) bool {
	if l.MaxFences <= 0 {
		return true
	}
	return n < l.MaxFences
}

// AllowObjectCount 是否允许新对象。
func (l Limits) AllowObjectCount(n int) bool {
	if l.MaxObjects <= 0 {
		return true
	}
	return n < l.MaxObjects
}
