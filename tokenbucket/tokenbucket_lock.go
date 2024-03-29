//Tokenbucket with mutex
//Inspired by similar Java Implemenation https://www.codereliant.io/rate-limiting-deep-dive/

package tokenbucket

import (
	"sync"
	"time"
)

type TokenBucketLock struct {
	capacity   uint64
	tokens     uint64
	refillRate float64
	// Store as Unix timestamp to be able to use atomic operations
	lastRefill int64
	//https://stackoverflow.com/questions/44949467/when-do-you-embed-mutex-in-struct-in-go
	sync.Mutex
}

func NewTokenBucketLock(capacity uint64, refillRate float64, lastRefill time.Time) *TokenBucketLock {
	return &TokenBucketLock{
		//total capacity of tokens to give out
		capacity: capacity,
		//tokens currently available
		tokens: capacity,
		//how many new tokens per second are made available
		refillRate: refillRate,
		lastRefill: lastRefill.Unix(),
	}
}

func (bucket *TokenBucketLock) refillTokens(now time.Time) {
	nowUnix := now.Unix()
	duration := nowUnix - bucket.lastRefill
	tokensToAdd := uint64(bucket.refillRate * float64(duration))

	if tokensToAdd > 0 {
		bucket.lastRefill = now.Unix()
		newTokens := bucket.tokens + tokensToAdd
		if newTokens > bucket.capacity {
			newTokens = bucket.capacity
		}
		bucket.tokens = newTokens
	}
}

func (bucket *TokenBucketLock) IsAllowed(amount uint64, now time.Time) bool {
	bucket.Lock()
	//Defer: Hold the lock and immediately release it before returning
	defer bucket.Unlock()
	bucket.refillTokens(now)
	if bucket.tokens >= amount {
		bucket.tokens -= amount
		return true
	}
	return false
}

var _ TokenBucketInterface = (*TokenBucketLock)(nil)
