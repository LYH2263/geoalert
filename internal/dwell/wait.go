package dwell

import (
	"context"
	"time"
)

// Wait 等待对象进入后驻留达到 threshold；尊重 ctx 取消。
func (t *Tracker) Wait(ctx context.Context, objectID, fenceID string, threshold time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	t.mu.Lock()
	clk := t.clock
	s, ok := t.sess[key{Object: objectID, Fence: fenceID}]
	t.mu.Unlock()
	if clk == nil {
		return ErrNoClock
	}
	if !ok {
		if err := waitUntil(ctx, func() bool {
			t.mu.Lock()
			_, ok := t.sess[key{Object: objectID, Fence: fenceID}]
			t.mu.Unlock()
			return ok
		}, 5*time.Millisecond); err != nil {
			return err
		}
		t.mu.Lock()
		s = t.sess[key{Object: objectID, Fence: fenceID}]
		t.mu.Unlock()
		if s == nil {
			return ctx.Err()
		}
	}
	deadline := s.entered.Add(threshold)
	now := clk.Now()
	remain := deadline.Sub(now)
	if remain <= 0 {
		return nil
	}
	return Sleep(ctx, remain)
}

// Sleep 可取消睡眠。
func Sleep(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return ctx.Err()
	}

	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

func waitUntil(ctx context.Context, ready func() bool, tick time.Duration) error {
	tm := time.NewTicker(tick)
	defer tm.Stop()
	for {
		if ready() {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-tm.C:
		}
	}
}
