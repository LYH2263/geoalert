package geoalert

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/LYH2263/go-geoalert/internal/audit"
	"github.com/LYH2263/go-geoalert/internal/clock"
	"github.com/LYH2263/go-geoalert/internal/dwell"
	"github.com/LYH2263/go-geoalert/internal/edge"
	"github.com/LYH2263/go-geoalert/internal/metrics"
	"github.com/LYH2263/go-geoalert/internal/persist"
)

// Engine 地理围栏引擎。
type Engine struct {
	mu       sync.RWMutex
	opts     Options
	fences   map[string]Fence
	active   map[string]bool
	tracks   map[string]*objectTrack
	edge     *edge.Detector
	dwell    *dwell.Tracker
	alerts   []Alert
	seq      uint64
	closed   atomic.Bool
	metrics  *metrics.Registry
	store    persist.Store
	auditor  *audit.Logger
	alertLog *audit.Logger
	ingests  atomic.Int64
}

type objectTrack struct {
	id     string
	points []TrackPoint
	inside map[string]bool
}

// New 创建引擎；Clock 为空时 dwell 路径返回 ErrNoClock。
func New(opts Options) (*Engine, error) {
	opts = opts.withDefaults()
	e := &Engine{
		opts:    opts,
		fences:  make(map[string]Fence),
		active:  make(map[string]bool),
		tracks:  make(map[string]*objectTrack),
		edge:    edge.New(),
		dwell:   dwell.New(),
		metrics: metrics.NewRegistry(),
		store:   opts.Store,
	}
	if opts.Clock != nil {
		e.dwell.SetClock(opts.Clock)
	}
	if e.store == nil && opts.PersistPath != "" {
		st, err := persist.OpenFile(opts.PersistPath)
		if err != nil {
			return nil, err
		}
		e.store = st
	}
	if e.store == nil {
		e.store = persist.NewMemory()
	}
	if opts.Auditor != nil {
		e.auditor = opts.Auditor
	} else if opts.AuditPath != "" {
		_ = os.MkdirAll(filepath.Dir(opts.AuditPath), 0o755)
		a, err := audit.Open(opts.AuditPath)
		if err != nil {
			return nil, err
		}
		e.auditor = a
	}
	if opts.AlertLogPath != "" {
		_ = os.MkdirAll(filepath.Dir(opts.AlertLogPath), 0o755)
		al, err := audit.Open(opts.AlertLogPath)
		if err != nil {
			return nil, err
		}
		e.alertLog = al
	}
	return e, nil
}

// Clock 返回当前时钟（可为 nil）。
func (e *Engine) Clock() clock.Clock {
	return e.opts.Clock
}

// RegisterFence 注册围栏；持久化成功后才标记 active。Vertices 会被拷贝。
func (e *Engine) RegisterFence(f Fence) error {
	if e.closed.Load() {
		return ErrClosed
	}
	norm, err := normalizeFence(f)
	if err != nil {
		return err
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, ok := e.fences[norm.ID]; ok {
		return ErrExists
	}
	if err := e.store.SaveFence(toPersistFence(norm)); err != nil {
		return wrapPersist(err)
	}

	if f.Kind == FencePolygon {
		norm.Vertices = f.Vertices
	}
	e.fences[norm.ID] = norm
	e.active[norm.ID] = true
	e.metrics.IncFences(1)
	if e.auditor != nil {
		_ = e.auditor.Write("register", norm.ID, norm.Kind.String())
	}
	return nil
}

// UnregisterFence 移除围栏。
func (e *Engine) UnregisterFence(id string) error {
	if e.closed.Load() {
		return ErrClosed
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, ok := e.fences[id]; !ok {
		return ErrNotFound
	}
	delete(e.fences, id)
	delete(e.active, id)
	e.edge.ClearFence(id)
	e.dwell.ClearFence(id)
	_ = e.store.DeleteFence(id)
	e.metrics.IncFences(-1)
	if e.auditor != nil {
		_ = e.auditor.Write("unregister", id, "")
	}
	return nil
}

// GetFence 返回围栏快照。
func (e *Engine) GetFence(id string) (FenceView, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	f, ok := e.fences[id]
	if !ok {
		return FenceView{}, false
	}
	return fenceView(f, e.active[id]), true
}

// ListFences 列出全部围栏快照。
func (e *Engine) ListFences() []FenceView {
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := make([]FenceView, 0, len(e.fences))
	for _, f := range e.fences {
		out = append(out, fenceView(f, e.active[f.ID]))
	}
	return out
}

func (e *Engine) nextAlertID() string {
	e.seq++
	return fmt.Sprintf("a-%d", e.seq)
}
