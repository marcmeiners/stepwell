package main

import (
	"fmt"
	"stepwell/core"
	"stepwell/tokenbucket"
	"time"
)

func main() {
	testTokenBucketHelia()
}

func testTokenBucketHelia() {
	now := time.Now()
	tb := tokenbucket.NewTokenBucketHelia(1, 0.5, now)

	requestTimes := []time.Duration{0, time.Second, 2 * time.Second, 3 * time.Second, 5 * time.Second}

	for _, duration := range requestTimes {
		requestTime := now.Add(duration)
		allowed := tb.IsAllowed(1, requestTime)
		fmt.Printf("Request at %s allowed: %v\n", requestTime.Format(time.RFC3339), allowed)
	}
}

func testStepWell() {
	stepwell := core.NewStepwell(8, time.Now(), 1, 10, 5)
	now := time.Now()
	requestTimes := []time.Duration{0, time.Second, 2 * time.Second, 3 * time.Second, 5 * time.Second}

	for _, duration := range requestTimes {
		requestTime := now.Add(duration)
		allowed := stepwell.IsAllowed(1, 1, requestTime)
		fmt.Printf("Request at %s allowed: %v\n", requestTime.Format(time.RFC3339), allowed)
	}
}

// handleCoreRequests processes requests for a given core
func handleCoreRequests(stepwell *core.StepWell, coreID uint64) {
	now := time.Now()
	requestTimes := []time.Duration{0, 100 * time.Millisecond, 200 * time.Millisecond, 300 * time.Millisecond, 400 * time.Millisecond}

	for _, duration := range requestTimes {
		requestTime := now.Add(duration)
		allowed := stepwell.IsAllowed(coreID, 1, requestTime)
		fmt.Printf("Request at %s for core %d allowed: %v\n", requestTime.Format(time.RFC3339), coreID, allowed)
		// Simulate processing time or delay between requests for this core
		time.Sleep(time.Millisecond * 100)
	}
}

func testStepWellOverflow() {
	numCores := uint64(8)
	capacity := uint64(3)
	refillRate := float64(2)
	bucketType := 1

	stepwell := core.NewStepwell(numCores, time.Now(), bucketType, capacity, refillRate)

	// Launch separate goroutines for core 0 and core 7 using the defined functions
	go handleCoreRequests(stepwell, 0) // Core 0
	go handleCoreRequests(stepwell, 2) // Core 2
	go handleCoreRequests(stepwell, 4) // Core 4
	go handleCoreRequests(stepwell, 7) // Core 7

	// Wait enough time for both goroutines to complete their work
	time.Sleep(2 * time.Second)
}
