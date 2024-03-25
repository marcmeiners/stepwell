//Inspired by: https://doi.org/10.1145/3548606.3560582

package tokenbucket

import (
	"time"
)

type TokenBucketHelia struct {
	capacity uint64
	//refill rate: packets/second -> take inverse to avoid further divisions in IsAllowed()
	refillRateInverse float64
	timestamp         time.Time
}

func NewTokenBucketHelia(capacity uint64, refillRate float64, timestamp time.Time) *TokenBucketHelia {
	return &TokenBucketHelia{
		capacity:          capacity,
		refillRateInverse: 1 / refillRate,
		timestamp:         timestamp,
	}
}

func (bucket *TokenBucketHelia) IsAllowed(amount uint64, now time.Time) bool {
	T := time.Duration(float64(bucket.capacity) * bucket.refillRateInverse * float64(time.Second))
	packetTime := time.Duration(float64(amount) * bucket.refillRateInverse * float64(time.Second))

	latestTimestamp := bucket.timestamp
	if now.After(bucket.timestamp) {
		latestTimestamp = now
	}

	if !latestTimestamp.Add(packetTime).After(now.Add(T)) {
		bucket.timestamp = latestTimestamp.Add(packetTime)
		return true
	} else {
		return false
	}
}

var _ TokenBucketInterface = (*TokenBucketHelia)(nil)
