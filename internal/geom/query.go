package geom

// ContainsAny 点集中是否有点在多边形内。
func ContainsAny(ring []Point, pts []Point) bool {
	for _, p := range pts {
		if PointInPolygon(p, ring) {
			return true
		}
	}
	return false
}

// ContainsAll 点集是否全部在多边形内。
func ContainsAll(ring []Point, pts []Point) bool {
	for _, p := range pts {
		if !PointInPolygon(p, ring) {
			return false
		}
	}
	return true
}

// FilterInside 过滤落在多边形内的点。
func FilterInside(ring []Point, pts []Point) []Point {
	out := make([]Point, 0, len(pts))
	for _, p := range pts {
		if PointInPolygon(p, ring) {
			out = append(out, p)
		}
	}
	return out
}

// NearestVertex 最近顶点索引。
func NearestVertex(ring []Point, p Point) (int, float64) {
	if len(ring) == 0 {
		return -1, 0
	}
	best := 0
	bestD := Distance(ring[0], p)
	for i := 1; i < len(ring); i++ {
		d := Distance(ring[i], p)
		if d < bestD {
			best = i
			bestD = d
		}
	}
	return best, bestD
}

// Translate 整体平移（度）。
func Translate(ring []Point, dLat, dLng float64) []Point {
	out := make([]Point, len(ring))
	for i, p := range ring {
		out[i] = Point{Lat: p.Lat + dLat, Lng: p.Lng + dLng}
	}
	return out
}
