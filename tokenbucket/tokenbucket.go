package tokenbucket

import (
	"time"
)

type TokenBucketInterface interface {
	IsAllowed(amount int64, now time.Time) bool
	GetCapacity() int64
	GetTokens() int64
}

func NewTokenBucketByType(bucketType int, capacity int64, refillRate float64, now time.Time) TokenBucketInterface {
	switch bucketType {
	case 1:
		return NewTokenBucketTrivial(capacity, refillRate, now)
	case 2:
		return NewTokenBucketAtomicLoops(capacity, refillRate, now)
	case 3:
		return NewTokenBucketLock(capacity, refillRate, now)
	case 4:
		return NewTokenBucketHelia(capacity, refillRate, now)
	default:
		return NewTokenBucketTrivial(capacity, refillRate, now)
	}
}
