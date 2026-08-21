package geom

// PointInCircle 点是否在圆内（含边界）。
func PointInCircle(p, center Point, radiusM float64) bool {
	if radiusM < 0 {
		return false
	}
	return HaversineMeters(center.Lat, center.Lng, p.Lat, p.Lng) <= radiusM
}

// CircleContainsRing 圆是否完全覆盖多边形顶点。
func CircleContainsRing(center Point, radiusM float64, ring []Point) bool {
	for _, p := range ring {
		if !PointInCircle(p, center, radiusM) {
			return false
		}
	}
	return true
}

// NearestOnCircle 圆上最近点（径向投影）。
func NearestOnCircle(p, center Point, radiusM float64) Point {
	d := HaversineMeters(center.Lat, center.Lng, p.Lat, p.Lng)
	if d == 0 {
		return Destination(center, 0, radiusM)
	}
	brg := BearingDegrees(center, p)
	return Destination(center, brg, radiusM)
}
