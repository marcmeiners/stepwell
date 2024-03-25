package main

import (
	"fmt"
	"stepwell/core"
	"stepwell/tokenbucket"
	"time"
)

func main() {
	testStepWell()
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
