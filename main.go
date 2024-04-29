package main

import (
	"fmt"
	"os"
	"stepwell/test"
	"strconv"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <testType> <numCores>")
		os.Exit(1)
	}

	testType := os.Args[1]
	//convert number of cores to integer
	numCoresArg, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Invalid number of cores:", os.Args[2])
		os.Exit(1)
	}

	numCores := uint64(numCoresArg)

	switch testType {
	case "testStepWellLoad":
		test.TestStepWellLoad(numCores)
	case "TestTokenBucketLoad":
		test.TestTokenBucketLoad(numCores)
	default:
		fmt.Printf("Unknown test type: %s\n", testType)
		os.Exit(1)
	}
}
