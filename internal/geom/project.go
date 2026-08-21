package geom

import "math"

// XY 简易墨卡托投影坐标（米），用于局部平面运算。
type XY struct {
	X, Y float64
}

// ToMercator 经纬度 → 墨卡托。
func ToMercator(p Point) XY {
	x := p.Lng * earthRadiusM * math.Pi / 180
	lat := p.Lat
	if lat > 85.05112878 {
		lat = 85.05112878
	}
	if lat < -85.05112878 {
		lat = -85.05112878
	}
	y := math.Log(math.Tan((90+lat)*math.Pi/360)) * earthRadiusM
	return XY{X: x, Y: y}
}

// FromMercator 墨卡托 → 经纬度。
func FromMercator(xy XY) Point {
	lng := xy.X / earthRadiusM * 180 / math.Pi
	lat := 180/math.Pi*2*math.Atan(math.Exp(xy.Y/earthRadiusM)) - 90
	return Point{Lat: lat, Lng: lng}
}

// RingToMercator 批量投影。
func RingToMercator(ring []Point) []XY {
	out := make([]XY, len(ring))
	for i, p := range ring {
		out[i] = ToMercator(p)
	}
	return out
}

// DistanceXY 平面距离。
func DistanceXY(a, b XY) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return math.Hypot(dx, dy)
}

// BufferRingMeters 沿法向粗略外扩（米→度近似后回写）。
func BufferRingMeters(ring []Point, meters float64) []Point {
	if len(ring) < 3 || meters == 0 {
		out := make([]Point, len(ring))
		copy(out, ring)
		return out
	}
	c := Centroid(ring)
	out := make([]Point, len(ring))
	for i, p := range ring {
		brg := BearingDegrees(c, p)
		d := Distance(c, p) + meters
		if d < 0 {
			d = 0
		}
		out[i] = Destination(c, brg, d)
	}
	return out
}
