package persist

// LatLng 持久化坐标。
type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// FenceRecord 围栏持久化记录。
type FenceRecord struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Kind         int      `json:"kind"`
	Vertices     []LatLng `json:"vertices,omitempty"`
	CenterLat    float64  `json:"center_lat,omitempty"`
	CenterLng    float64  `json:"center_lng,omitempty"`
	RadiusM      float64  `json:"radius_m,omitempty"`
	DwellAfterMs int64    `json:"dwell_after_ms,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	AlertOnEnter bool     `json:"alert_on_enter"`
	AlertOnExit  bool     `json:"alert_on_exit"`
	AlertOnDwell bool     `json:"alert_on_dwell"`
}

// Store 围栏持久化接口。
type Store interface {
	SaveFence(FenceRecord) error
	DeleteFence(id string) error
	ListFences() ([]FenceRecord, error)
	Flush() error
	Close() error
}
