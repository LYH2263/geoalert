package geoalert

import (
	"errors"
	"testing"
	"time"

	"github.com/LYH2263/go-geoalert/internal/clock"
	"github.com/LYH2263/go-geoalert/internal/persist"
)

// newTestEngine 构造一个最小可用的引擎（内存 Store + Fake 时钟）。
func newTestEngine(t *testing.T) *Engine {
	t.Helper()
	eng, err := New(Options{
		Store: persist.NewMemory(),
		Clock: clock.NewFake(time.Unix(1_700_000_000, 0)),
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return eng
}

func validTriangle() []LatLng {
	return []LatLng{
		{Lat: 31.20, Lng: 121.40},
		{Lat: 31.22, Lng: 121.40},
		{Lat: 31.21, Lng: 121.42},
	}
}

// 回归 1：空 Vertices 的多边形必须在规范化阶段被拒收，不得进入 active 表。
func TestRegisterFence_RejectsEmptyPolygon(t *testing.T) {
	eng := newTestEngine(t)

	// 空 Vertices 切片 —— 运营误配的典型场景。
	err := eng.RegisterFence(Fence{
		ID:       "empty",
		Kind:     FencePolygon,
		Vertices: []LatLng{},
	})
	if !errors.Is(err, ErrInvalidFence) {
		t.Fatalf("want ErrInvalidFence for empty polygon, got %v", err)
	}

	// nil Vertices 同样必须拒收。
	err = eng.RegisterFence(Fence{
		ID:   "nil-verts",
		Kind: FencePolygon,
	})
	if !errors.Is(err, ErrInvalidFence) {
		t.Fatalf("want ErrInvalidFence for nil polygon, got %v", err)
	}

	// 顶点不足 3 个也要拒收。
	err = eng.RegisterFence(Fence{
		ID:   "two-verts",
		Kind: FencePolygon,
		Vertices: []LatLng{
			{Lat: 31.20, Lng: 121.40},
			{Lat: 31.22, Lng: 121.40},
		},
	})
	if !errors.Is(err, ErrInvalidFence) {
		t.Fatalf("want ErrInvalidFence for <3 vertices, got %v", err)
	}

	// 关键断言：失败不得留半残 active 表项。
	for _, f := range eng.ListFences() {
		if f.Active {
			t.Fatalf("ListFences leaked active fence after failed register: %+v", f)
		}
	}
	if fv, ok := eng.GetFence("empty"); ok {
		t.Fatalf("GetFence leaked half-resident empty fence: %+v", fv)
	}
}

// 回归 1b：即使直接构造半残 Fence 调用 contains，RingAnchor 也不得 panic。
// 这兜底“空环曾让 ring[0] 越界 panic”的根因。
func TestContains_DoesNotPanicOnEmptyRing(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("contains panicked on empty ring: %v", r)
		}
	}()

	emptyPoly := Fence{ID: "leak", Kind: FencePolygon, Vertices: []LatLng{}}
	if emptyPoly.contains(31.21, 121.41) {
		t.Fatal("empty polygon should not contain any point")
	}

	nilPoly := Fence{ID: "leak2", Kind: FencePolygon}
	if nilPoly.contains(31.21, 121.41) {
		t.Fatal("nil-vertices polygon should not contain any point")
	}
}

// 回归 3：持久化失败时不得在内存里留半残 active 表项。
func TestRegisterFence_NoResidentEntryOnPersistFailure(t *testing.T) {
	mem := persist.NewMemory()
	eng, err := New(Options{
		Store: mem,
		Clock: clock.NewFake(time.Unix(1_700_000_000, 0)),
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	// 让后续 SaveFence 失败。
	mem.SetFail(true)

	err = eng.RegisterFence(Fence{
		ID:       "broken-persist",
		Kind:     FencePolygon,
		Vertices: validTriangle(),
	})
	if !errors.Is(err, ErrPersist) {
		t.Fatalf("want ErrPersist, got %v", err)
	}

	// 失败后，内存表不得挂这条 active 围栏。
	if fv, ok := eng.GetFence("broken-persist"); ok {
		t.Fatalf("GetFence leaked half-resident fence after persist failure: %+v", fv)
	}
	for _, f := range eng.ListFences() {
		if f.ID == "broken-persist" {
			t.Fatalf("ListFences leaked half-resident fence: %+v", f)
		}
	}

	// 恢复存储后可正常注册同 ID（证明没有残留占位）。
	mem.SetFail(false)
	if err := eng.RegisterFence(Fence{
		ID:       "broken-persist",
		Kind:     FencePolygon,
		Vertices: validTriangle(),
	}); err != nil {
		t.Fatalf("re-register after recovery: %v", err)
	}
}

// 回归 2：合法多边形注册后，点位 ingest 能正常评估，不会因 RingAnchor 崩溃。
func TestRegisterFence_IngestOnValidPolygonDoesNotCrash(t *testing.T) {
	eng := newTestEngine(t)

	if err := eng.RegisterFence(Fence{
		ID:       "tri",
		Kind:     FencePolygon,
		Vertices: validTriangle(),
	}); err != nil {
		t.Fatalf("register: %v", err)
	}

	// 点位在围栏内 —— 触发 contains -> RingAnchor 路径。
	if err := eng.Ingest(TrackPoint{
		ObjectID: "obj-1",
		At:       time.Unix(1_700_000_000, 0),
		Lat:      31.21,
		Lng:      121.41,
	}); err != nil {
		t.Fatalf("ingest: %v", err)
	}
}
