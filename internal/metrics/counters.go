package metrics

import "sync/atomic"

// Registry 计数器。
type Registry struct {
	fences  atomic.Int64
	objects atomic.Int64
	ingests atomic.Int64
	enters  atomic.Int64
	exits   atomic.Int64
	dwells  atomic.Int64
}

func NewRegistry() *Registry { return &Registry{} }

func (r *Registry) IncFences(n int64)  { r.fences.Add(n) }
func (r *Registry) IncObjects(n int64) { r.objects.Add(n) }
func (r *Registry) IncIngests(n int64) { r.ingests.Add(n) }
func (r *Registry) IncEnters(n int64)  { r.enters.Add(n) }
func (r *Registry) IncExits(n int64)   { r.exits.Add(n) }
func (r *Registry) IncDwells(n int64)  { r.dwells.Add(n) }

func (r *Registry) Fences() int64  { return r.fences.Load() }
func (r *Registry) Objects() int64 { return r.objects.Load() }
func (r *Registry) Ingests() int64 { return r.ingests.Load() }
func (r *Registry) Enters() int64  { return r.enters.Load() }
func (r *Registry) Exits() int64   { return r.exits.Load() }
func (r *Registry) Dwells() int64  { return r.dwells.Load() }

// Snapshot 全部计数。
func (r *Registry) Snapshot() map[string]int64 {
	return map[string]int64{
		"fences":  r.Fences(),
		"objects": r.Objects(),
		"ingests": r.Ingests(),
		"enters":  r.Enters(),
		"exits":   r.Exits(),
		"dwells":  r.Dwells(),
	}
}
