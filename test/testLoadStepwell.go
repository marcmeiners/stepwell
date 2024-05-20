package test

import (
	"fmt"
	"stepwell/extensions"
	"stepwell/stepwell"
	"sync"
	"time"
)

// handleCoreRequests processes requests for a given core, using side channels to stop the routines
func handleCoreRequests(stepwell *stepwell.StepWell, coreID uint64, stopChan <-chan struct{}, testRunning *bool, sumIsAllowed *int64, lock *sync.Mutex) {
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
			allowed := stepwell.IsAllowed(coreID, 1, requestTime)
			if allowed && *testRunning {
				num_allowed++
			}
		}
	}
}

func TestStepWellLoad(numCores uint64, bucketType int, duration int, refillRateInt int, capacityInt int) {
	capacity := int64(capacityInt)
	refillRate := float64(refillRateInt)
	numSeconds := time.Duration(duration) * time.Second

	testRunning := false
	var lock sync.Mutex
	totalAllowed := int64(0)

	stopChans := make([]chan struct{}, numCores)

	stepwell := stepwell.NewStepwell(numCores, time.Now(), bucketType, capacity, refillRate)

	for i := 0; i < int(numCores); i++ {
		stopChans[i] = make(chan struct{})
		go handleCoreRequests(stepwell, uint64(i), stopChans[i], &testRunning, &totalAllowed, &lock)
	}

	time.Sleep(1 * time.Second)
	testRunning = true
	time.Sleep(numSeconds)
	testRunning = false

	for _, stopChan := range stopChans {
		close(stopChan)
	}

	time.Sleep(1 * time.Second)

	expected_tokens := float64(numSeconds.Seconds()) * refillRate

	fmt.Println("Test completed.")
	fmt.Printf("Expected: %.2f Actual: %d", expected_tokens, totalAllowed)
}
