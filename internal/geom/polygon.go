package geom

// PointInPolygon 射线法判定点是否在多边形内（含边界视为内）。
func PointInPolygon(p Point, ring []Point) bool {
	n := len(ring)
	if n < 3 {
		return false
	}
	if OnBoundary(p, ring) {
		return true
	}
	inside := false
	j := n - 1
	for i := 0; i < n; i++ {
		yi, yj := ring[i].Lat, ring[j].Lat
		xi, xj := ring[i].Lng, ring[j].Lng
		intersect := ((yi > p.Lat) != (yj > p.Lat)) &&
			(p.Lng < (xj-xi)*(p.Lat-yi)/(yj-yi+1e-30)+xi)
		if intersect {
			inside = !inside
		}
		j = i
	}
	return inside
}

// OnBoundary 点是否落在任一边上。
func OnBoundary(p Point, ring []Point) bool {
	n := len(ring)
	if n < 2 {
		return false
	}
	for i := 0; i < n; i++ {
		a := ring[i]
		b := ring[(i+1)%n]
		if pointOnSegment(p, a, b) {
			return true
		}
	}
	return false
}

func pointOnSegment(p, a, b Point) bool {
	cross := (p.Lng-a.Lng)*(b.Lat-a.Lat) - (p.Lat-a.Lat)*(b.Lng-a.Lng)
	if cross < -1e-9 || cross > 1e-9 {
		return false
	}
	dot := (p.Lng-a.Lng)*(b.Lng-a.Lng) + (p.Lat-a.Lat)*(b.Lat-a.Lat)
	if dot < 0 {
		return false
	}
	len2 := (b.Lng-a.Lng)*(b.Lng-a.Lng) + (b.Lat-a.Lat)*(b.Lat-a.Lat)
	return dot <= len2+1e-9
}

// CloseRing 若首尾不同则追加闭合点。
func CloseRing(ring []Point) []Point {
	if len(ring) == 0 {
		return ring
	}
	if ring[0].Equals(ring[len(ring)-1], 1e-12) {
		out := make([]Point, len(ring))
		copy(out, ring)
		return out
	}
	out := make([]Point, len(ring)+1)
	copy(out, ring)
	out[len(ring)] = ring[0]
	return out
}

// EnsureOpen 去掉闭合重复点。
func EnsureOpen(ring []Point) []Point {
	n := len(ring)
	if n >= 2 && ring[0].Equals(ring[n-1], 1e-12) {
		out := make([]Point, n-1)
		copy(out, ring[:n-1])
		return out
	}
	out := make([]Point, n)
	copy(out, ring)
	return out
}
