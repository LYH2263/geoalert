package edge

import "time"

// History 记录最近若干次边沿，便于回放诊断。
type History struct {
	max int
	buf []Transition
}

func NewHistory(max int) *History {
	if max <= 0 {
		max = 64
	}
	return &History{max: max, buf: make([]Transition, 0, max)}
}

func (h *History) Push(tr Transition) {
	if tr.Kind == None || tr.Kind == Stay {
		return
	}
	h.buf = append(h.buf, tr)
	if len(h.buf) > h.max {
		overflow := len(h.buf) - h.max
		copy(h.buf, h.buf[overflow:])
		h.buf = h.buf[:h.max]
	}
}

func (h *History) All() []Transition {
	out := make([]Transition, len(h.buf))
	copy(out, h.buf)
	return out
}

func (h *History) Since(t time.Time) []Transition {
	out := make([]Transition, 0, len(h.buf))
	for _, tr := range h.buf {
		if !tr.At.Before(t) {
			out = append(out, tr)
		}
	}
	return out
}

func (h *History) Clear() {
	h.buf = h.buf[:0]
}
