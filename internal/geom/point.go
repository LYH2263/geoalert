package geom

// Point 经纬度。
type Point struct {
	Lat float64
	Lng float64
}

// Valid 粗检纬度经度范围。
func (p Point) Valid() bool {
	return p.Lat >= -90 && p.Lat <= 90 && p.Lng >= -180 && p.Lng <= 180
}

// Equals 近似相等。
func (p Point) Equals(o Point, eps float64) bool {
	if eps <= 0 {
		eps = 1e-9
	}
	dlat := p.Lat - o.Lat
	dlng := p.Lng - o.Lng
	if dlat < 0 {
		dlat = -dlat
	}
	if dlng < 0 {
		dlng = -dlng
	}
	return dlat <= eps && dlng <= eps
}
