package stepwellplus

import (
	"stepwell/tokenbucket"
	"time"
)

//TODO: add dynamic resource share allocator core

type StepWellPlusInterface interface {
	IsAllowed(port uint64, amount int64, now time.Time) bool
}

type StepWellPlus struct {
	Cores      []*StepWellPlusNode
	numCores   uint64
	Capacity   int64
	refillRate float64
	bucketType int
}

type StepWellPlusNode struct {
	TokenBucket tokenbucket.TokenBucketInterface
}

func NewStepwellPlus(numCores uint64, now time.Time, bucketType int, capacity int64, refillRate float64) *StepWellPlus {
	if numCores <= 0 {
		return nil
	}

	var cores []*StepWellPlusNode

	for i := uint64(0); i < numCores; i++ {
		node := &StepWellPlusNode{TokenBucket: tokenbucket.NewTokenBucketByType(bucketType, capacity/int64(numCores), refillRate/float64(numCores), now)}
		cores = append(cores, node)
	}

	return &StepWellPlus{
		Cores:      cores,
		numCores:   numCores,
		Capacity:   capacity,
		refillRate: refillRate,
		bucketType: bucketType,
	}
}

func (stepwellplus *StepWellPlus) IsAllowed(port uint64, amount int64, now time.Time) bool {
	return stepwellplus.Cores[port].TokenBucket.IsAllowed(amount, now)
}

var _ StepWellPlusInterface = (*StepWellPlus)(nil)
