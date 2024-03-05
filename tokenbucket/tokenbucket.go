package tokenbucket

import(
	"time"
)

type TokenBucketInterface interface {
    refillTokens(now time.Time)
    IsAllowed(amount uint64, now time.Time) bool
}

type TokenBucket struct {
	capacity uint64
	tokens	uint64
	refillRate float64
	// Store as Unix timestamp to be able to use atomic operations
	lastRefill int64
}