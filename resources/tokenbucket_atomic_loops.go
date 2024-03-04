import(
	"time"
	"sync.atomic"
)

type TokenBucketAtomicLoops struct {
	TokenBucket
}

func new(capacity uint64, refillRate float64, lastRefill time.Time){
	return &TokenBucketAtomicLoops{
		//total capacity of tokens to give out
		capacity: capacity
		//tokens currently available
		tokens: capacity
		//how many new tokens per second are made available
		refillRate: refillRate
		lastRefill: lastRefill
	}
}

func (bucket *TokenBucketAtomicLoops) refillTokens(now time.Time){
	duration := now.Sub(atomic.LoadInt64(bucket.lastRefill))
	tokensToAdd := bucket.refillRate * duration.Seconds()
	
	if(tokensToAdd > 0){
		atomic.StoreInt64(&bucket.lastRefill, now())	
		for(true) {
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

func (bucket *TokenBucketAtomicLoops) isAllowed(amount uint64, now time.Time) {
	bucket.refillTokens()
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

var _ TokenBucket = (*TokenBucketAtomicLoops)(nil)