package geom

import "math"

// CrossTrackDistanceM 点到大圆弧的横向距离（近似，米）。
func CrossTrackDistanceM(p, a, b Point) float64 {
	d13 := HaversineMeters(a.Lat, a.Lng, p.Lat, p.Lng) / earthRadiusM
	θ13 := BearingDegrees(a, p) * math.Pi / 180
	θ12 := BearingDegrees(a, b) * math.Pi / 180
	return math.Asin(math.Sin(d13)*math.Sin(θ13-θ12)) * earthRadiusM
}

// AlongTrackDistanceM 点在 a→b 上的沿线距离。
func AlongTrackDistanceM(p, a, b Point) float64 {
	d13 := HaversineMeters(a.Lat, a.Lng, p.Lat, p.Lng) / earthRadiusM
	θ13 := BearingDegrees(a, p) * math.Pi / 180
	θ12 := BearingDegrees(a, b) * math.Pi / 180
	return math.Acos(math.Cos(d13)/math.Cos(math.Asin(math.Sin(d13)*math.Sin(θ13-θ12)))) * earthRadiusM
}

// NearestOnSegment 线段上最近点（球面粗近似，用平面墨卡托）。
func NearestOnSegment(p, a, b Point) Point {
	pa := ToMercator(p)
	aa := ToMercator(a)
	bb := ToMercator(b)
	dx := bb.X - aa.X
	dy := bb.Y - aa.Y
	len2 := dx*dx + dy*dy
	if len2 == 0 {
		return a
	}
	t := ((pa.X-aa.X)*dx + (pa.Y-aa.Y)*dy) / len2
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	return FromMercator(XY{X: aa.X + t*dx, Y: aa.Y + t*dy})
}

// DistanceToRingM 点到多边形边界最小距离。
func DistanceToRingM(p Point, ring []Point) float64 {
	pts := EnsureOpen(ring)
	n := len(pts)
	if n == 0 {
		return 0
	}
	best := Distance(p, pts[0])
	for i := 0; i < n; i++ {
		q := NearestOnSegment(p, pts[i], pts[(i+1)%n])
		d := Distance(p, q)
		if d < best {
			best = d
		}
	}
	return best
}
