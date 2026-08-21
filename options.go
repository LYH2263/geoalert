package geoalert

import (
	"time"

	"github.com/LYH2263/go-geoalert/internal/audit"
	"github.com/LYH2263/go-geoalert/internal/clock"
	"github.com/LYH2263/go-geoalert/internal/persist"
)

// Options 创建 Engine。
type Options struct {
	Clock               clock.Clock
	Store               persist.Store
	Auditor             *audit.Logger
	PersistPath         string
	AuditPath           string
	AlertLogPath        string
	MaxTrackPoints      int
	MaxAlerts           int
	DefaultDwell        time.Duration
	StrictSelfIntersect bool
}

func (o *Options) withDefaults() Options {
	out := *o
	if out.MaxTrackPoints <= 0 {
		out.MaxTrackPoints = 256
	}
	if out.MaxAlerts <= 0 {
		out.MaxAlerts = 1024
	}
	if out.DefaultDwell < 0 {
		out.DefaultDwell = 0
	}
	return out
}
