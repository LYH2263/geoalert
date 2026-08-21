package geoalert

import (
	"time"

	"github.com/LYH2263/go-geoalert/internal/clone"
	"github.com/LYH2263/go-geoalert/internal/persist"
)

func toPersistFence(f Fence) persist.FenceRecord {
	return persist.FenceRecord{
		ID:           f.ID,
		Name:         f.Name,
		Kind:         int(f.Kind),
		Vertices:     toPersistVerts(f.Vertices),
		CenterLat:    f.Center.Lat,
		CenterLng:    f.Center.Lng,
		RadiusM:      f.RadiusM,
		DwellAfterMs: f.DwellAfter.Milliseconds(),
		Tags:         clone.Strings(f.Tags),
		AlertOnEnter: f.AlertOnEnter,
		AlertOnExit:  f.AlertOnExit,
		AlertOnDwell: f.AlertOnDwell,
	}
}

func toPersistVerts(vs []LatLng) []persist.LatLng {
	out := make([]persist.LatLng, len(vs))
	for i, v := range vs {
		out[i] = persist.LatLng{Lat: v.Lat, Lng: v.Lng}
	}
	return out
}

func fromPersistFence(r persist.FenceRecord) Fence {
	f := Fence{
		ID:           r.ID,
		Name:         r.Name,
		Kind:         FenceKind(r.Kind),
		Center:       LatLng{Lat: r.CenterLat, Lng: r.CenterLng},
		RadiusM:      r.RadiusM,
		DwellAfter:   time.Duration(r.DwellAfterMs) * time.Millisecond,
		Tags:         clone.Strings(r.Tags),
		AlertOnEnter: r.AlertOnEnter,
		AlertOnExit:  r.AlertOnExit,
		AlertOnDwell: r.AlertOnDwell,
	}
	f.Vertices = make([]LatLng, len(r.Vertices))
	for i, v := range r.Vertices {
		f.Vertices[i] = LatLng{Lat: v.Lat, Lng: v.Lng}
	}
	return f
}

// LoadPersisted 从 Store 装载围栏（覆盖同 ID）。
func (e *Engine) LoadPersisted() error {
	if e.closed.Load() {
		return ErrClosed
	}
	recs, err := e.store.ListFences()
	if err != nil {
		return wrapPersist(err)
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, r := range recs {
		f := fromPersistFence(r)
		norm, err := normalizeFence(f)
		if err != nil {
			continue
		}
		e.fences[norm.ID] = norm
		e.active[norm.ID] = true
	}
	return nil
}
