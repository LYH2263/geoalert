package geoalert

import (
	"testing"

	"github.com/LYH2263/go-geoalert/internal/persist"
)

// newTestEngine 构造一个最小可用引擎（内存存储，无外部 IO）。
func newTestEngine(t *testing.T) *Engine {
	t.Helper()
	e, err := New(Options{Store: persist.NewMemory()})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return e
}

// emitRawLocked 绕过公开 API 注入一条告警，便于在受控状态测试队列拷贝语义。
func (e *Engine) emitRawLocked(a Alert) {
	e.mu.Lock()
	defer e.mu.Unlock()
	a.ID = e.nextAlertID()
	a.Seq = e.seq
	e.alerts = append(e.alerts, a)
}

// TestAlerts_ReturnsCopy 断言 Alerts() 返回独立拷贝：调用方改返回值不污染内部队列。
func TestAlerts_ReturnsCopy(t *testing.T) {
	e := newTestEngine(t)
	e.emitRawLocked(Alert{FenceID: "f1", ObjectID: "obj-A", Kind: AlertEnter})

	got := e.Alerts()
	if len(got) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(got))
	}
	original := got[0].ObjectID

	// 调用方在返回切片上做标注修改。
	got[0].ObjectID = "ANNOTATED"

	// 再拉一次，内部必须仍是原值。
	again := e.Alerts()
	if again[0].ObjectID != original {
		t.Fatalf("internal alert mutated: want %q, got %q", original, again[0].ObjectID)
	}
	if got[0].ObjectID != "ANNOTATED" {
		t.Fatalf("caller copy should keep annotation, got %q", got[0].ObjectID)
	}

	// 返回切片底层数组也不应与内部共享（cap 之外的写也不应串味）。
	if cap(got) > 0 && len(e.alerts) > 0 {
		_ = append(got[:0], Alert{ObjectID: "X"})
		if e.alerts[0].ObjectID == "X" {
			t.Fatalf("Alerts() returned slice shares backing array with internal queue")
		}
	}
}

// TestSnapshot_AlertsAreCopy 断言 Snapshot.Alerts 与内部队列解耦。
func TestSnapshot_AlertsAreCopy(t *testing.T) {
	e := newTestEngine(t)
	e.emitRawLocked(Alert{FenceID: "f1", ObjectID: "obj-B", Kind: AlertExit})

	snap := e.Snapshot()
	if len(snap.Alerts) != 1 {
		t.Fatalf("expected 1 alert in snapshot, got %d", len(snap.Alerts))
	}
	original := snap.Alerts[0].ObjectID

	snap.Alerts[0].ObjectID = "MARKED"

	again := e.Snapshot()
	if again.Alerts[0].ObjectID != original {
		t.Fatalf("internal alert mutated via snapshot: want %q, got %q", original, again.Alerts[0].ObjectID)
	}
}

// TestDrainAlerts_IndependentFromLaterQueue 断言 DrainAlerts 返回的切片
// 不与后续队列追加共享底层数组（旧 bug：cloneAlerts 未拷贝时复用槽位会串味）。
func TestDrainAlerts_IndependentFromLaterQueue(t *testing.T) {
	e := newTestEngine(t)
	e.emitRawLocked(Alert{FenceID: "f1", ObjectID: "obj-C", Kind: AlertEnter})

	drained := e.DrainAlerts()
	if len(drained) != 1 || drained[0].ObjectID != "obj-C" {
		t.Fatalf("drained = %+v, want obj-C", drained)
	}

	// 排空后再追加新告警；若共享底层数组，drained[0] 会被新值覆盖。
	e.emitRawLocked(Alert{FenceID: "f1", ObjectID: "obj-D", Kind: AlertExit})

	if drained[0].ObjectID != "obj-C" {
		t.Fatalf("drained slice aliasing internal queue: want obj-C, got %q", drained[0].ObjectID)
	}

	cur := e.Alerts()
	if len(cur) != 1 || cur[0].ObjectID != "obj-D" {
		t.Fatalf("queue after drain+append = %+v, want obj-D", cur)
	}
}
