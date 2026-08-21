package geom

import "math"

// RingAreaM2 球面多边形近似面积（平方米，正值）。
func RingAreaM2(ring []Point) float64 {
	pts := EnsureOpen(ring)
	n := len(pts)
	if n < 3 {
		return 0
	}
	var sum float64
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		sum += pts[i].Lng * pts[j].Lat
		sum -= pts[j].Lng * pts[i].Lat
	}
	if sum < 0 {
		sum = -sum
	}
	return sum * 0.5 * 111320 * 111320
}

// Centroid 平面质心（度坐标）。
func Centroid(ring []Point) Point {
	pts := EnsureOpen(ring)
	n := len(pts)
	if n == 0 {
		return Point{}
	}
	var a, cx, cy float64
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		cross := pts[i].Lng*pts[j].Lat - pts[j].Lng*pts[i].Lat
		a += cross
		cx += (pts[i].Lng + pts[j].Lng) * cross
		cy += (pts[i].Lat + pts[j].Lat) * cross
	}
	a *= 0.5
	if math.Abs(a) < 1e-18 {
		return pts[0]
	}
	return Point{Lat: cy / (6 * a), Lng: cx / (6 * a)}
}

// PerimeterMeters 周长。
func PerimeterMeters(ring []Point) float64 {
	pts := EnsureOpen(ring)
	n := len(pts)
	if n < 2 {
		return 0
	}
	var sum float64
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		sum += Distance(pts[i], pts[j])
	}
	return sum
}

// Winding 符号面积：正为逆时针。
func Winding(ring []Point) float64 {
	pts := EnsureOpen(ring)
	n := len(pts)
	var sum float64
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		sum += pts[i].Lng*pts[j].Lat - pts[j].Lng*pts[i].Lat
	}
	return sum / 2
}
