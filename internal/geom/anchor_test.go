package geom

import "testing"

// 回归：空环不得 panic，必须返回 ok=false。
// 修复前 RingAnchor 无条件 ring[0]，空切片直接越界 panic。
func TestRingAnchor_EmptyRingNoPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("RingAnchor panicked on empty ring: %v", r)
		}
	}()

	if _, ok := RingAnchor(nil); ok {
		t.Fatal("nil ring must return ok=false")
	}
	if _, ok := RingAnchor([]Point{}); ok {
		t.Fatal("empty ring must return ok=false")
	}

	// 非空环仍应返回首顶点 + ok=true。
	p := Point{Lat: 1, Lng: 2}
	got, ok := RingAnchor([]Point{p, {Lat: 3, Lng: 4}})
	if !ok {
		t.Fatal("non-empty ring must return ok=true")
	}
	if got != p {
		t.Fatalf("anchor = %+v, want %+v", got, p)
	}
}
