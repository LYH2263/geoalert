package geoalert

import "errors"

var (
	ErrClosed         = errors.New("geoalert: closed")
	ErrExists         = errors.New("geoalert: fence already exists")
	ErrNotFound       = errors.New("geoalert: fence not found")
	ErrInvalidFence   = errors.New("geoalert: invalid fence")
	ErrInvalidPoint   = errors.New("geoalert: invalid lat/lng")
	ErrInvalidObject  = errors.New("geoalert: invalid object id")
	ErrNoClock        = errors.New("geoalert: no clock")
	ErrPersist        = errors.New("geoalert: persist failed")
	ErrCanceled       = errors.New("geoalert: canceled")
	ErrEmptyTrack     = errors.New("geoalert: empty track")
	ErrDwellThreshold = errors.New("geoalert: invalid dwell threshold")
)
