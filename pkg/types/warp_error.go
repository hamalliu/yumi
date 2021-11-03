package types

import "fmt"

// WarpError ...
type WarpError string

// Warp ...
func (we WarpError) Warp(err error) error {
	return fmt.Errorf("%s: %w", we, err)
}
