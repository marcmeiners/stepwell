// Test behavior of Stepwell / token buckets under certain amount of requests in certain time intervals

package test

import (
	"fmt"
	"stepwell/tokenbucket"
	"time"
)

func handleRequests(tokenbucket tokenbucket.TokenBucketInterface, coreID uint64, stopChan <-chan struct{}) {
	num_allowed := 0
	for {
		select {
		case <-stopChan: // Stop signal received
			fmt.Printf("Stopping requests for core %d\n", coreID)
			return
		default:
			requestTime := time.Now()
			allowed := tokenbucket.IsAllowed(1, requestTime)
			if allowed {
				num_allowed++
			}
		}
	}
}

func measureTokenAmount(tokenbucket tokenbucket.TokenBucketInterface, stopChan <-chan struct{}) {
	sleepDuration := time.Duration(1) * time.Nanosecond
	for {
		select {
		case <-stopChan: // Stop signal received
			fmt.Println("Stopping token measurement.")
			return
		default:
			tokens := tokenbucket.GetTokens()
			if tokens < 0 {
				fmt.Printf("Number of tokens left in the Token Bucket: %d\n", tokens)
			}
			time.Sleep(sleepDuration)
		}
	}
}

func TestTokenBucketLoad() {
	numCores := uint64(64)
	capacity := int64(10)
	refillRate := float64(1)
	bucketType := 1
	duration := 60 * time.Second

	tokenbucket := tokenbucket.NewTokenBucketByType(bucketType, capacity, refillRate, time.Now())

	stopChans := make([]chan struct{}, numCores)
	stopChanMeasurement := make(chan struct{})

	for i := 0; i < int(numCores); i++ {
		stopChans[i] = make(chan struct{})
		go handleRequests(tokenbucket, uint64(i), stopChans[i])
	}

	go measureTokenAmount(tokenbucket, stopChanMeasurement)

	time.Sleep(duration)

	for _, stopChan := range stopChans {
		close(stopChan)
	}
	close(stopChanMeasurement)

	// Wait a bit for goroutines to clean up before ending the test
	time.Sleep(1 * time.Second)
	fmt.Println("Test completed.")
}
