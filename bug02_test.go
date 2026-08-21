package geoalert_test

import (
	"testing"
	"time"

	"github.com/LYH2263/go-geoalert"
	"github.com/LYH2263/go-geoalert/internal/clock"
	"github.com/LYH2263/go-geoalert/internal/persist"
)

func TestBug02_AlertsSliceAlias(t *testing.T) {
	clk := clock.NewFake(time.Date(2026, 8, 21, 10, 0, 0, 0, time.UTC))
	eng, err := geoalert.New(geoalert.Options{Clock: clk, Store: persist.NewMemory()})
	if err != nil {
		t.Fatal(err)
	}
	defer eng.Close()
	f := geoalert.SampleCircleFence("c", geoalert.LatLng{Lat: 31.2, Lng: 121.5}, 500)
	if err := eng.RegisterFence(f); err != nil {
		t.Fatal(err)
	}
	if err := eng.Ingest(geoalert.TrackPoint{ObjectID: "truck", At: clk.Now(), Lat: 31.2, Lng: 121.5}); err != nil {
		t.Fatal(err)
	}
	a := eng.Alerts()
	if len(a) == 0 {
		t.Fatal("expected alert")
	}
	a[0].ObjectID = "mutated"
	b := eng.Alerts()
	if b[0].ObjectID == "mutated" {
		t.Fatal("Alerts shared backing array with caller")
	}
}
