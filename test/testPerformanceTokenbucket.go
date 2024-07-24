package test

import (
	"fmt"
	"stepwell/extensions"
	"stepwell/tokenbucket"
	"sync"
	"time"
)

func handleRequestsPerformance(tokenbucket tokenbucket.TokenBucketInterface, coreID uint64, stopChan <-chan struct{}, numIters int64, testRunning *bool, sumFinished *int64, lock *sync.Mutex) {
	err := extensions.PinToCore(int(coreID))
	if err != nil {
		fmt.Printf("Failed to pin goroutine to core %d: %v\n", coreID, err)
	}
	num_executed := int64(0)
	for {
		select {
		case <-stopChan: // Stop signal received
			fmt.Printf("Stopping requests for core %d\n", coreID)
			return
		default:
			requestTime := time.Now()
			tokenbucket.IsAllowed(1, requestTime)
			if *testRunning {
				num_executed++
			}
			if num_executed == numIters {
				lock.Lock()
				*sumFinished++
				lock.Unlock()
				return
			}
		}
	}
}

func TestTokenBucketPerformance(numCores uint64, bucketType int, duration int, refillRateInt int, capacityInt int) {
	capacity := int64(capacityInt)
	refillRate := float64(refillRateInt)
	numIters := int64(duration)

	testRunning := false
	var lock sync.Mutex
	sumFinished := int64(0)

	stopChans := make([]chan struct{}, numCores)

	tokenbucket := tokenbucket.NewTokenBucketByType(bucketType, capacity, refillRate, time.Now())

	for i := 0; i < int(numCores); i++ {
		stopChans[i] = make(chan struct{})
		go handleRequestsPerformance(tokenbucket, uint64(i), stopChans[i], numIters, &testRunning, &sumFinished, &lock)
	}

	time.Sleep(500 * time.Millisecond)
	testRunning = true
	start := time.Now()

	var measuredDuration time.Duration

	//start time measurement
	for {
		if sumFinished == int64(numCores) {
			measuredDuration = time.Since(start)
			break
		}
	}

	for _, stopChan := range stopChans {
		close(stopChan)
	}

	time.Sleep(500 * time.Millisecond)

	fmt.Printf("Time: %d", measuredDuration.Nanoseconds())
}
