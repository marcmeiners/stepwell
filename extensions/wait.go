package extensions

import (
	"math/rand"
	"time"
)

func RandomWait() {
	rand.Seed(time.Now().UnixNano())
	randomDuration := time.Duration(rand.Intn(5)) * time.Nanosecond
	// Sleep for the random duration
	time.Sleep(randomDuration)
}

func ShortWait() {
	nanosecond := time.Duration(1) * time.Nanosecond
	time.Sleep(nanosecond)
}
