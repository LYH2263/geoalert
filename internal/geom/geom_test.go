package geom_test

import (
	"testing"

	"github.com/LYH2263/go-geoalert/internal/geom"
)

func TestPointInPolygon(t *testing.T) {
	ring := []geom.Point{
		{Lat: 0, Lng: 0}, {Lat: 0, Lng: 2}, {Lat: 2, Lng: 2}, {Lat: 2, Lng: 0},
	}
	if !geom.PointInPolygon(geom.Point{Lat: 1, Lng: 1}, ring) {
		t.Fatal("center should be inside")
	}
	if geom.PointInPolygon(geom.Point{Lat: 3, Lng: 3}, ring) {
		t.Fatal("outside")
	}
}

func TestHaversine(t *testing.T) {
	// ~111 km per degree lat
	d := geom.HaversineMeters(0, 0, 1, 0)
	if d < 110000 || d > 112000 {
		t.Fatalf("unexpected distance %v", d)
	}
}

func TestSelfIntersect(t *testing.T) {
	bow := []geom.Point{
		{Lat: 0, Lng: 0}, {Lat: 1, Lng: 1}, {Lat: 0, Lng: 1}, {Lat: 1, Lng: 0},
	}
	if !geom.SelfIntersects(bow) {
		t.Fatal("expected self-intersect")
	}
	sq := []geom.Point{
		{Lat: 0, Lng: 0}, {Lat: 0, Lng: 1}, {Lat: 1, Lng: 1}, {Lat: 1, Lng: 0},
	}
	if geom.SelfIntersects(sq) {
		t.Fatal("square should not self-intersect")
	}
}

func TestPointInCircle(t *testing.T) {
	c := geom.Point{Lat: 31.2, Lng: 121.5}
	if !geom.PointInCircle(geom.Point{Lat: 31.2, Lng: 121.5}, c, 100) {
		t.Fatal("center")
	}
	far := geom.Destination(c, 90, 500)
	if geom.PointInCircle(far, c, 100) {
		t.Fatal("far point")
	}
}
