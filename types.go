package geoalert

import "time"

// FenceKind 围栏几何类型。
type FenceKind int

const (
	FencePolygon FenceKind = iota
	FenceCircle
)

func (k FenceKind) String() string {
	switch k {
	case FencePolygon:
		return "polygon"
	case FenceCircle:
		return "circle"
	default:
		return "unknown"
	}
}

// LatLng 经纬度点（度）。
type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// Fence 地理围栏定义。Vertices 仅多边形有效；Circle 使用 Center+RadiusM。
type Fence struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Kind         FenceKind     `json:"kind"`
	Vertices     []LatLng     `json:"vertices,omitempty"`
	Center       LatLng       `json:"center,omitempty"`
	RadiusM      float64      `json:"radius_m,omitempty"`
	DwellAfter   time.Duration `json:"dwell_after,omitempty"`
	Tags         []string     `json:"tags,omitempty"`
	AlertOnEnter bool         `json:"alert_on_enter"`
	AlertOnExit  bool         `json:"alert_on_exit"`
	AlertOnDwell bool         `json:"alert_on_dwell"`
}

// TrackPoint 对象在某一时刻的位置。
type TrackPoint struct {
	ObjectID string            `json:"object_id"`
	At       time.Time        `json:"at"`
	Lat      float64          `json:"lat"`
	Lng      float64          `json:"lng"`
	SpeedMps float64          `json:"speed_mps,omitempty"`
	Heading  float64          `json:"heading,omitempty"`
	Meta     map[string]string `json:"meta,omitempty"`
}

// AlertKind 告警类型。
type AlertKind int

const (
	AlertEnter AlertKind = iota + 1
	AlertExit
	AlertDwell
)

func (k AlertKind) String() string {
	switch k {
	case AlertEnter:
		return "enter"
	case AlertExit:
		return "exit"
	case AlertDwell:
		return "dwell"
	default:
		return "unknown"
	}
}

// Alert 围栏事件。
type Alert struct {
	ID       string        `json:"id"`
	FenceID  string        `json:"fence_id"`
	ObjectID string        `json:"object_id"`
	Kind     AlertKind     `json:"kind"`
	At       time.Time     `json:"at"`
	Lat      float64       `json:"lat"`
	Lng      float64       `json:"lng"`
	DwellFor time.Duration `json:"dwell_for,omitempty"`
	Seq      uint64        `json:"seq"`
}

// ObjectState 对象相对某围栏的内外状态。
type ObjectState int

const (
	StateUnknown ObjectState = iota
	StateOutside
	StateInside
)

func (s ObjectState) String() string {
	switch s {
	case StateOutside:
		return "outside"
	case StateInside:
		return "inside"
	default:
		return "unknown"
	}
}

// FenceView 对外围栏快照（Vertices 已拷贝）。
type FenceView struct {
	ID         string        `json:"id"`
	Name       string        `json:"name"`
	Kind       FenceKind     `json:"kind"`
	Vertices   []LatLng      `json:"vertices,omitempty"`
	Center     LatLng        `json:"center,omitempty"`
	RadiusM    float64       `json:"radius_m,omitempty"`
	DwellAfter time.Duration  `json:"dwell_after,omitempty"`
	Tags       []string      `json:"tags,omitempty"`
	Active     bool          `json:"active"`
}

// TrackView 轨迹摘要。
type TrackView struct {
	ObjectID string     `json:"object_id"`
	Points   int        `json:"points"`
	Last     TrackPoint `json:"last"`
	Inside   []string   `json:"inside,omitempty"`
}

// EngineStats 运行计数。
type EngineStats struct {
	Fences       int   `json:"fences"`
	Objects      int   `json:"objects"`
	Ingests      int64 `json:"ingests"`
	Enters       int64 `json:"enters"`
	Exits        int64 `json:"exits"`
	Dwells       int64 `json:"dwells"`
	AlertsQueued int   `json:"alerts_queued"`
	Closed       bool  `json:"closed"`
}
