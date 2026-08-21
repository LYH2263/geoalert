package geoalert

import (
	"sort"

	"github.com/LYH2263/go-geoalert/internal/geom"
)

// FenceBounds 返回围栏包围盒（度）。
func (e *Engine) FenceBounds(id string) (minLat, maxLat, minLng, maxLng float64, ok bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	f, found := e.fences[id]
	if !found {
		return 0, 0, 0, 0, false
	}
	var b geom.BBox
	switch f.Kind {
	case FencePolygon:
		b = geom.BBoxOfRing(toGeom(f.Vertices))
	case FenceCircle:
		b = geom.BBoxOfCircle(geom.Point{Lat: f.Center.Lat, Lng: f.Center.Lng}, f.RadiusM)
	default:
		return 0, 0, 0, 0, false
	}
	return b.MinLat, b.MaxLat, b.MinLng, b.MaxLng, true
}

// ObjectsInside 列出当前判定在指定围栏内的对象。
func (e *Engine) ObjectsInside(fenceID string) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	var out []string
	for _, t := range e.tracks {
		if t.inside[fenceID] {
			out = append(out, t.id)
		}
	}
	sort.Strings(out)
	return out
}

// EvaluatePoint 不改状态：仅判断点相对围栏内外。
func (e *Engine) EvaluatePoint(fenceID string, lat, lng float64) (bool, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	f, ok := e.fences[fenceID]
	if !ok {
		return false, ErrNotFound
	}
	return f.contains(lat, lng), nil
}
