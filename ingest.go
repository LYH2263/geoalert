package geoalert

import (
	"context"
	"errors"
	"time"

	"github.com/LYH2263/go-geoalert/internal/dwell"
	"github.com/LYH2263/go-geoalert/internal/edge"
	"github.com/LYH2263/go-geoalert/internal/validate"
)

// Ingest 摄入单点并评估全部活跃围栏。
func (e *Engine) Ingest(p TrackPoint) error {
	return e.IngestContext(context.Background(), p)
}

// IngestContext 摄入单点；ctx 在评估前检查。
func (e *Engine) IngestContext(ctx context.Context, p TrackPoint) error {
	if err := ctx.Err(); err != nil {
		return ErrCanceled
	}

	if err := validate.ObjectID(p.ObjectID); err != nil {
		return ErrInvalidObject
	}
	if err := validate.LatLng(p.Lat, p.Lng); err != nil {
		return ErrInvalidPoint
	}
	if p.At.IsZero() {
		if e.opts.Clock == nil {
			return ErrNoClock
		}
		p.At = e.opts.Clock.Now()
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	return e.ingestLocked(ctx, p)
}

// IngestBatch 批量摄入；每点之间检查 ctx。
func (e *Engine) IngestBatch(ctx context.Context, points []TrackPoint) error {
	if e.closed.Load() {
		return ErrClosed
	}
	for i := range points {
		if err := ctx.Err(); err != nil {
			return ErrCanceled
		}
		if err := e.IngestContext(ctx, points[i]); err != nil {
			return err
		}
	}
	return nil
}

func (e *Engine) ingestLocked(ctx context.Context, p TrackPoint) error {
	if err := ctx.Err(); err != nil {
		return ErrCanceled
	}
	t := e.ensureTrackLocked(p.ObjectID)
	e.appendPointLocked(t, p)
	e.ingests.Add(1)
	e.metrics.IncIngests(1)

	for id, f := range e.fences {
		if !e.active[id] {
			continue
		}
		if err := ctx.Err(); err != nil {
			return ErrCanceled
		}
		inside := f.contains(p.Lat, p.Lng)
		tr := e.edge.Observe(p.ObjectID, id, edge.Presence(inside), p.At)
		t.inside[id] = inside

		switch tr.Kind {
		case edge.Enter:
			if f.AlertOnEnter {
				e.emitLocked(Alert{
					FenceID:  id,
					ObjectID: p.ObjectID,
					Kind:     AlertEnter,
					At:       p.At,
					Lat:      p.Lat,
					Lng:      p.Lng,
				})
				e.metrics.IncEnters(1)
			}
			if f.AlertOnDwell || f.DwellAfter > 0 {
				th := f.DwellAfter
				if th <= 0 {
					th = e.opts.DefaultDwell
				}
				if th > 0 {
					e.dwell.Enter(p.ObjectID, id, p.At, th)
				}
			}
		case edge.Exit:
			if f.AlertOnExit {
				e.emitLocked(Alert{
					FenceID:  id,
					ObjectID: p.ObjectID,
					Kind:     AlertExit,
					At:       p.At,
					Lat:      p.Lat,
					Lng:      p.Lng,
				})
				e.metrics.IncExits(1)
			}
			e.dwell.Exit(p.ObjectID, id)
		case edge.Stay:
			if inside && (f.AlertOnDwell || f.DwellAfter > 0) {
				if err := e.checkDwellLocked(ctx, f, p); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (e *Engine) checkDwellLocked(ctx context.Context, f Fence, p TrackPoint) error {
	if e.opts.Clock == nil {
		return ErrNoClock
	}
	th := f.DwellAfter
	if th <= 0 {
		th = e.opts.DefaultDwell
	}
	if th <= 0 {
		return nil
	}
	ev, ok, err := e.dwell.Check(ctx, p.ObjectID, f.ID, p.At)
	if err != nil {
		if errors.Is(err, dwell.ErrNoClock) {
			return ErrNoClock
		}
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return ErrCanceled
		}
		return err
	}
	if !ok || !f.AlertOnDwell {
		return nil
	}
	e.emitLocked(Alert{
		FenceID:  f.ID,
		ObjectID: p.ObjectID,
		Kind:     AlertDwell,
		At:       p.At,
		Lat:      p.Lat,
		Lng:      p.Lng,
		DwellFor: ev.Duration,
	})
	e.metrics.IncDwells(1)
	return nil
}

// WaitDwell 等待对象在围栏内驻留达到阈值（尊重 ctx）。
func (e *Engine) WaitDwell(ctx context.Context, objectID, fenceID string) error {
	if e.closed.Load() {
		return ErrClosed
	}
	if e.opts.Clock == nil {
		return ErrNoClock
	}
	e.mu.RLock()
	f, ok := e.fences[fenceID]
	th := time.Duration(0)
	if ok {
		th = f.DwellAfter
		if th <= 0 {
			th = e.opts.DefaultDwell
		}
	}
	e.mu.RUnlock()
	if !ok {
		return ErrNotFound
	}
	if th <= 0 {
		return ErrDwellThreshold
	}
	return e.dwell.Wait(ctx, objectID, fenceID, th)
}
