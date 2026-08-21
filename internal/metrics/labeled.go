package metrics

import "sync"

// Labeled 带标签的子计数（围栏维度）。
type Labeled struct {
	mu   sync.Mutex
	byID map[string]*Registry
}

func NewLabeled() *Labeled {
	return &Labeled{byID: make(map[string]*Registry)}
}

func (l *Labeled) For(id string) *Registry {
	l.mu.Lock()
	defer l.mu.Unlock()
	r, ok := l.byID[id]
	if !ok {
		r = NewRegistry()
		l.byID[id] = r
	}
	return r
}

func (l *Labeled) All() map[string]map[string]int64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make(map[string]map[string]int64, len(l.byID))
	for id, r := range l.byID {
		out[id] = r.Snapshot()
	}
	return out
}
