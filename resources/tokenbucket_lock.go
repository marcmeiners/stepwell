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

func new(capacity uint64, refillPeriod time.Duration){
	return &TokenBucketLock{
		//total capacity of tokens to give out
		capacity: capacity,
		//tokens currently available
		tokens: capacity,
		refillPeriod: refillPeriod,
		lastRefill: time.Now()
	}
}

func (bucket *TokenBucketLock) refillTokens(){
	tokensToAdd := uint64(time.Since(bucket.lastRefill).Milliseconds() / bucket.refillPeriod.Milliseconds())
	
	if(tokensToAdd > 0){
		bucket.lastRefill = time.Now()
		newTokens := bucket.tokens + tokensToAdd
		if newTokens > bucket.capacity {
			newTokens = bucket.capacity
		}
		bucket.tokens = newTokens
	}
}

func (bucket *TokenBucketLock) isAllowed() {
	bucket.Lock()
	//Defer: Hold the lock and immediately release it before returning
	defer bucket.Unlock()
	bucket.refillTokens()
	if(bucket.tokens > 0){
		bucket.tokens--
		return true
	}
	return false
}

