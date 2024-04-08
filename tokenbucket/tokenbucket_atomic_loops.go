package tokenbucket

import (
	"sync/atomic"
	"time"
)

type TokenBucketAtomicLoops struct {
	capacity   uint64
	tokens     uint64
	refillRate float64
	// Store as Unix timestamp to be able to use atomic operations
	lastRefill int64
}

func NewTokenBucketAtomicLoops(capacity uint64, refillRate float64, lastRefill time.Time) *TokenBucketAtomicLoops {
	return &TokenBucketAtomicLoops{
		//total capacity of tokens to give out
		capacity: capacity,
		//tokens currently available
		tokens: capacity,
		//how many new tokens per second are made available
		refillRate: refillRate,
		lastRefill: lastRefill.Unix(),
	}
}

func (bucket *TokenBucketAtomicLoops) refillTokens(now time.Time) {
	lastRefillUnix := atomic.LoadInt64(&bucket.lastRefill)
	duration := now.Unix() - lastRefillUnix
	tokensToAdd := uint64(bucket.refillRate * float64(duration))

	if tokensToAdd > 0 {
		atomic.StoreInt64(&bucket.lastRefill, now.Unix())
		for {
			currentTokens := atomic.LoadUint64(&bucket.tokens)
			newTokens := currentTokens + tokensToAdd
			if newTokens > bucket.capacity {
				newTokens = bucket.capacity
			}
			if atomic.CompareAndSwapUint64(&bucket.tokens, currentTokens, newTokens) {
				break
			}
		}
	}
}

func (bucket *TokenBucketAtomicLoops) GetCapacity() uint64 {
	return bucket.capacity
}

func (bucket *TokenBucketAtomicLoops) GetTokens() uint64 {
	return bucket.tokens
}

func (bucket *TokenBucketAtomicLoops) IsAllowed(amount uint64, now time.Time) bool {
	bucket.refillTokens(now)
	for {
		currentTokens := atomic.LoadUint64(&bucket.tokens)
		if currentTokens < amount {
			return false
		}
		if atomic.CompareAndSwapUint64(&bucket.tokens, currentTokens, currentTokens-amount) {
			return true
		}
	}
}

var _ TokenBucketInterface = (*TokenBucketAtomicLoops)(nil)
