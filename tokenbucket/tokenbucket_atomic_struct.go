package tokenbucket

import (
	"math"
	"sync/atomic"
	"time"
	"unsafe"
)

type tokenBucketContents struct {
	tokens     int64
	lastRefill int64
}
type TokenBucketAtomicStructs struct {
	capacity   int64
	contents   *tokenBucketContents
	refillRate float64
}

func NewTokenBucketAtomicStructs(capacity int64, refillRate float64,
	lastRefill time.Time) *TokenBucketAtomicStructs {
	return &TokenBucketAtomicStructs{
		//total capacity of tokens to give out
		capacity: capacity,
		//tokens currently available
		contents: &tokenBucketContents{tokens: capacity, lastRefill: lastRefill.Unix()},
		//how many new tokens per second are made available
		refillRate: refillRate,
	}
}

func (bucket *TokenBucketAtomicStructs) refillTokens(now time.Time) {
	for {
		lastContents := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&bucket.
			contents)))
		contents := (*tokenBucketContents)(lastContents)
		duration := now.UnixNano() - contents.lastRefill
		tokensToAdd := int64(math.Floor(float64(bucket.refillRate) / 1_000_000_000 * float64(
			duration)))

		if tokensToAdd > 0 {
			newTokens := contents.tokens + tokensToAdd
			if newTokens > bucket.capacity {
				newTokens = bucket.capacity
			}
			newStruct := tokenBucketContents{
				tokens:     newTokens,
				lastRefill: now.UnixNano(),
			}
			if atomic.CompareAndSwapPointer(
				(*unsafe.Pointer)(unsafe.Pointer(&bucket.contents)), lastContents,
				unsafe.Pointer(&newStruct)) {
				break
			}
		} else {
			break
		}
	}
}

func (bucket *TokenBucketAtomicStructs) GetCapacity() int64 {
	return bucket.capacity
}

func (bucket *TokenBucketAtomicStructs) GetTokens() int64 {
	return bucket.contents.tokens
}

func (bucket *TokenBucketAtomicStructs) IsAllowed(amount int64, now time.Time) bool {
	bucket.refillTokens(now)
	for {
		lastContents := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(
			&bucket.contents)))
		contents := (*tokenBucketContents)(lastContents)
		currentTokens := contents.tokens
		if currentTokens < amount {
			return false
		}
		newStruct := tokenBucketContents{
			tokens:     currentTokens - amount,
			lastRefill: contents.lastRefill,
		}

		if atomic.CompareAndSwapPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&bucket.contents)), lastContents,
			unsafe.Pointer(&newStruct)) {
			return true
		}
	}
}

var _ TokenBucketInterface = (*TokenBucketAtomicStructs)(nil)
