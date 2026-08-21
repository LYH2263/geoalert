package edge

import (
	"testing"
	"time"
)

// Clear 后 Observe 不应 panic：state 必须保持可写映射而非 nil。
// 复现运维脚本 Close 后误点 Ingest 触发的 nil-map 写入崩溃。
func TestDetectorClearKeepsStateWritable(t *testing.T) {
	d := New()
	d.Observe("o1", "f1", Inside, time.Now())

	d.Clear()

	// 旧实现这里 d.state 为 nil，写操作会 panic。
	tr := d.Observe("o2", "f1", Inside, time.Now())
	if tr.Kind != Enter {
		t.Fatalf("after Clear, fresh inside observation = %s, want enter", tr.Kind)
	}
	if d.state == nil {
		t.Fatal("state is nil after Clear; expected non-nil writable map")
	}
	// Inside 查询在清空后应返回未见过的状态。
	if _, seen := d.Inside("o2", "f1"); !seen {
		t.Fatal("Inside should report seen after a fresh Observe")
	}
}
