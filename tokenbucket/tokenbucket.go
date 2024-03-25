package tokenbucket

import (
	"time"
)

type TokenBucketInterface interface {
	refillTokens(now time.Time)
	IsAllowed(amount uint64, now time.Time) bool
}

func NewTokenBucketByType(bucketType int, capacity uint64, refillRate float64, lastRefill time.Time) TokenBucketInterface {
	switch bucketType {
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
