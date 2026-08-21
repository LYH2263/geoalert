package audit

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBug09_AuditRotateClosesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.jsonl")
	l, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	if err := l.Write("register", "fence-a", "circle"); err != nil {
		t.Fatal(err)
	}
	if err := l.Rotate(); err != nil {
		t.Fatalf("rotate failed (handle still held?): %v", err)
	}
	list, err := ListRotated(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) == 0 {
		t.Fatal("expected rotated audit file")
	}
	if err := l.Write("ingest", "obj-1", "enter"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("new audit file missing: %v", err)
	}
}
