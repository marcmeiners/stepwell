//Tokenbucket with mutex
//Inspired by similar Java Implemenation https://www.codereliant.io/rate-limiting-deep-dive/

import(
	"time"
	"sync"
)

type TokenBucket struct {
	//https://stackoverflow.com/questions/44949467/when-do-you-embed-mutex-in-struct-in-go
	sync.Mutex
	capacity uint64
	tokens	uint64
	refillPeriod time.Duration
	lastRefill time.Time
}

func new(capacity uint64, refillPeriod time.Duration){
	return &TokenBucket{
		//total capacity of tokens to give out
		capacity: capacity,
		//tokens currently available
		tokens: capacity,
		refillPeriod: refillPeriod,
		lastRefill: time.Now()
	}
}

func (bucket *TokenBucket) refillTokens(){
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

func (bucket *TokenBucket) isAllowed() {
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


