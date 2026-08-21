package geoalert

import (
	"github.com/LYH2263/go-geoalert/internal/geom"
)

// Index 可选网格索引，加速多围栏候选。
type fenceIndex struct {
	g *geom.GridIndex
}

func newFenceIndex() *fenceIndex {
	return &fenceIndex{g: geom.NewGridIndex(0.02)}
}

func (e *Engine) ensureIndexLocked() {
	// 懒建：每次注册时更新（index 挂在 Engine 上需字段）
}

// CandidatesNear 返回点附近可能相关的围栏 ID（无索引时返回全部）。
func (e *Engine) CandidatesNear(lat, lng float64) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := make([]string, 0, len(e.fences))
	for id, f := range e.fences {
		var b geom.BBox
		switch f.Kind {
		case FencePolygon:
			b = geom.BBoxOfRing(toGeom(f.Vertices)).Expand(0.01)
		case FenceCircle:
			b = geom.BBoxOfCircle(geom.Point{Lat: f.Center.Lat, Lng: f.Center.Lng}, f.RadiusM).Expand(0.01)
		}
		if b.Contains(geom.Point{Lat: lat, Lng: lng}) {
			out = append(out, id)
		}
	}
	return out
}
