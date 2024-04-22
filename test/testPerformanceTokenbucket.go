package test

import (
	"fmt"
	"stepwell/tokenbucket"
	"sync"
	"time"
)

func handleRequestsPerformance(tokenbucket tokenbucket.TokenBucketInterface, coreID uint64, stopChan <-chan struct{}, numIters int64, testRunning *bool, sumFinished *int64, lock *sync.Mutex) {
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

func TestTokenBucketPerformance() {
	numCores := uint64(8)
	capacity := int64(10)
	refillRate := float64(1)
	bucketType := 4
	numIters := int64(1000)

	testRunning := false
	var lock sync.Mutex
	sumFinished := int64(0)

	stopChans := make([]chan struct{}, numCores)

	tokenbucket := tokenbucket.NewTokenBucketByType(bucketType, capacity, refillRate, time.Now())

	for i := 0; i < int(numCores); i++ {
		stopChans[i] = make(chan struct{})
		go handleRequestsPerformance(tokenbucket, uint64(i), stopChans[i], numIters, &testRunning, &sumFinished, &lock)
	}

	time.Sleep(1 * time.Second)
	testRunning = true
	start := time.Now()

	var duration time.Duration

	//start time measurement
	for {
		if sumFinished == int64(numCores) {
			duration = time.Since(start)
			break
		}
	}

	for _, stopChan := range stopChans {
		close(stopChan)
	}

	time.Sleep(1 * time.Second)

	fmt.Println("Test completed. Execution time in nanoseconds: ")
	fmt.Println(duration.Nanoseconds())
}
