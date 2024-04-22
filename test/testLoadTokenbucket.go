package test

import (
	"fmt"
	"stepwell/tokenbucket"
	"sync"
	"time"
)

func handleRequests(tokenbucket tokenbucket.TokenBucketInterface, coreID uint64, stopChan <-chan struct{}, testRunning *bool, sumIsAllowed *int64, lock *sync.Mutex) {
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

func TestTokenBucketLoad() {
	numCores := uint64(128)
	capacity := int64(10)
	refillRate := float64(1)
	bucketType := 1
	duration := 60 * time.Second

	testRunning := false
	var lock sync.Mutex
	totalAllowed := int64(0)

	stopChans := make([]chan struct{}, numCores)

	tokenbucket := tokenbucket.NewTokenBucketByType(bucketType, capacity, refillRate, time.Now())

	for i := 0; i < int(numCores); i++ {
		stopChans[i] = make(chan struct{})
		go handleRequests(tokenbucket, uint64(i), stopChans[i], &testRunning, &totalAllowed, &lock)
	}

	time.Sleep(1 * time.Second)
	testRunning = true
	time.Sleep(duration)
	testRunning = false

	for _, stopChan := range stopChans {
		close(stopChan)
	}

	time.Sleep(1 * time.Second)

	fmt.Println("Test completed.")
	fmt.Printf("Test Time: %.2f, Refill Rate: %.2f, Number of Tokens issued overall:%d", duration.Seconds(), refillRate, totalAllowed)
}
