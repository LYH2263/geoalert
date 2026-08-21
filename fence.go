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

		verts := cloneLatLngs(f.Vertices)
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
		ring := toGeom(f.Vertices)
		if _, ok := geom.RingAnchor(ring); !ok {
			return false
		}
		return geom.PointInPolygon(geom.Point{Lat: lat, Lng: lng}, ring)
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
