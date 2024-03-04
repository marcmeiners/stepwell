package tokenbucket

import(
	"time"
)

type TokenBucket interface {
    refillTokens()
    isAllowed() bool
}

type TokenBucket struct {
	capacity uint64
	tokens	uint64
	refillPeriod time.Duration
	lastRefill time.Time
}