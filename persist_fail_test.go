package geoalert

import (
	"errors"
	"testing"
	"time"

	"github.com/LYH2263/go-geoalert/internal/clock"
	"github.com/LYH2263/go-geoalert/internal/persist"
)

// TestRegisterFence_PersistFailNoDirtyActive 回归：持久化失败时，内存 fences/active
// 不得被脏写激活，ListFences 看不到该围栏，Ingest 不会对它判进出。
func TestRegisterFence_PersistFailNoDirtyActive(t *testing.T) {
	mem := persist.NewMemory()
	mem.SetFail(true)
	e, err := New(Options{Store: mem, Clock: clock.NewFake(time.Now())})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	f := Fence{
		ID:           "f1",
		Name:         "fence1",
		Kind:         FenceCircle,
		Center:       LatLng{Lat: 30.0, Lng: 120.0},
		RadiusM:      1000,
		AlertOnEnter: true,
		AlertOnExit:  true,
	}
	if err := e.RegisterFence(f); !errors.Is(err, ErrPersist) {
		t.Fatalf("RegisterFence want ErrPersist, got %v", err)
	}

	// 内存层不得出现该围栏。
	if views := e.ListFences(); len(views) != 0 {
		t.Fatalf("ListFences after failed persist = %v, want empty", views)
	}
	if _, ok := e.GetFence("f1"); ok {
		t.Fatalf("GetFence after failed persist should be absent")
	}
	if s := e.Stats(); s.Fences != 0 {
		t.Fatalf("Stats.Fences = %d, want 0", s.Fences)
	}

	// 持久化层也不得残留该围栏（LoadPersisted 不应把它救活）。
	if recs, err := mem.ListFences(); err != nil || len(recs) != 0 {
		t.Fatalf("store.ListFences = %v (err=%v), want empty", recs, err)
	}
	if err := e.LoadPersisted(); err != nil {
		t.Fatalf("LoadPersisted: %v", err)
	}
	if views := e.ListFences(); len(views) != 0 {
		t.Fatalf("ListFences after LoadPersisted = %v, want empty", views)
	}

	// Ingest 不得产生任何告警。
	before := e.Alerts()
	_ = e.Ingest(TrackPoint{ObjectID: "obj1", At: time.Now(), Lat: 30.0, Lng: 120.0})
	after := e.Alerts()
	if len(before) != len(after) {
		t.Fatalf("alerts changed after ingest on failed-persist fence: before=%d after=%d", len(before), len(after))
	}
}

// TestRegisterFence_PersistOkActive 对照组：持久化成功后围栏可见且活跃。
func TestRegisterFence_PersistOkActive(t *testing.T) {
	mem := persist.NewMemory()
	e, err := New(Options{Store: mem, Clock: clock.NewFake(time.Now())})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	f := Fence{
		ID:           "f1",
		Kind:         FenceCircle,
		Center:       LatLng{Lat: 30.0, Lng: 120.0},
		RadiusM:      1000,
		AlertOnEnter: true,
		AlertOnExit:  true,
	}
	if err := e.RegisterFence(f); err != nil {
		t.Fatalf("RegisterFence: %v", err)
	}
	views := e.ListFences()
	if len(views) != 1 || !views[0].Active {
		t.Fatalf("ListFences = %+v, want 1 active", views)
	}
}
