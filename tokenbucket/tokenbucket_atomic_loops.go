package tokenbucket

import (
	"math"
	"sync/atomic"
	"time"
)

type TokenBucketAtomicLoops struct {
	capacity   int64
	tokens     int64
	refillRate float64
	// Store as Unix timestamp to be able to use atomic operations
	lastRefill int64
}

func NewTokenBucketAtomicLoops(capacity int64, refillRate float64, lastRefill time.Time) *TokenBucketAtomicLoops {
	return &TokenBucketAtomicLoops{
		//total capacity of tokens to give out
		capacity: capacity,
		//tokens currently available
		tokens: capacity,
		//how many new tokens per second are made available
		refillRate: refillRate,
		lastRefill: lastRefill.UnixNano(),
	}
}

func (bucket *TokenBucketAtomicLoops) refillTokens(now time.Time) {
	lastRefillUnixNano := atomic.LoadInt64(&bucket.lastRefill)
	duration := now.UnixNano() - lastRefillUnixNano
	tokensToAdd := int64(math.Floor(float64(bucket.refillRate) / 1_000_000_000 * float64(duration)))

	if tokensToAdd > 0 {
		atomic.StoreInt64(&bucket.lastRefill, now.UnixNano())
		for {
			currentTokens := atomic.LoadInt64(&bucket.tokens)
			newTokens := currentTokens + tokensToAdd
			if newTokens > bucket.capacity {
				newTokens = bucket.capacity
			}
			if atomic.CompareAndSwapInt64(&bucket.tokens, currentTokens, newTokens) {
				break
			}
		}
	}
}

func (bucket *TokenBucketAtomicLoops) GetCapacity() int64 {
	return bucket.capacity
}

func (bucket *TokenBucketAtomicLoops) GetTokens() int64 {
	return bucket.tokens
}

func (bucket *TokenBucketAtomicLoops) IsAllowed(amount int64, now time.Time) bool {
	bucket.refillTokens(now)
	for {
		currentTokens := atomic.LoadInt64(&bucket.tokens)
		if currentTokens < amount {
			return false
		}
		if atomic.CompareAndSwapInt64(&bucket.tokens, currentTokens, currentTokens-amount) {
			return true
		}
	}
}

var _ TokenBucketInterface = (*TokenBucketAtomicLoops)(nil)
