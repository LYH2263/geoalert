package geoalert_test

import (
	"errors"
	"testing"
	"time"

	"github.com/LYH2263/go-geoalert"
	"github.com/LYH2263/go-geoalert/internal/clock"
	"github.com/LYH2263/go-geoalert/internal/persist"
)

func TestBug03_IngestAfterCloseNoPanic(t *testing.T) {
	clk := clock.NewFake(time.Date(2026, 8, 21, 10, 0, 0, 0, time.UTC))
	eng, err := geoalert.New(geoalert.Options{Clock: clk, Store: persist.NewMemory()})
	if err != nil {
		t.Fatal(err)
	}
	if err := eng.RegisterFence(geoalert.SampleCircleFence("c", geoalert.LatLng{Lat: 0, Lng: 0}, 1000)); err != nil {
		t.Fatal(err)
	}
	if err := eng.Close(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Ingest after Close panicked via edge path: %v", r)
		}
	}()
	err = eng.Ingest(geoalert.TrackPoint{ObjectID: "o", At: time.Now(), Lat: 0, Lng: 0})
	if !errors.Is(err, geoalert.ErrClosed) {
		t.Fatalf("want ErrClosed, got %v", err)
	}
}
