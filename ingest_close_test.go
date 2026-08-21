package geoalert

import (
	"errors"
	"testing"
	"time"

	"github.com/LYH2263/go-geoalert/internal/clock"
)

// 运维脚本 Close 后误点 Ingest：应返回 ErrClosed 而非在
// edge 边沿状态 nil map 上写崩。
func TestIngestAfterCloseReturnsErrClosed(t *testing.T) {
	e, err := New(Options{Clock: clock.NewFake(time.Now())})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := e.RegisterFence(SampleCircleFence("f1", LatLng{Lat: 30, Lng: 120}, 1000)); err != nil {
		t.Fatalf("register: %v", err)
	}

	if err := e.Ingest(TrackPoint{
		ObjectID: "o1", At: time.Now(),
		Lat: 30.0001, Lng: 120.0001,
	}); err != nil {
		t.Fatalf("ingest before close: %v", err)
	}

	if err := e.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	// 误点 Ingest：必须 guard，返回 ErrClosed，不得 panic。
	if err := e.Ingest(TrackPoint{
		ObjectID: "o1", At: time.Now(),
		Lat: 30.0001, Lng: 120.0001,
	}); !errors.Is(err, ErrClosed) {
		t.Fatalf("ingest after close = %v, want ErrClosed", err)
	}

	// 批量入口同样不得泄漏到边沿写入。
	if err := e.IngestBatch(nil, []TrackPoint{{
		ObjectID: "o1", At: time.Now(),
		Lat: 30.0001, Lng: 120.0001,
	}}); !errors.Is(err, ErrClosed) {
		t.Fatalf("ingest batch after close = %v, want ErrClosed", err)
	}
}
