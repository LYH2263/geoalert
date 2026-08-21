package geoalert

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/LYH2263/go-geoalert/internal/clock"
)

// validCircle 返回一个合法的圆形围栏（上海某处，半径 1000m）。
func validCircle(id string) Fence {
	return Fence{
		ID:      id,
		Kind:    FenceCircle,
		Center:  LatLng{Lat: 31.2304, Lng: 121.4737},
		RadiusM: 1000,
	}
}

// TestIngestBatchRespectsCancel 防止 IngestBatch 退化为 context.Background()：
// 取消 ctx 后必须尽快返回 ErrCanceled，而非把剩余点逐个啃完。
func TestIngestBatchRespectsCancel(t *testing.T) {
	eng, err := New(Options{Clock: clock.NewFake(time.UnixMilli(1_700_000_000_000))})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := eng.RegisterFence(validCircle("yard")); err != nil {
		t.Fatalf("RegisterFence: %v", err)
	}

	pts := make([]TrackPoint, 100_000)
	for i := range pts {
		pts[i] = TrackPoint{ObjectID: "truck-1", At: time.UnixMilli(1_700_000_000_000), Lat: 31.2304, Lng: 121.4737}
	}

	ctx, cancel := context.WithCancel(context.Background())
	// 立即取消：逐点循环下一次 ctx.Err() 即应短路。
	cancel()

	done := make(chan error, 1)
	start := time.Now()
	go func() { done <- eng.IngestBatch(ctx, pts) }()
	select {
	case err := <-done:
		if !errors.Is(err, ErrCanceled) {
			t.Fatalf("期望 ErrCanceled，实际 %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("IngestBatch 取消后未在 2s 内返回（仍在啃剩余坐标点）")
	}

	// 已取消的批量不应推进多少点：轨迹点数必须远小于提交量。
	got, _ := eng.Track("truck-1")
	if len(got) > len(pts)/100 {
		t.Errorf("取消后仍摄入了 %d 个点，期望几乎不推进", len(got))
	}
	// 且整体应在远小于逐点处理 10w 点的时间内返回。
	if d := time.Since(start); d > time.Second {
		t.Errorf("IngestBatch 取消耗时 %v，疑似仍在逐点循环", d)
	}
}
