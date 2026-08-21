package geoalert_test

import (
	"testing"
	"time"

	"github.com/LYH2263/go-geoalert"
	"github.com/LYH2263/go-geoalert/internal/clock"
	"github.com/LYH2263/go-geoalert/internal/persist"
)

func TestBug01_RegisterFenceVerticesAlias(t *testing.T) {
	clk := clock.NewFake(time.Date(2026, 8, 21, 10, 0, 0, 0, time.UTC))
	eng, err := geoalert.New(geoalert.Options{Clock: clk, Store: persist.NewMemory()})
	if err != nil {
		t.Fatal(err)
	}
	defer eng.Close()
	verts := []geoalert.LatLng{
		{Lat: 0, Lng: 0}, {Lat: 0, Lng: 1}, {Lat: 1, Lng: 1}, {Lat: 1, Lng: 0},
	}
	if err := eng.RegisterFence(geoalert.Fence{
		ID: "poly", Kind: geoalert.FencePolygon, Vertices: verts, AlertOnEnter: true,
	}); err != nil {
		t.Fatal(err)
	}
	verts[0].Lat = 99
	v, ok := eng.GetFence("poly")
	if !ok {
		t.Fatal("missing fence")
	}
	if v.Vertices[0].Lat == 99 {
		t.Fatal("RegisterFence aliased caller Vertices")
	}
}
