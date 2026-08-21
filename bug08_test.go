package geoalert_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/LYH2263/go-geoalert"
	"github.com/LYH2263/go-geoalert/internal/clock"
	"github.com/LYH2263/go-geoalert/internal/persist"
)

func TestBug08_WaitDwellHonorsCancel(t *testing.T) {
	clk := clock.NewFake(time.Date(2026, 8, 21, 10, 0, 0, 0, time.UTC))
	eng, err := geoalert.New(geoalert.Options{Clock: clk, Store: persist.NewMemory()})
	if err != nil {
		t.Fatal(err)
	}
	defer eng.Close()
	f := geoalert.SampleCircleFence("yard", geoalert.LatLng{Lat: 31.2, Lng: 121.5}, 200)
	f.AlertOnDwell = true
	f.DwellAfter = 3 * time.Second
	if err := eng.RegisterFence(f); err != nil {
		t.Fatal(err)
	}
	if err := eng.Ingest(geoalert.TrackPoint{ObjectID: "bus", At: clk.Now(), Lat: 31.2, Lng: 121.5}); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	start := time.Now()
	err = eng.WaitDwell(ctx, "bus", "yard")
	elapsed := time.Since(start)
	if elapsed > time.Second {
		t.Fatalf("WaitDwell ignored cancel, blocked %v", elapsed)
	}
	if err == nil || (!errors.Is(err, context.Canceled) && !errors.Is(err, geoalert.ErrCanceled)) {
		t.Fatalf("want canceled, got %v", err)
	}
}
