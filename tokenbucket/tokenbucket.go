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

func NewTokenBucketByType(type int, capacity uint64, refillRate float64, lastRefill time.Time){
	switch type {
		case 1:
			return NewTokenBucketTrivial(capacity, refillRate, lastRefill)
		case 2:
			return NewTokenBucketAtomicLoops(capacity, refillRate, lastRefill)
		case 3:
			return NewTokenBucketLock(capacity, refillRate, lastRefill)
		default:
			return NewTokenBucketTrivial(capacity, refillRate, lastRefill)
	}
}