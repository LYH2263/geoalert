package geoalert

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/LYH2263/go-geoalert/internal/clock"
	"github.com/LYH2263/go-geoalert/internal/persist"
)

func testEng(t *testing.T) *Engine {
	t.Helper()
	clk := clock.NewFake(time.Date(2026, 8, 21, 10, 0, 0, 0, time.UTC))
	eng, err := New(Options{Clock: clk, Store: persist.NewMemory()})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = eng.Close() })
	return eng
}

func TestEnterExitPolygon(t *testing.T) {
	eng := testEng(t)
	verts := []LatLng{
		{Lat: 31.22, Lng: 121.47},
		{Lat: 31.22, Lng: 121.48},
		{Lat: 31.23, Lng: 121.48},
		{Lat: 31.23, Lng: 121.47},
	}
	if err := eng.RegisterFence(Fence{
		ID: "yard", Kind: FencePolygon, Vertices: verts,
		AlertOnEnter: true, AlertOnExit: true,
	}); err != nil {
		t.Fatal(err)
	}
	clk := eng.Clock().(*clock.Fake)
	outside := TrackPoint{ObjectID: "t1", At: clk.Now(), Lat: 31.21, Lng: 121.46}
	if err := eng.Ingest(outside); err != nil {
		t.Fatal(err)
	}
	clk.Advance(time.Second)
	inside := TrackPoint{ObjectID: "t1", At: clk.Now(), Lat: 31.225, Lng: 121.475}
	if err := eng.Ingest(inside); err != nil {
		t.Fatal(err)
	}
	alerts := eng.Alerts()
	if len(alerts) != 1 || alerts[0].Kind != AlertEnter {
		t.Fatalf("want enter, got %#v", alerts)
	}
	clk.Advance(time.Second)
	if err := eng.Ingest(TrackPoint{ObjectID: "t1", At: clk.Now(), Lat: 31.21, Lng: 121.46}); err != nil {
		t.Fatal(err)
	}
	alerts = eng.Alerts()
	if len(alerts) != 2 || alerts[1].Kind != AlertExit {
		t.Fatalf("want exit, got %#v", alerts)
	}
}

func TestRegisterFenceClonesVertices(t *testing.T) {
	eng := testEng(t)
	verts := []LatLng{
		{Lat: 0, Lng: 0}, {Lat: 0, Lng: 1}, {Lat: 1, Lng: 1}, {Lat: 1, Lng: 0},
	}
	if err := eng.RegisterFence(Fence{ID: "p", Kind: FencePolygon, Vertices: verts, AlertOnEnter: true}); err != nil {
		t.Fatal(err)
	}
	verts[0].Lat = 99
	v, ok := eng.GetFence("p")
	if !ok || v.Vertices[0].Lat == 99 {
		t.Fatalf("vertices aliased: %#v", v.Vertices)
	}
}

func TestAlertsReturnsCopy(t *testing.T) {
	eng := testEng(t)
	_ = eng.RegisterFence(SampleCircleFence("c", LatLng{Lat: 31.2, Lng: 121.5}, 500))
	clk := eng.Clock().(*clock.Fake)
	_ = eng.Ingest(TrackPoint{ObjectID: "o", At: clk.Now(), Lat: 31.2, Lng: 121.5})
	a := eng.Alerts()
	if len(a) == 0 {
		t.Fatal("expected alert")
	}
	a[0].ObjectID = "mutated"
	b := eng.Alerts()
	if b[0].ObjectID == "mutated" {
		t.Fatal("alerts shared backing array")
	}
}

func TestCloseThenIngest(t *testing.T) {
	eng := testEng(t)
	_ = eng.RegisterFence(SampleCircleFence("c", LatLng{Lat: 0, Lng: 0}, 1000))
	_ = eng.Close()
	err := eng.Ingest(TrackPoint{ObjectID: "o", At: time.Now(), Lat: 0, Lng: 0})
	if !errors.Is(err, ErrClosed) {
		t.Fatalf("want ErrClosed, got %v", err)
	}
}

func TestDwellRequiresClock(t *testing.T) {
	eng, err := New(Options{Store: persist.NewMemory()})
	if err != nil {
		t.Fatal(err)
	}
	defer eng.Close()
	f := SampleCircleFence("c", LatLng{Lat: 0, Lng: 0}, 100)
	f.AlertOnDwell = true
	f.DwellAfter = time.Second
	if err := eng.RegisterFence(f); err != nil {
		t.Fatal(err)
	}
	err = eng.Ingest(TrackPoint{ObjectID: "o", At: time.Now(), Lat: 0, Lng: 0})
	if err != nil {
		t.Fatal(err)
	}
	err = eng.Ingest(TrackPoint{ObjectID: "o", At: time.Now().Add(2 * time.Second), Lat: 0, Lng: 0})
	if !errors.Is(err, ErrNoClock) {
		t.Fatalf("want ErrNoClock, got %v", err)
	}
}

func TestSelfIntersectRejected(t *testing.T) {
	eng := testEng(t)
	// 蝴蝶结自交
	verts := []LatLng{
		{Lat: 0, Lng: 0}, {Lat: 1, Lng: 1}, {Lat: 0, Lng: 1}, {Lat: 1, Lng: 0},
	}
	err := eng.RegisterFence(Fence{ID: "x", Kind: FencePolygon, Vertices: verts})
	if !errors.Is(err, ErrInvalidFence) {
		t.Fatalf("want ErrInvalidFence, got %v", err)
	}
}

func TestPersistFailureNotActive(t *testing.T) {
	store := persist.NewMemory()
	store.SetFail(true)
	eng, err := New(Options{Clock: clock.Real{}, Store: store})
	if err != nil {
		t.Fatal(err)
	}
	defer eng.Close()
	err = eng.RegisterFence(SampleCircleFence("c", LatLng{Lat: 0, Lng: 0}, 10))
	if !errors.Is(err, ErrPersist) {
		t.Fatalf("want ErrPersist, got %v", err)
	}
	if len(eng.ListFences()) != 0 {
		t.Fatal("fence should not be active")
	}
}

func TestIngestBatchHonorsCancel(t *testing.T) {
	eng := testEng(t)
	_ = eng.RegisterFence(SampleCircleFence("c", LatLng{Lat: 0, Lng: 0}, 1000))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	pts := []TrackPoint{
		{ObjectID: "o", At: time.Now(), Lat: 0, Lng: 0},
		{ObjectID: "o", At: time.Now(), Lat: 0.001, Lng: 0},
	}
	err := eng.IngestBatch(ctx, pts)
	if !errors.Is(err, ErrCanceled) {
		t.Fatalf("want ErrCanceled, got %v", err)
	}
}

func TestDwellAlert(t *testing.T) {
	eng := testEng(t)
	f := SampleCircleFence("c", LatLng{Lat: 31.2, Lng: 121.5}, 200)
	f.AlertOnDwell = true
	f.DwellAfter = 5 * time.Second
	if err := eng.RegisterFence(f); err != nil {
		t.Fatal(err)
	}
	clk := eng.Clock().(*clock.Fake)
	_ = eng.Ingest(TrackPoint{ObjectID: "o", At: clk.Now(), Lat: 31.2, Lng: 121.5})
	clk.Advance(6 * time.Second)
	_ = eng.Ingest(TrackPoint{ObjectID: "o", At: clk.Now(), Lat: 31.2, Lng: 121.5})
	var found bool
	for _, a := range eng.Alerts() {
		if a.Kind == AlertDwell {
			found = true
		}
	}
	if !found {
		t.Fatalf("want dwell alert, got %#v", eng.Alerts())
	}
}
