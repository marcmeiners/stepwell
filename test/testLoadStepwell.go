// Test behavior of Stepwell / token buckets under certain amount of requests in certain time intervals

package test

import (
	"fmt"
	"stepwell/stepwell"
	"time"
)

// handleCoreRequests processes requests for a given core, using side channels to stop the routines
func handleCoreRequests(stepwell *stepwell.StepWell, coreID uint64, stopChan <-chan struct{}) {
	num_allowed := 0
	for {
		select {
		case <-stopChan: // Stop signal received
			fmt.Printf("Stopping requests for core %d\n", coreID)
			return
		default:
			requestTime := time.Now()
			allowed := stepwell.IsAllowed(coreID, 1, requestTime)
			if allowed {
				num_allowed++
			}
		}
	}
}

func measureTokenAmountStepwell(stepwell *stepwell.StepWell, stopChan <-chan struct{}) {
	sleepDuration := time.Duration(1) * time.Nanosecond
	for {
		select {
		case <-stopChan: // Stop signal received
			fmt.Println("Stopping token measurement.")
			return
		default:
			//Go all the way up to the root token bucket of the stepwell structure
			tokens := stepwell.Cores[0].Parent.Parent.Parent.Parent.TokenBucket.GetTokens()
			if tokens < 0 {
				fmt.Printf("Number of tokens left in the Token Bucket: %d\n", tokens)
			}

			time.Sleep(sleepDuration)
		}
	}
}

func TestStepWellLoad() {
	numCores := uint64(64)
	capacity := int64(10)
	refillRate := float64(1)
	bucketType := 1
	duration := 60 * time.Second

	stopChans := make([]chan struct{}, numCores)
	stopChanMeasurement := make(chan struct{})

	stepwell := stepwell.NewStepwell(numCores, time.Now(), bucketType, capacity, refillRate)

	for i := 0; i < int(numCores); i++ {
		stopChans[i] = make(chan struct{})
		go handleCoreRequests(stepwell, uint64(i), stopChans[i])
	}

	go measureTokenAmountStepwell(stepwell, stopChanMeasurement)

	time.Sleep(duration)

	for _, stopChan := range stopChans {
		close(stopChan)
	}
	close(stopChanMeasurement)

	// Wait a bit for goroutines to clean up before ending the test
	time.Sleep(1 * time.Second)
	fmt.Println("Test completed.")
}
