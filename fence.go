package geoalert

import (
	"github.com/LYH2263/go-geoalert/internal/clone"
	"github.com/LYH2263/go-geoalert/internal/geom"
	"github.com/LYH2263/go-geoalert/internal/validate"
)

// normalizeFence 校验并拷贝围栏几何，避免调用方切片别名污染内部状态。
func normalizeFence(f Fence) (Fence, error) {
	if err := validate.FenceID(f.ID); err != nil {
		return Fence{}, wrapInvalidFence(err)
	}
	out := f
	out.Tags = clone.Strings(f.Tags)
	switch f.Kind {
	case FencePolygon:
		if len(f.Vertices) < 3 {
			return Fence{}, wrapInvalidFence(validate.ErrTooFewVertices)
		}
		verts := cloneLatLngs(f.Vertices)
		lats := make([]float64, len(verts))
		lngs := make([]float64, len(verts))
		for i, v := range verts {
			lats[i], lngs[i] = v.Lat, v.Lng
		}
		if err := validate.CheckPolygon(lats, lngs); err != nil {
			return Fence{}, wrapInvalidFence(err)
		}
		pts := toGeom(verts)
		if geom.SelfIntersects(pts) {
			return Fence{}, wrapInvalidFence(validate.ErrSelfIntersect)
		}
		out.Vertices = verts
		out.Center = LatLng{}
		out.RadiusM = 0
	case FenceCircle:
		if err := validate.LatLng(f.Center.Lat, f.Center.Lng); err != nil {
			return Fence{}, wrapInvalidFence(err)
		}
		if f.RadiusM <= 0 {
			return Fence{}, wrapInvalidFence(validate.ErrBadRadius)
		}
		out.Vertices = nil
		out.Center = f.Center
		out.RadiusM = f.RadiusM
	default:
		return Fence{}, wrapInvalidFence(validate.ErrBadKind)
	}
	if !out.AlertOnEnter && !out.AlertOnExit && !out.AlertOnDwell {
		out.AlertOnEnter = true
		out.AlertOnExit = true
	}
	return out, nil
}

// cloneLatLngs 深拷贝顶点切片。LatLng 为纯值类型（仅两个 float64），
// 故 make+copy 即等价于深拷贝；nil 入参返回 nil，避免把空值替换成空切片。
func cloneLatLngs(src []LatLng) []LatLng {
	if src == nil {
		return nil
	}
	dst := make([]LatLng, len(src))
	copy(dst, src)
	return dst
}

func toGeom(vs []LatLng) []geom.Point {
	out := make([]geom.Point, len(vs))
	for i, v := range vs {
		out[i] = geom.Point{Lat: v.Lat, Lng: v.Lng}
	}
	return out
}

func (f Fence) contains(lat, lng float64) bool {
	switch f.Kind {
	case FencePolygon:
		return geom.PointInPolygon(geom.Point{Lat: lat, Lng: lng}, toGeom(f.Vertices))
	case FenceCircle:
		d := geom.HaversineMeters(f.Center.Lat, f.Center.Lng, lat, lng)
		return d <= f.RadiusM
	default:
		return false
	}
}

func fenceView(f Fence, active bool) FenceView {
	return FenceView{
		ID:         f.ID,
		Name:       f.Name,
		Kind:       f.Kind,
		Vertices:   cloneLatLngs(f.Vertices),
		Center:     f.Center,
		RadiusM:    f.RadiusM,
		DwellAfter: f.DwellAfter,
		Tags:       clone.Strings(f.Tags),
		Active:     active,
	}
}
