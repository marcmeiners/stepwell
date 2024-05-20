//go:build !linux
// +build !linux

package extensions

import "errors"

func PinToCore(coreID int) error {
	return errors.New("PinToCore function is not supported on this platform")
}
