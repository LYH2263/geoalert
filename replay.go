package geoalert

import (
	"time"

	"github.com/LYH2263/go-geoalert/internal/geom"
)

// ReplayOptions 轨迹回放参数。
type ReplayOptions struct {
	Step time.Duration
}

// Replay 按时间顺序重放已缓存轨迹点（不新增存储，仅再评估边沿——用于诊断）。
// 实现：清空该对象边沿状态后按点再 Ingest。
func (e *Engine) Replay(objectID string, opts ReplayOptions) error {
	if e.closed.Load() {
		return ErrClosed
	}
	pts, ok := e.Track(objectID)
	if !ok || len(pts) == 0 {
		return ErrEmptyTrack
	}
	e.mu.Lock()
	e.edge.ClearObject(objectID)
	e.dwell.Exit(objectID, "") // no-op for empty fence; clear per fence below
	for id := range e.fences {
		e.dwell.Exit(objectID, id)
	}
	if t, ok := e.tracks[objectID]; ok {
		t.inside = make(map[string]bool)
		t.points = t.points[:0]
	}
	e.mu.Unlock()

	for _, p := range pts {
		if err := e.Ingest(p); err != nil {
			return err
		}
		if opts.Step > 0 && e.opts.Clock != nil {
			// 仅占位：真实等待由调用方控制
			_ = opts.Step
		}
	}
	return nil
}

// DistanceToFence 对象当前位置到围栏的粗略距离（米）；在内为 0。
func (e *Engine) DistanceToFence(objectID, fenceID string) (float64, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	f, ok := e.fences[fenceID]
	if !ok {
		return 0, ErrNotFound
	}
	t, ok := e.tracks[objectID]
	if !ok || len(t.points) == 0 {
		return 0, ErrEmptyTrack
	}
	p := t.points[len(t.points)-1]
	if f.contains(p.Lat, p.Lng) {
		return 0, nil
	}
	switch f.Kind {
	case FenceCircle:
		d := geom.HaversineMeters(f.Center.Lat, f.Center.Lng, p.Lat, p.Lng)
		if d < f.RadiusM {
			return 0, nil
		}
		return d - f.RadiusM, nil
	case FencePolygon:
		pts := toGeom(f.Vertices)
		_, d := geom.NearestVertex(pts, geom.Point{Lat: p.Lat, Lng: p.Lng})
		return d, nil
	default:
		return 0, ErrInvalidFence
	}
}
