# Stepwell

Stepwell is a scalable rate-limiting system designed for high-speed network applications. It leverages a hierarchical structure of non-locking token buckets to ensure accurate rate enforcement across multiple processing cores, minimizing synchronization overhead and enhancing throughput. This repo also contains two different atomic token bucket implementations which are faster than Stepwell and can be used on systems that support atomic operations.

## Different Rate Limiters in this Project

1. [**Baseline Token Bucket**](tokenbucket/tokenbucket_trivial.go): A simple, non-thread-safe implementation.
2. [**Locked Token Bucket**](tokenbucket/tokenbucket_lock.go): Ensures thread safety using mutex locks.
3. [**Atomic Token Bucket**](tokenbucket/tokenbucket_atomic_struct.go): Uses atomic operations to manage concurrency without locks.
4. [**Timestamp Token Bucket**](tokenbucket/tokenbucket_helia.go): An advanced atomic token bucket design storing only a single timestamp for efficient token management.
5. [**Stepwell**](stepwell/stepwell.go): A hierarchical structure of baseline token buckets that works without locking and atomic operations.

## Usage

Stepwell can be integrated into your existing Go projects. Below is an example of how to use the Stepwell system.

```go
package main

import (
	"fmt"
	"time"
	"stepwell/stepwell"
)

func main() {
	// Configuration
	numCores := uint64(1)
	bucketType := 1
	capacity := int64(100)
	refillRate := float64(10)

	// Initialize Stepwell
	stepwellSystem := stepwell.NewStepwell(numCores, time.Now(), bucketType, capacity, refillRate)

	// Simulate a single request
	requestTime := time.Now()
	allowed := stepwellSystem.IsAllowed(0, 1, requestTime)
	if allowed {
		fmt.Println("Request allowed")
	} else {
		fmt.Println("Request denied")
	}

	fmt.Println("Test completed.")
}
```

## Evaluation

Stepwell's performance has been evaluated through high-load and performance tests. The system demonstrates significant improvements in scalability and efficiency compared to traditional rate-limiting methods that incorporate locking.

## Contributing

Contributions are welcome! Please fork the repository and submit pull requests for any enhancements or bug fixes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contact

For any questions or inquiries, please open an issue on the [GitHub repository](https://github.com/marcmeiners/stepwell).