package geoalert

import "github.com/LYH2263/go-geoalert/internal/clone"

func (e *Engine) ensureTrackLocked(id string) *objectTrack {
	t, ok := e.tracks[id]
	if ok {
		return t
	}
	t = &objectTrack{
		id:     id,
		points: make([]TrackPoint, 0, 8),
		inside: make(map[string]bool),
	}
	e.tracks[id] = t
	e.metrics.IncObjects(1)
	return t
}

func (e *Engine) appendPointLocked(t *objectTrack, p TrackPoint) {
	cp := p
	cp.Meta = clone.StringMap(p.Meta)
	t.points = append(t.points, cp)
	if max := e.opts.MaxTrackPoints; max > 0 && len(t.points) > max {
		overflow := len(t.points) - max
		copy(t.points, t.points[overflow:])
		t.points = t.points[:max]
	}
}

// Track 返回对象轨迹拷贝。
func (e *Engine) Track(objectID string) ([]TrackPoint, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	t, ok := e.tracks[objectID]
	if !ok {
		return nil, false
	}
	out := make([]TrackPoint, len(t.points))
	for i, p := range t.points {
		out[i] = p
		out[i].Meta = clone.StringMap(p.Meta)
	}
	return out, true
}

// ListTracks 轨迹摘要。
func (e *Engine) ListTracks() []TrackView {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.listTracksLocked()
}

func (e *Engine) listTracksLocked() []TrackView {
	out := make([]TrackView, 0, len(e.tracks))
	for _, t := range e.tracks {
		v := TrackView{ObjectID: t.id, Points: len(t.points)}
		if n := len(t.points); n > 0 {
			v.Last = t.points[n-1]
			v.Last.Meta = clone.StringMap(v.Last.Meta)
		}
		for fid, inside := range t.inside {
			if inside {
				v.Inside = append(v.Inside, fid)
			}
		}
		out = append(out, v)
	}
	return out
}

func (e *Engine) clearTracksLocked() {

	for id := range e.tracks {
		delete(e.tracks, id)
	}
}
