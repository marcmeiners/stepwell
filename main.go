package main

import (
    "fmt"
    "time"
    "stepwell/tokenbucket"
)

func main() {
    tb := tokenbucket.newTokenBucketLock(10, 1, time.Now())

    now := time.Now()
    requestTimes := []time.Duration{0, time.Second, 2 * time.Second, 3 * time.Second, 5 * time.Second}

    for _, duration := range requestTimes {
        requestTime := now.Add(duration)
        allowed := tb.IsAllowed(1, requestTime)
        fmt.Printf("Request at %s allowed: %v\n", requestTime.Format(time.RFC3339), allowed)
    }
}