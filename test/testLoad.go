// Test behavior of Stepwell / token buckets under certain amount of requests in certain time intervals

package test

import (
	"fmt"
	"math/rand"
	"stepwell/stepwell"
	"stepwell/tokenbucket"
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
			rand.Seed(time.Now().UnixNano())
			randomDuration := time.Duration(rand.Intn(6)) * time.Nanosecond
			// Sleep for the random duration
			time.Sleep(randomDuration)
			requestTime := time.Now()
			allowed := stepwell.IsAllowed(coreID, 1, requestTime)
			if allowed {
				num_allowed++
				//fmt.Printf("Num Allowed for core %d: %d\n", coreID, num_allowed)
				//fmt.Printf("Request at %s for core %d allowed: %v\n", requestTime.Format(time.RFC3339), coreID, allowed)
			}
		}
	}
}

func measureTokenAmount(stepwell *stepwell.StepWell, stopChan <-chan struct{}) {
	sleepDuration := time.Duration(100) * time.Nanosecond
	for {
		select {
		case <-stopChan: // Stop signal received
			fmt.Println("Stopping token measurement.")
			return
		default:
			//Go all the way up to the root token bucket of the stepwell structure
			tokens := stepwell.Cores[0].Parent.Parent.Parent.Parent.TokenBucket.GetTokens()
			capacity := stepwell.Cores[0].Parent.Parent.Parent.Parent.TokenBucket.GetCapacity()
			if tokens > 0 {
				fmt.Printf("Capacity of root TokenBucket: %d, Actual Tokens: %d\n",
					capacity,
					tokens)
			}

			time.Sleep(sleepDuration)
		}
	}
}

func TestStepWellOverflow() {
	numCores := uint64(16)
	capacity := uint64(10)
	refillRate := float64(200)
	bucketType := 1
	duration := 60 * time.Second
	// Can be adapted to only overflow a certain set of cores
	coreIDs := []uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	numActiveCores := 16

	stopChans := make([]chan struct{}, numActiveCores)
	stopChanMeasurement := make(chan struct{})

	stepwell := stepwell.NewStepwell(numCores, time.Now(), bucketType, capacity, refillRate)

	for i, coreID := range coreIDs {
		stopChans[i] = make(chan struct{})
		go handleCoreRequests(stepwell, coreID, stopChans[i])
	}

	go measureTokenAmount(stepwell, stopChanMeasurement)

	time.Sleep(duration)

	for _, stopChan := range stopChans {
		close(stopChan)
	}
	close(stopChanMeasurement)

	// Wait a bit for goroutines to clean up before ending the test
	time.Sleep(1 * time.Second)
	fmt.Println("Test completed.")
}

func handleBucketRequests(tokenbucket tokenbucket.TokenBucketInterface, coreID uint64, stopChan <-chan struct{}) {
	num_allowed := 0
	for {
		select {
		case <-stopChan: // Stop signal received
			fmt.Printf("Stopping requests for core %d\n", coreID)
			return
		default:
			rand.Seed(time.Now().UnixNano())
			randomDuration := time.Duration(rand.Intn(2)) * time.Nanosecond
			// Sleep for the random duration
			time.Sleep(randomDuration)
			requestTime := time.Now()
			allowed := tokenbucket.IsAllowed(1, requestTime)
			if allowed {
				num_allowed++
			}
		}
	}
}

func measureTokenAmountTokenbucket(tokenbucket tokenbucket.TokenBucketInterface, stopChan <-chan struct{}) {
	sleepDuration := time.Duration(5) * time.Nanosecond
	for {
		select {
		case <-stopChan: // Stop signal received
			fmt.Println("Stopping token measurement.")
			return
		default:
			tokens := tokenbucket.GetTokens()
			capacity := tokenbucket.GetCapacity()
			// Don't yet manage to "overflow" the token bucket without lock / atomic operations
			if tokens > 8 {
				fmt.Printf("Capacity of TokenBucket: %d, Actual Tokens: %d\n",
					capacity,
					tokens)
			}
			time.Sleep(sleepDuration)
		}
	}
}

func TestTokenBucketOverflow() {
	numCores := uint64(100)
	capacity := uint64(10)
	refillRate := float64(10000)
	bucketType := 1
	duration := 60 * time.Second

	tokenbucket := tokenbucket.NewTokenBucketByType(bucketType, capacity, refillRate, time.Now())

	stopChans := make([]chan struct{}, numCores)
	stopChanMeasurement := make(chan struct{})

	for i := 0; i < int(numCores); i++ {
		stopChans[i] = make(chan struct{})
		go handleBucketRequests(tokenbucket, uint64(i), stopChans[i])
	}

	go measureTokenAmountTokenbucket(tokenbucket, stopChanMeasurement)

	time.Sleep(duration)

	for _, stopChan := range stopChans {
		close(stopChan)
	}
	close(stopChanMeasurement)

	// Wait a bit for goroutines to clean up before ending the test
	time.Sleep(1 * time.Second)
	fmt.Println("Test completed.")
}

// generateRequestTimings generates a slice of time.Duration for request timings.
func GenerateRequestTimings(startDelayMS, intervalMS, count int) []time.Duration {
	var timings []time.Duration
	for i := 0; i < count; i++ {
		timings = append(timings, time.Duration(startDelayMS+i*intervalMS)*time.Millisecond)
	}
	return timings
}

// testTokenBucket conducts a test for token bucket implementations.
func TestTokenBucket(bucketType int, capacity uint64, refillRate float64, requestTimings []time.Duration) {
	now := time.Now()
	tb := tokenbucket.NewTokenBucketByType(bucketType, capacity, refillRate, now)
	for _, duration := range requestTimings {
		requestTime := now.Add(duration)
		allowed := tb.IsAllowed(1, requestTime)
		fmt.Printf("TokenBucket: Request at %s allowed: %v\n", requestTime.Format(time.RFC3339), allowed)
	}
}
