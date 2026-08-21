package edge_test

import (
	"testing"
	"time"

	"github.com/LYH2263/go-geoalert/internal/edge"
)

func TestDetectorEnterExit(t *testing.T) {
	d := edge.New()
	now := time.Now()
	tr := d.Observe("o", "f", edge.Outside, now)
	if tr.Kind != edge.None {
		t.Fatalf("first outside: %v", tr.Kind)
	}
	tr = d.Observe("o", "f", edge.Inside, now.Add(time.Second))
	if tr.Kind != edge.Enter {
		t.Fatalf("enter: %v", tr.Kind)
	}
	tr = d.Observe("o", "f", edge.Inside, now.Add(2*time.Second))
	if tr.Kind != edge.Stay {
		t.Fatalf("stay: %v", tr.Kind)
	}
	tr = d.Observe("o", "f", edge.Outside, now.Add(3*time.Second))
	if tr.Kind != edge.Exit {
		t.Fatalf("exit: %v", tr.Kind)
	}
}
