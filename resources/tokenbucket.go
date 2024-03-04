package tokenbucket

import(
	"time"
)

type TokenBucket interface {
    refillTokens(now time.Time)
    isAllowed(amount uint64, now time.Time) bool
}

type TokenBucket struct {
	capacity uint64
	tokens	uint64
	refillRate float64
	lastRefill time.Time
}