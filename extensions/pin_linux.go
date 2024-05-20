//go:build linux
// +build linux

package extensions

import (
	"runtime"

	"golang.org/x/sys/unix"
)

func PinToCore(coreID int) error {
	runtime.LockOSThread() // Lock the current goroutine to its current OS thread

	var cpuset unix.CPUSet
	cpuset.Zero()
	cpuset.Set(coreID)

	pid := unix.Gettid() // Get the thread ID of the calling thread
	return unix.SchedSetaffinity(pid, &cpuset)
}
