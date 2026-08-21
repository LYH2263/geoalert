package geom

import "math"

// Simplify 道格拉斯-普克简化（epsilon 为度）。
func Simplify(ring []Point, epsilon float64) []Point {
	if len(ring) < 3 || epsilon <= 0 {
		out := make([]Point, len(ring))
		copy(out, ring)
		return out
	}
	return douglasPeucker(ring, epsilon)
}

func douglasPeucker(points []Point, eps float64) []Point {
	if len(points) < 3 {
		out := make([]Point, len(points))
		copy(out, points)
		return out
	}
	maxDist := -1.0
	idx := 0
	start, end := points[0], points[len(points)-1]
	for i := 1; i < len(points)-1; i++ {
		d := perpDist(points[i], start, end)
		if d > maxDist {
			maxDist = d
			idx = i
		}
	}
	if maxDist > eps {
		left := douglasPeucker(points[:idx+1], eps)
		right := douglasPeucker(points[idx:], eps)
		out := make([]Point, 0, len(left)+len(right)-1)
		out = append(out, left[:len(left)-1]...)
		out = append(out, right...)
		return out
	}
	return []Point{start, end}
}

func perpDist(p, a, b Point) float64 {
	if a.Equals(b, 1e-15) {
		return Distance(p, a) / 111320.0
	}
	num := math.Abs((b.Lng-a.Lng)*(a.Lat-p.Lat) - (a.Lng-p.Lng)*(b.Lat-a.Lat))
	den := math.Hypot(b.Lng-a.Lng, b.Lat-a.Lat)
	if den == 0 {
		return 0
	}
	return num / den
}

// Densify 沿边按最大间距插值（米）。
func Densify(ring []Point, maxStepM float64) []Point {
	if len(ring) < 2 || maxStepM <= 0 {
		out := make([]Point, len(ring))
		copy(out, ring)
		return out
	}
	pts := EnsureOpen(ring)
	out := make([]Point, 0, len(pts)*2)
	for i := 0; i < len(pts); i++ {
		a := pts[i]
		b := pts[(i+1)%len(pts)]
		out = append(out, a)
		d := Distance(a, b)
		if d <= maxStepM {
			continue
		}
		n := int(math.Ceil(d / maxStepM))
		brg := BearingDegrees(a, b)
		for k := 1; k < n; k++ {
			out = append(out, Destination(a, brg, d*float64(k)/float64(n)))
		}
	}
	return out
}
