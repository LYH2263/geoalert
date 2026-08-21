package edge

import "time"

// Presence 当前是否在围栏内。
type Presence bool

const (
	Outside Presence = false
	Inside  Presence = true
)

// Kind 边沿类型。
type Kind int

const (
	None Kind = iota
	Enter
	Exit
	Stay
)

func (k Kind) String() string {
	switch k {
	case Enter:
		return "enter"
	case Exit:
		return "exit"
	case Stay:
		return "stay"
	default:
		return "none"
	}
}

// Transition 一次观察结果。
type Transition struct {
	ObjectID string
	FenceID  string
	Kind     Kind
	At       time.Time
	Inside   bool
}

type key struct {
	Object string
	Fence  string
}

type cell struct {
	inside bool
	seen   bool
	at     time.Time
}
