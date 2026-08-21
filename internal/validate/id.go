package validate

import (
	"errors"
	"strings"
	"unicode"
)

var (
	ErrEmptyID         = errors.New("empty id")
	ErrBadID           = errors.New("invalid id chars")
	ErrTooFewVertices  = errors.New("polygon needs >= 3 vertices")
	ErrSelfIntersect   = errors.New("polygon self-intersects")
	ErrBadRadius       = errors.New("radius must be > 0")
	ErrBadKind         = errors.New("unknown fence kind")
	ErrBadLat          = errors.New("latitude out of range")
	ErrBadLng          = errors.New("longitude out of range")
	ErrDupVertex       = errors.New("duplicate consecutive vertices")
	ErrEmptyObject     = errors.New("empty object id")
)

// FenceID 校验围栏 ID。
func FenceID(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return ErrEmptyID
	}
	for _, r := range id {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' || r == '.' {
			continue
		}
		return ErrBadID
	}
	return nil
}

// ObjectID 校验对象 ID。
func ObjectID(id string) error {
	if strings.TrimSpace(id) == "" {
		return ErrEmptyObject
	}
	return FenceID(id)
}

// LatLng 校验坐标。
func LatLng(lat, lng float64) error {
	if lat < -90 || lat > 90 {
		return ErrBadLat
	}
	if lng < -180 || lng > 180 {
		return ErrBadLng
	}
	return nil
}
