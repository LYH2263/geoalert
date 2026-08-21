package geom

import "math"

// BBox 轴对齐包围盒（经纬度）。
type BBox struct {
	MinLat, MaxLat float64
	MinLng, MaxLng float64
}

// BBoxOfRing 计算多边形包围盒。
func BBoxOfRing(ring []Point) BBox {
	if len(ring) == 0 {
		return BBox{}
	}
	b := BBox{MinLat: ring[0].Lat, MaxLat: ring[0].Lat, MinLng: ring[0].Lng, MaxLng: ring[0].Lng}
	for _, p := range ring[1:] {
		if p.Lat < b.MinLat {
			b.MinLat = p.Lat
		}
		if p.Lat > b.MaxLat {
			b.MaxLat = p.Lat
		}
		if p.Lng < b.MinLng {
			b.MinLng = p.Lng
		}
		if p.Lng > b.MaxLng {
			b.MaxLng = p.Lng
		}
	}
	return b
}

// Contains 点是否在盒内。
func (b BBox) Contains(p Point) bool {
	return p.Lat >= b.MinLat && p.Lat <= b.MaxLat && p.Lng >= b.MinLng && p.Lng <= b.MaxLng
}

// Intersects 两包围盒是否相交。
func (b BBox) Intersects(o BBox) bool {
	return !(b.MaxLat < o.MinLat || b.MinLat > o.MaxLat || b.MaxLng < o.MinLng || b.MinLng > o.MaxLng)
}

// Expand 扩展边距（度）。
func (b BBox) Expand(deg float64) BBox {
	b.MinLat -= deg
	b.MaxLat += deg
	b.MinLng -= deg
	b.MaxLng += deg
	return b
}

// AreaDeg2 面积（度²）。
func (b BBox) AreaDeg2() float64 {
	return (b.MaxLat - b.MinLat) * (b.MaxLng - b.MinLng)
}

// BBoxOfCircle 粗略圆包围盒。
func BBoxOfCircle(center Point, radiusM float64) BBox {
	dlat := radiusM / 111320.0
	cl := math.Abs(center.Lat)
	cos := math.Cos(cl * math.Pi / 180)
	if cos < 0.2 {
		cos = 0.2
	}
	dlng := radiusM / (111320.0 * cos)
	return BBox{
		MinLat: center.Lat - dlat,
		MaxLat: center.Lat + dlat,
		MinLng: center.Lng - dlng,
		MaxLng: center.Lng + dlng,
	}
}
