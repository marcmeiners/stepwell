//Tokenbucket with mutex
//Inspired by similar Java Implemenation https://www.codereliant.io/rate-limiting-deep-dive/

package tokenbucket

import (
	"math"
	"stepwell/extensions"
	"time"
)

type TokenBucketTrivial struct {
	capacity   int64
	tokens     int64
	refillRate float64
	// Store as Unix timestamp to be able to use atomic operations
	lastRefill int64
}

func NewTokenBucketTrivial(capacity int64, refillRate float64, lastRefill time.Time) *TokenBucketTrivial {
	return &TokenBucketTrivial{
		//total capacity of tokens to give out
		capacity: capacity,
		//tokens currently available
		tokens: capacity,
		//how many new tokens per second are made available
		refillRate: refillRate,
		lastRefill: lastRefill.UnixNano(),
	}
}

func (bucket *TokenBucketTrivial) refillTokens(now time.Time) {
	nowUnixNano := now.UnixNano()
	duration := nowUnixNano - bucket.lastRefill
	tokensToAdd := int64(math.Floor(float64(bucket.refillRate) / 1_000_000_000 * float64(duration)))

	if tokensToAdd > 0 {
		bucket.lastRefill = now.UnixNano()
		newTokens := bucket.tokens + tokensToAdd
		if newTokens > bucket.capacity {
			newTokens = bucket.capacity
		}
		bucket.tokens = newTokens
	}
}

func (bucket *TokenBucketTrivial) SetRefillRate(refillRate float64) {
	bucket.refillRate = refillRate
}

func (bucket *TokenBucketTrivial) GetCapacity() int64 {
	return bucket.capacity
}

func (bucket *TokenBucketTrivial) GetTokens() int64 {
	return bucket.tokens
}

func (bucket *TokenBucketTrivial) IsAllowed(amount int64, now time.Time) bool {
	bucket.refillTokens(now)
	if bucket.tokens >= amount {
		//Wait a few nanoseconds to show concurrency effect
		extensions.ShortWait()
		bucket.tokens -= amount
		return true
	}
	return false
}

var _ TokenBucketInterface = (*TokenBucketTrivial)(nil)
