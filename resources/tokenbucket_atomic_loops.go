import(
	"time"
	"sync.atomic"
)

type TokenBucketAtomicLoops struct {
	TokenBucket
}

func new(capacity uint64, refillPeriod time.Duration){
	return &TokenBucketAtomic{
		//total capacity of tokens to give out
		capacity: capacity,
		//tokens currently available
		tokens: capacity,
		refillPeriod: refillPeriod,
		lastRefill: time.Now()
	}
}

func (bucket *TokenBucketAtomicLoops) refillTokens(){
	tokensToAdd := uint64(time.Since(atomic.LoadInt64(&bucket.lastRefill)).Milliseconds() / bucket.refillPeriod.Milliseconds())
	
	if(tokensToAdd > 0){
		atomic.StoreInt64(&bucket.lastRefill, time.Now())	
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

func (bucket *TokenBucketAtomicLoops) isAllowed() {
	bucket.refillTokens()
	for {
		currentTokens := atomic.LoadInt64(&bucket.tokens)
		if currentTokens <= 0 { //use "<=" because bucket could slightly overflow without using the mutex
			return false
		}
		if atomic.CompareAndSwapInt64(&bucket.tokens, currentTokens, currentTokens-1) {
			return true
		}
	}
}

var _ TokenBucket = (*TokenBucketAtomicLoops)(nil)