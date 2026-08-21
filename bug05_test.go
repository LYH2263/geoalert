package geoalert_test

import (
	"errors"
	"testing"
	"time"

	"github.com/LYH2263/go-geoalert"
	"github.com/LYH2263/go-geoalert/internal/clock"
	"github.com/LYH2263/go-geoalert/internal/persist"
)

func TestBug05_SelfIntersectWrapsInvalidFence(t *testing.T) {
	clk := clock.NewFake(time.Date(2026, 8, 21, 10, 0, 0, 0, time.UTC))
	eng, err := geoalert.New(geoalert.Options{Clock: clk, Store: persist.NewMemory()})
	if err != nil {
		t.Fatal(err)
	}
	defer eng.Close()
	verts := []geoalert.LatLng{
		{Lat: 0, Lng: 0}, {Lat: 1, Lng: 1}, {Lat: 0, Lng: 1}, {Lat: 1, Lng: 0},
	}
	err = eng.RegisterFence(geoalert.Fence{ID: "bowtie", Kind: geoalert.FencePolygon, Vertices: verts})
	if err == nil {
		t.Fatal("expected invalid self-intersect fence")
	}
	if !errors.Is(err, geoalert.ErrInvalidFence) {
		t.Fatalf("want errors.Is ErrInvalidFence, got %v", err)
	}
}
