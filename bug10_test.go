package geoalert_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/LYH2263/go-geoalert"
	"github.com/LYH2263/go-geoalert/internal/clock"
	"github.com/LYH2263/go-geoalert/internal/persist"
)

func TestBug10_CloseFlushesAlertLogBeforeClearTracks(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "alerts.jsonl")
	clk := clock.NewFake(time.Date(2026, 8, 21, 10, 0, 0, 0, time.UTC))
	eng, err := geoalert.New(geoalert.Options{
		Clock:        clk,
		Store:        persist.NewMemory(),
		AlertLogPath: logPath,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := eng.RegisterFence(geoalert.SampleCircleFence("c", geoalert.LatLng{Lat: 31.2, Lng: 121.5}, 500)); err != nil {
		t.Fatal(err)
	}
	if err := eng.Ingest(geoalert.TrackPoint{ObjectID: "van", At: clk.Now(), Lat: 31.2, Lng: 121.5}); err != nil {
		t.Fatal(err)
	}
	if err := eng.Close(); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	body := string(b)
	if !strings.Contains(body, "engine-close\ttracks\t1") {
		t.Fatalf("Close cleared tracks before flushing alert log; body=%q", body)
	}
}
