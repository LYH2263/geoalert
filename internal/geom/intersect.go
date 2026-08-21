package geom

// SelfIntersects 检测简单多边形是否自交（相邻边共享顶点不算）。
func SelfIntersects(ring []Point) bool {
	n := len(ring)
	if n < 4 {
		return false
	}
	pts := EnsureOpen(ring)
	n = len(pts)
	for i := 0; i < n; i++ {
		a1 := pts[i]
		a2 := pts[(i+1)%n]
		for j := i + 1; j < n; j++ {
			b1 := pts[j]
			b2 := pts[(j+1)%n]
			if i == j || (i+1)%n == j || i == (j+1)%n {
				continue
			}
			if segmentsIntersectProper(a1, a2, b1, b2) {
				return true
			}
		}
	}
	return false
}

func segmentsIntersectProper(a1, a2, b1, b2 Point) bool {
	d1 := direction(b1, b2, a1)
	d2 := direction(b1, b2, a2)
	d3 := direction(a1, a2, b1)
	d4 := direction(a1, a2, b2)
	if ((d1 > 0 && d2 < 0) || (d1 < 0 && d2 > 0)) && ((d3 > 0 && d4 < 0) || (d3 < 0 && d4 > 0)) {
		return true
	}
	return false
}

func direction(a, b, c Point) float64 {
	return (c.Lng-a.Lng)*(b.Lat-a.Lat) - (c.Lat-a.Lat)*(b.Lng-a.Lng)
}

// SegmentIntersect 两线段是否相交（含端点）。
func SegmentIntersect(a1, a2, b1, b2 Point) bool {
	if segmentsIntersectProper(a1, a2, b1, b2) {
		return true
	}
	return pointOnSegment(a1, b1, b2) || pointOnSegment(a2, b1, b2) ||
		pointOnSegment(b1, a1, a2) || pointOnSegment(b2, a1, a2)
}
