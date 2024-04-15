//Tokenbucket with mutex
//Inspired by similar Java Implemenation https://www.codereliant.io/rate-limiting-deep-dive/

package tokenbucket

import (
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
		lastRefill: lastRefill.Unix(),
	}
}

func (bucket *TokenBucketTrivial) refillTokens(now time.Time) {
	nowUnix := now.Unix()
	duration := nowUnix - bucket.lastRefill
	tokensToAdd := int64(bucket.refillRate * float64(duration))

	if tokensToAdd > 0 {
		bucket.lastRefill = now.Unix()
		newTokens := bucket.tokens + tokensToAdd
		if newTokens > bucket.capacity {
			newTokens = bucket.capacity
		}
		bucket.tokens = newTokens
	}
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
