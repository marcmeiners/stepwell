//Tokenbucket with mutex
//Inspired by similar Java Implemenation https://www.codereliant.io/rate-limiting-deep-dive/

package tokenbucket

import(
	"time"
	"sync"
)

type TokenBucketLock struct {
	//https://stackoverflow.com/questions/44949467/when-do-you-embed-mutex-in-struct-in-go
	sync.Mutex
	TokenBucket
}

func newTokenBucketLock(capacity uint64, refillRate float64, lastRefill time.Time){
	return &TokenBucketLock{
		//total capacity of tokens to give out
		capacity: capacity
		//tokens currently available
		tokens: capacity
		//how many new tokens per second are made available
		refillRate: refillRate
		lastRefill: lastRefill
	}
}

func (bucket *TokenBucketLock) refillTokens(now time.Time){
	duration := now.Sub(bucket.lastRefill)
	tokensToAdd := bucket.refillRate * duration.Seconds()
		
	if(tokensToAdd > 0){
		bucket.lastRefill = now
		newTokens := bucket.tokens + tokensToAdd
		if newTokens > bucket.capacity {
			newTokens = bucket.capacity
		}
		bucket.tokens = newTokens
	}
}

func (bucket *TokenBucketLock) isAllowed(amount uint64, now time.Time) bool {
	bucket.Lock()
	//Defer: Hold the lock and immediately release it before returning
	defer bucket.Unlock()
	bucket.refillTokens()
	if(bucket.tokens >= amount){
		bucket.tokens -= amount
		return true
	}
	return false
}

var _ TokenBucket = (*TokenBucketLock)(nil)