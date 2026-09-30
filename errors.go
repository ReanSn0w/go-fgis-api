package fgis

import (
	"errors"
	"fmt"
)

var ErrClosed = errors.New("FGIS client is closed")

// HTTPError omits upstream response bodies because they may contain personal
// data or authentication material.
type HTTPError struct {
	Method     string
	Path       string
	StatusCode int
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("FGIS %s %s returned HTTP %d", e.Method, e.Path, e.StatusCode)
}
