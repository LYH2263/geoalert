package geoalert_test

import (
	"errors"
	"testing"
	"time"

	"github.com/LYH2263/go-geoalert"
	"github.com/LYH2263/go-geoalert/internal/clock"
	"github.com/LYH2263/go-geoalert/internal/persist"
)

func TestBug06_PersistFailNotActivate(t *testing.T) {
	store := persist.NewMemory()
	store.SetFail(true)
	clk := clock.NewFake(time.Date(2026, 8, 21, 10, 0, 0, 0, time.UTC))
	eng, err := geoalert.New(geoalert.Options{Clock: clk, Store: store})
	if err != nil {
		t.Fatal(err)
	}
	defer eng.Close()
	err = eng.RegisterFence(geoalert.SampleCircleFence("c", geoalert.LatLng{Lat: 0, Lng: 0}, 10))
	if !errors.Is(err, geoalert.ErrPersist) {
		t.Fatalf("want ErrPersist, got %v", err)
	}
	if n := len(eng.ListFences()); n != 0 {
		t.Fatalf("fence activated despite persist failure: %d", n)
	}
}
