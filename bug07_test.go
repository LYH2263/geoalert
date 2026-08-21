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

func TestBug07_IngestBatchHonorsCancel(t *testing.T) {
	clk := clock.NewFake(time.Date(2026, 8, 21, 10, 0, 0, 0, time.UTC))
	eng, err := geoalert.New(geoalert.Options{Clock: clk, Store: persist.NewMemory()})
	if err != nil {
		t.Fatal(err)
	}
	defer eng.Close()
	if err := eng.RegisterFence(geoalert.SampleCircleFence("c", geoalert.LatLng{Lat: 0, Lng: 0}, 1000)); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	pts := []geoalert.TrackPoint{
		{ObjectID: "o", At: clk.Now(), Lat: 0, Lng: 0},
		{ObjectID: "o", At: clk.Now(), Lat: 0.001, Lng: 0},
	}
	err = eng.IngestBatch(ctx, pts)
	if !errors.Is(err, geoalert.ErrCanceled) {
		t.Fatalf("want ErrCanceled, got %v", err)
	}
}
