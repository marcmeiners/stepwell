package test

import (
	"fmt"
	"stepwell/stepwell"
	"sync"
	"time"
)

// handleCoreRequests processes requests for a given core, using side channels to stop the routines
func handleCoreRequestsPerformance(stepwell *stepwell.StepWell, coreID uint64, stopChan <-chan struct{}, numIters int64, testRunning *bool, sumFinished *int64, lock *sync.Mutex) {
	num_executed := int64(0)
	for {
		select {
		case <-stopChan: // Stop signal received
			fmt.Printf("Stopping requests for core %d\n", coreID)
			return
		default:
			requestTime := time.Now()
			stepwell.IsAllowed(coreID, 1, requestTime)
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

func TestStepWellPerformance() {
	numCores := uint64(8)
	capacity := int64(10)
	refillRate := float64(1)
	bucketType := 1
	numIters := int64(1000)

	testRunning := false
	var lock sync.Mutex
	sumFinished := int64(0)

	stopChans := make([]chan struct{}, numCores)

	stepwell := stepwell.NewStepwell(numCores, time.Now(), bucketType, capacity, refillRate)

	for i := 0; i < int(numCores); i++ {
		stopChans[i] = make(chan struct{})
		go handleCoreRequestsPerformance(stepwell, uint64(i), stopChans[i], numIters, &testRunning, &sumFinished, &lock)
	}

	time.Sleep(1 * time.Second)

	testRunning = true
	start := time.Now()
	var duration time.Duration

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