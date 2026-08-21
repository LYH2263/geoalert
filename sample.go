package geoalert

import (
	"github.com/LYH2263/go-geoalert/internal/geom"
)

// SampleSquare 以中心点生成边长约 sideM 的正方形围栏顶点（近似）。
func SampleSquare(center LatLng, sideM float64) []LatLng {
	half := sideM / 2
	n := geom.Destination(geom.Point{Lat: center.Lat, Lng: center.Lng}, 0, half)
	e := geom.Destination(geom.Point{Lat: center.Lat, Lng: center.Lng}, 90, half)
	s := geom.Destination(geom.Point{Lat: center.Lat, Lng: center.Lng}, 180, half)
	w := geom.Destination(geom.Point{Lat: center.Lat, Lng: center.Lng}, 270, half)
	// 用轴对齐近似：取四角
	ne := geom.Point{Lat: n.Lat, Lng: e.Lng}
	se := geom.Point{Lat: s.Lat, Lng: e.Lng}
	sw := geom.Point{Lat: s.Lat, Lng: w.Lng}
	nw := geom.Point{Lat: n.Lat, Lng: w.Lng}
	return []LatLng{
		{Lat: nw.Lat, Lng: nw.Lng},
		{Lat: ne.Lat, Lng: ne.Lng},
		{Lat: se.Lat, Lng: se.Lng},
		{Lat: sw.Lat, Lng: sw.Lng},
	}
}

// SampleCircleFence 圆形围栏快捷构造。
func SampleCircleFence(id string, center LatLng, radiusM float64) Fence {
	return Fence{
		ID:           id,
		Kind:         FenceCircle,
		Center:       center,
		RadiusM:      radiusM,
		AlertOnEnter: true,
		AlertOnExit:  true,
	}
}

// SamplePolygonFence 多边形围栏快捷构造。
func SamplePolygonFence(id string, verts []LatLng) Fence {
	return Fence{
		ID:           id,
		Kind:         FencePolygon,
		Vertices:     verts,
		AlertOnEnter: true,
		AlertOnExit:  true,
	}
}
