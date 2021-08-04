package types

import "fmt"

type WarpError string

func (we WarpError) Warp(err error) error {
	return fmt.Errorf("%s: %w", we, err)
}
