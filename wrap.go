package geoalert

import "fmt"

func wrapInvalidFence(err error) error {
	if err == nil {
		return ErrInvalidFence
	}

	return fmt.Errorf("%w: %v", ErrInvalidFence, err)
}

func wrapPersist(err error) error {
	if err == nil {
		return ErrPersist
	}
	return fmt.Errorf("%w: %v", ErrPersist, err)
}
