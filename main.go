package main

import (
	"fmt"
	"stepwell/core"
	"stepwell/tokenbucket"
	"time"
)

func main() {
	testStepWellOverflow()
}

func testTokenBucket() {
	tb := tokenbucket.NewTokenBucketLock(10, 1, time.Now())

	now := time.Now()
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

func testStepWellOverflow() {
	// Initialize StepWell with a small capacity and refill rate
	// to ensure it will overflow with a burst of requests.
	numCores := uint64(8)
	capacity := uint64(5)
	refillRate := float64(1)
	bucketType := 1
	stepwell := core.NewStepwell(numCores, time.Now(), bucketType, capacity, refillRate)

	now := time.Now()
	// Simulate a burst of requests in quick succession
	requestTimes := []time.Duration{0, 100 * time.Millisecond, 200 * time.Millisecond, 300 * time.Millisecond, 400 * time.Millisecond}

	for _, duration := range requestTimes {
		requestTime := now.Add(duration)
		// Simulate multiple requests at each time
		allowed := stepwell.IsAllowed(0, 1, requestTime)
		fmt.Printf("Request at %s for core %d allowed: %v\n", requestTime.Format(time.RFC3339), 0, allowed)
		allowed = stepwell.IsAllowed(1, 1, requestTime)
		fmt.Printf("Request at %s for core %d allowed: %v\n", requestTime.Format(time.RFC3339), 1, allowed)
		allowed = stepwell.IsAllowed(7, 1, requestTime)
		fmt.Printf("Request at %s for core %d allowed: %v\n", requestTime.Format(time.RFC3339), 7, allowed)
	}
}
