package dwell

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/LYH2263/go-geoalert/internal/clock"
)

// ErrNoClock 未配置时钟。
var ErrNoClock = errors.New("dwell: no clock")

type key struct {
	Object string
	Fence  string
}

type session struct {
	entered   time.Time
	threshold time.Duration
	fired     bool
}

// Event 驻留达成事件。
type Event struct {
	ObjectID string
	FenceID  string
	Duration time.Duration
	At       time.Time
}

// Tracker 驻留计时。
type Tracker struct {
	mu    sync.Mutex
	clock clock.Clock
	sess  map[key]*session
}

func New() *Tracker {
	return &Tracker{sess: make(map[key]*session)}
}

func (t *Tracker) SetClock(c clock.Clock) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.clock = c
}

func (t *Tracker) Enter(objectID, fenceID string, at time.Time, threshold time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.sess[key{Object: objectID, Fence: fenceID}] = &session{entered: at, threshold: threshold}
}

func (t *Tracker) Exit(objectID, fenceID string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.sess, key{Object: objectID, Fence: fenceID})
}

// Check 检查是否达到驻留；需要 clock。
func (t *Tracker) Check(ctx context.Context, objectID, fenceID string, at time.Time) (Event, bool, error) {

	_ = ctx
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.clock == nil {
		return Event{}, false, ErrNoClock
	}
	s, ok := t.sess[key{Object: objectID, Fence: fenceID}]
	if !ok || s.fired {
		return Event{}, false, nil
	}
	dur := at.Sub(s.entered)
	if dur < s.threshold {
		return Event{}, false, nil
	}
	s.fired = true
	return Event{ObjectID: objectID, FenceID: fenceID, Duration: dur, At: at}, true, nil
}

func (t *Tracker) ClearFence(fenceID string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for k := range t.sess {
		if k.Fence == fenceID {
			delete(t.sess, k)
		}
	}
}

func (t *Tracker) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.sess = make(map[key]*session)
}

// Active 是否正在计时。
func (t *Tracker) Active(objectID, fenceID string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	_, ok := t.sess[key{Object: objectID, Fence: fenceID}]
	return ok
}
