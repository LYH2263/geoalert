package persist

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// Memory 内存 Store。
type Memory struct {
	mu     sync.Mutex
	fences map[string]FenceRecord
	fail   bool
}

func NewMemory() *Memory {
	return &Memory{fences: make(map[string]FenceRecord)}
}

// SetFail 测试用：后续 Save 失败。
func (m *Memory) SetFail(v bool) { m.fail = v }

func (m *Memory) SaveFence(r FenceRecord) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.fail {
		return errPersist
	}
	cp := r
	cp.Vertices = cloneVerts(r.Vertices)
	cp.Tags = cloneStr(r.Tags)
	m.fences[r.ID] = cp
	return nil
}

func (m *Memory) DeleteFence(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.fences, id)
	return nil
}

func (m *Memory) ListFences() ([]FenceRecord, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]FenceRecord, 0, len(m.fences))
	for _, r := range m.fences {
		cp := r
		cp.Vertices = cloneVerts(r.Vertices)
		cp.Tags = cloneStr(r.Tags)
		out = append(out, cp)
	}
	return out, nil
}

func (m *Memory) Flush() error { return nil }
func (m *Memory) Close() error { return nil }

var errPersist = &persistError{"persist failed"}

type persistError struct{ s string }

func (e *persistError) Error() string { return e.s }

func cloneVerts(src []LatLng) []LatLng {
	if src == nil {
		return nil
	}
	dst := make([]LatLng, len(src))
	copy(dst, src)
	return dst
}

func cloneStr(src []string) []string {
	if src == nil {
		return nil
	}
	dst := make([]string, len(src))
	copy(dst, src)
	return dst
}

// FileStore JSON 文件 Store。
type FileStore struct {
	mu   sync.Mutex
	path string
	mem  *Memory
}

func OpenFile(path string) (*FileStore, error) {
	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	fs := &FileStore{path: path, mem: NewMemory()}
	b, err := os.ReadFile(path)
	if err == nil && len(b) > 0 {
		var list []FenceRecord
		if err := json.Unmarshal(b, &list); err != nil {
			return nil, err
		}
		for _, r := range list {
			_ = fs.mem.SaveFence(r)
		}
	} else if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return fs, nil
}

func (f *FileStore) SaveFence(r FenceRecord) error {
	if err := f.mem.SaveFence(r); err != nil {
		return err
	}
	return f.Flush()
}

func (f *FileStore) DeleteFence(id string) error {
	if err := f.mem.DeleteFence(id); err != nil {
		return err
	}
	return f.Flush()
}

func (f *FileStore) ListFences() ([]FenceRecord, error) {
	return f.mem.ListFences()
}

func (f *FileStore) Flush() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	list, err := f.mem.ListFences()
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	tmp := f.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, f.path)
}

func (f *FileStore) Close() error { return f.Flush() }
