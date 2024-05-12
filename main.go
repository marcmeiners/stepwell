package main

import (
	"fmt"
	"os"
	"stepwell/test"
	"strconv"
)

func main() {
	if len(os.Args) < 7 { // Ensure there are at least four arguments
		fmt.Println("Usage: go run main.go <testType> <numCores> <bucketType> <duration> <refillRate> <capacity>")
		os.Exit(1)
	}

	testType := os.Args[1]
	// Convert number of cores to integer
	numCoresArg, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Invalid number of cores:", os.Args[2])
		os.Exit(1)
	}
	numCores := uint64(numCoresArg)

	// Convert bucket type to integer
	bucketTypeArg, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("Invalid bucket type:", os.Args[3])
		os.Exit(1)
	}
	bucketType := int(bucketTypeArg)

	// Convert duration to integer
	durationArg, err := strconv.Atoi(os.Args[4])
	if err != nil {
		fmt.Println("Invalid duration:", os.Args[4])
		os.Exit(1)
	}
	duration := int(durationArg)

	// Convert refill rate to integer
	refillRateArg, err := strconv.Atoi(os.Args[5])
	if err != nil {
		fmt.Println("Invalid refill rate:", os.Args[5])
		os.Exit(1)
	}
	refillRateInt := int(refillRateArg)

	// Convert capacity to integer
	capacityArg, err := strconv.Atoi(os.Args[6])
	if err != nil {
		fmt.Println("Invalid capacity:", os.Args[6])
		os.Exit(1)
	}
	capacityInt := int(capacityArg)

	switch testType {
	case "TestStepWellLoad":
		test.TestStepWellLoad(numCores, bucketType, duration, refillRateInt, capacityInt)
	case "TestTokenBucketLoad":
		test.TestTokenBucketLoad(numCores, bucketType, duration, refillRateInt, capacityInt)
	case "TestStepWellPerformance":
		test.TestStepWellPerformance(numCores, bucketType, duration, refillRateInt, capacityInt)
	case "TestTokenBucketPerformance":
		test.TestTokenBucketPerformance(numCores, bucketType, duration, refillRateInt, capacityInt)
	default:
		fmt.Printf("Unknown test type: %s\n", testType)
		os.Exit(1)
	}
}
