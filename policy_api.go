package geoalert

import "github.com/LYH2263/go-geoalert/internal/policy"

// FilterFencesByTags 按标签过滤围栏视图。
func (e *Engine) FilterFencesByTags(require, exclude []string) []FenceView {
	f := policy.TagFilter{Require: require, Exclude: exclude}
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := make([]FenceView, 0)
	for _, fence := range e.fences {
		if f.Match(fence.Tags) {
			out = append(out, fenceView(fence, e.active[fence.ID]))
		}
	}
	return out
}

// MergeFenceTags 合并标签到已有围栏。
func (e *Engine) MergeFenceTags(id string, tags []string) error {
	if e.closed.Load() {
		return ErrClosed
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	f, ok := e.fences[id]
	if !ok {
		return ErrNotFound
	}
	f.Tags = policy.MergeTags(f.Tags, tags)
	if err := e.store.SaveFence(toPersistFence(f)); err != nil {
		return wrapPersist(err)
	}
	e.fences[id] = f
	return nil
}
