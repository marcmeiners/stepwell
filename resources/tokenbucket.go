package tokenbucket

import(
	"time"
)

type TokenBucket interface {
    refillTokens()
    isAllowed(amount uint64) bool
}

type TokenBucket struct {
	capacity uint64
	tokens	uint64
	refillRate float64
	lastRefill time.Time
}