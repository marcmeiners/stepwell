//Inspired by: https://doi.org/10.1145/3548606.3560582

package tokenbucket

import (
	"sync/atomic"
	"time"
)

type TokenBucketHelia struct {
	capacity int64
	//refill rate: packets/second -> take inverse to avoid further divisions in IsAllowed()
	refillRateInverse float64
	timestamp         int64
}

func NewTokenBucketHelia(capacity int64, refillRate float64, timestamp time.Time) *TokenBucketHelia {
	return &TokenBucketHelia{
		capacity:          capacity,
		refillRateInverse: 1 / refillRate,
		timestamp:         timestamp.UnixNano(),
	}
}

func (bucket *TokenBucketHelia) GetCapacity() int64 {
	return bucket.capacity
}

func (bucket *TokenBucketHelia) GetTokens() int64 {
	now := time.Now()
	nowUnix := now.UnixNano()
	latestTimestamp := atomic.LoadInt64(&bucket.timestamp)
	if nowUnix >= latestTimestamp {
		return 0
	} else {
		duration := time.Duration(latestTimestamp - nowUnix)
		durationInSeconds := float64(duration) / float64(time.Second)
		return int64(durationInSeconds / bucket.refillRateInverse)
	}
}

// time.Duration is a type having int64 as its underlying type, which stores the duration in nanoseconds.
func (bucket *TokenBucketHelia) IsAllowed(amount int64, now time.Time) bool {
	T := time.Duration(float64(bucket.capacity) * bucket.refillRateInverse * float64(time.Second))
	packetTime := time.Duration(float64(amount) * bucket.refillRateInverse * float64(time.Second))

	nowUnix := now.UnixNano()
	for {
		latestTimestamp := atomic.LoadInt64(&bucket.timestamp)
		newTimestamp := int64(0)

		if nowUnix > latestTimestamp {
			newTimestamp = nowUnix + int64(packetTime)
		} else {
			newTimestamp = latestTimestamp + int64(packetTime)
		}

		if newTimestamp > nowUnix+int64(T) {
			return false
		}

		if atomic.CompareAndSwapInt64(&bucket.timestamp, latestTimestamp, newTimestamp) {
			return true
		}
	}
}

var _ TokenBucketInterface = (*TokenBucketHelia)(nil)
