package geoalert_test

import (
	"testing"
	"time"

	"github.com/LYH2263/go-geoalert"
	"github.com/LYH2263/go-geoalert/internal/clock"
	"github.com/LYH2263/go-geoalert/internal/persist"
)

func TestBug04_EmptyPolygonNoHalfFence(t *testing.T) {
	clk := clock.NewFake(time.Date(2026, 8, 21, 12, 0, 0, 0, time.UTC))
	eng, err := geoalert.New(geoalert.Options{Clock: clk, Store: persist.NewMemory()})
	if err != nil {
		t.Fatal(err)
	}
	defer eng.Close()
	defer func() {
		if rec := recover(); rec != nil {
			t.Fatalf("Ingest/Register path panicked on empty polygon: %v", rec)
		}
	}()
	err = eng.RegisterFence(geoalert.Fence{
		ID: "empty-poly", Kind: geoalert.FencePolygon, Vertices: nil, AlertOnEnter: true,
	})
	if err == nil {
		// 若错误地登记成功，Ingest 也必须安全，且最终不得留下半残围栏
		_ = eng.Ingest(geoalert.TrackPoint{ObjectID: "car", At: clk.Now(), Lat: 1, Lng: 1})
	}
	for _, f := range eng.ListFences() {
		if f.ID == "empty-poly" && f.Kind == geoalert.FencePolygon && len(f.Vertices) < 3 {
			t.Fatalf("half-broken empty polygon left in fence table: %+v", f)
		}
	}
	if _, ok := eng.GetFence("empty-poly"); ok {
		t.Fatal("empty polygon must not remain registered")
	}
}
