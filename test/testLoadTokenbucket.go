package test

import (
	"fmt"
	"stepwell/extensions"
	"stepwell/tokenbucket"
	"sync"
	"time"
)

func handleRequests(tokenbucket tokenbucket.TokenBucketInterface, coreID uint64, stopChan <-chan struct{}, testRunning *bool, sumIsAllowed *int64, lock *sync.Mutex) {
	err := extensions.PinToCore(int(coreID))
	if err != nil {
		fmt.Printf("Failed to pin goroutine to core %d: %v\n", coreID, err)
	}
	num_allowed := int64(0)
	for {
		select {
		case <-stopChan: // Stop signal received
			fmt.Printf("Stopping requests for core %d\n", coreID)
			lock.Lock() // Acquire the lock before modifying the shared resource
			*sumIsAllowed += num_allowed
			lock.Unlock()
			return
		default:
			requestTime := time.Now()
			allowed := tokenbucket.IsAllowed(1, requestTime)
			if allowed && *testRunning {
				num_allowed++
			}
		}
	}
}

func TestTokenBucketLoad(numCores uint64, bucketType int, duration int, refillRateInt int, capacityInt int) {
	capacity := int64(capacityInt)
	refillRate := float64(refillRateInt)
	numSeconds := time.Duration(duration) * time.Second

	testRunning := false
	var lock sync.Mutex
	totalAllowed := int64(0)

	stopChans := make([]chan struct{}, numCores)

	tokenbucket := tokenbucket.NewTokenBucketByType(bucketType, capacity, refillRate, time.Now())

	for i := 0; i < int(numCores); i++ {
		stopChans[i] = make(chan struct{})
		go handleRequests(tokenbucket, uint64(i), stopChans[i], &testRunning, &totalAllowed, &lock)
	}

	time.Sleep(500 * time.Millisecond)
	testRunning = true
	time.Sleep(numSeconds)
	testRunning = false

	for _, stopChan := range stopChans {
		close(stopChan)
	}

	time.Sleep(500 * time.Millisecond)

	expected_tokens := float64(numSeconds.Seconds()) * refillRate

	fmt.Println("Test completed.")
	fmt.Printf("Expected: %.2f Actual: %d", expected_tokens, totalAllowed)
}
