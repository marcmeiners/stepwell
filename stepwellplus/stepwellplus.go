package stepwellplus

import (
	"stepwell/tokenbucket"
	"sync/atomic"
	"time"
)

type StepWellPlusInterface interface {
	IsAllowed(port uint64, amount int64, now time.Time) bool
	StartWorker()
	StopWorker()
}

type StepWellPlus struct {
	Cores         []*StepWellPlusNode
	numCores      uint64
	refreshDelay  time.Duration
	workerRunning bool
	stopChan      chan struct{}
	Capacity      int64
	refillRate    float64
	bucketType    int
}

type StepWellPlusNode struct {
	TokenBucket tokenbucket.TokenBucketInterface
	requests    int64
}

func NewStepwellPlus(numCores uint64, refreshDelay time.Duration, now time.Time, bucketType int, capacity int64, refillRate float64) *StepWellPlus {
	if numCores <= 0 {
		return nil
	}

	var cores []*StepWellPlusNode

	for i := uint64(0); i < numCores; i++ {
		node := &StepWellPlusNode{TokenBucket: tokenbucket.NewTokenBucketByType(bucketType, capacity/int64(numCores), refillRate/float64(numCores), now)}
		cores = append(cores, node)
	}

	return &StepWellPlus{
		Cores:         cores,
		numCores:      numCores,
		refreshDelay:  refreshDelay,
		workerRunning: false,
		stopChan:      make(chan struct{}),
		Capacity:      capacity,
		refillRate:    refillRate,
		bucketType:    bucketType,
	}
}

func (stepwellplus *StepWellPlus) IsAllowed(port uint64, amount int64, now time.Time) bool {
	core := stepwellplus.Cores[port]
	atomic.AddInt64(&core.requests, 1)
	return core.TokenBucket.IsAllowed(amount, now)
}

func (stepwellplus *StepWellPlus) StartWorker() {
	if stepwellplus.workerRunning {
		return
	}
	stepwellplus.workerRunning = true
	stepwellplus.stopChan = make(chan struct{})

	go stepwellplus.startWorkerCore()
}

func (stepwellplus *StepWellPlus) StopWorker() {
	if !stepwellplus.workerRunning {
		return
	}
	stepwellplus.workerRunning = false

	close(stepwellplus.stopChan)
}

func (stepwellplus *StepWellPlus) startWorkerCore() {
	ticker := time.NewTicker(stepwellplus.refreshDelay)
	defer ticker.Stop()

	for {
		select {
		case <-stepwellplus.stopChan:
			return
		case <-ticker.C:
			totalRequests := int64(0)
			requestCounts := make([]int64, stepwellplus.numCores)

			for i, core := range stepwellplus.Cores {
				requests := atomic.LoadInt64(&core.requests)
				requestCounts[i] = requests
				totalRequests += requests
			}

			if totalRequests > 0 {
				for i, core := range stepwellplus.Cores {
					proportionalRate := stepwellplus.refillRate * (float64(requestCounts[i]) / float64(totalRequests))
					core.TokenBucket.SetRefillRate(proportionalRate)
					atomic.StoreInt64(&core.requests, 0)
				}
			}
		}
	}
}

var _ StepWellPlusInterface = (*StepWellPlus)(nil)
