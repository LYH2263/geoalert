package geom

import "math"

const earthRadiusM = 6371000.0

// HaversineMeters 两点大圆距离（米）。
func HaversineMeters(lat1, lng1, lat2, lng2 float64) float64 {
	φ1 := lat1 * math.Pi / 180
	φ2 := lat2 * math.Pi / 180
	Δφ := (lat2 - lat1) * math.Pi / 180
	Δλ := (lng2 - lng1) * math.Pi / 180
	a := math.Sin(Δφ/2)*math.Sin(Δφ/2) + math.Cos(φ1)*math.Cos(φ2)*math.Sin(Δλ/2)*math.Sin(Δλ/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusM * c
}

// Distance 两点距离。
func Distance(a, b Point) float64 {
	return HaversineMeters(a.Lat, a.Lng, b.Lat, b.Lng)
}

// BearingDegrees 初始方位角（度）。
func BearingDegrees(a, b Point) float64 {
	φ1 := a.Lat * math.Pi / 180
	φ2 := b.Lat * math.Pi / 180
	Δλ := (b.Lng - a.Lng) * math.Pi / 180
	y := math.Sin(Δλ) * math.Cos(φ2)
	x := math.Cos(φ1)*math.Sin(φ2) - math.Sin(φ1)*math.Cos(φ2)*math.Cos(Δλ)
	θ := math.Atan2(y, x) * 180 / math.Pi
	if θ < 0 {
		θ += 360
	}
	return θ
}

// Destination 沿方位角前进 distM 米后的点。
func Destination(start Point, bearingDeg, distM float64) Point {
	δ := distM / earthRadiusM
	θ := bearingDeg * math.Pi / 180
	φ1 := start.Lat * math.Pi / 180
	λ1 := start.Lng * math.Pi / 180
	φ2 := math.Asin(math.Sin(φ1)*math.Cos(δ) + math.Cos(φ1)*math.Sin(δ)*math.Cos(θ))
	λ2 := λ1 + math.Atan2(math.Sin(θ)*math.Sin(δ)*math.Cos(φ1), math.Cos(δ)-math.Sin(φ1)*math.Sin(φ2))
	lng := λ2 * 180 / math.Pi
	for lng > 180 {
		lng -= 360
	}
	for lng < -180 {
		lng += 360
	}
	return Point{Lat: φ2 * 180 / math.Pi, Lng: lng}
}

// Midpoint 两点中点（球面近似）。
func Midpoint(a, b Point) Point {
	φ1 := a.Lat * math.Pi / 180
	λ1 := a.Lng * math.Pi / 180
	φ2 := b.Lat * math.Pi / 180
	Δλ := (b.Lng - a.Lng) * math.Pi / 180
	Bx := math.Cos(φ2) * math.Cos(Δλ)
	By := math.Cos(φ2) * math.Sin(Δλ)
	φ3 := math.Atan2(math.Sin(φ1)+math.Sin(φ2), math.Sqrt((math.Cos(φ1)+Bx)*(math.Cos(φ1)+Bx)+By*By))
	λ3 := λ1 + math.Atan2(By, math.Cos(φ1)+Bx)
	return Point{Lat: φ3 * 180 / math.Pi, Lng: λ3 * 180 / math.Pi}
}
