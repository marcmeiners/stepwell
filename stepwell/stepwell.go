package stepwell

import (
	"stepwell/tokenbucket"
	"time"
)

//Problems to think about:

//  1. What about the case when a tokenBucket on the path returns false when calling IsAllowed()?
//     Should the tokens issued before be revoked again?
type StepWellInterface interface {
	//Get a token for all the buckets on the path to the single bucket which is the root of the tree and the bottom layer of the StepWell tree structure
	IsAllowed(port uint64, amount int64, now time.Time) bool
}

type StepWell struct {
	//Cores are the top layer token buckets or the leaves in the StepWell tree structure
	Cores      []*StepWellNode
	root       *StepWellNode
	numCores   uint64
	Capacity   int64
	refillRate float64
	bucketType int
}

// idea: use a tree structure similar to a linked list
// later on we use the tree flipped aroung, i.e. we start at the leaves and iteratively use "prev"
// also: we make it possible to use different tokenbucket types in the stepwell tree structure
type StepWellNode struct {
	TokenBucket tokenbucket.TokenBucketInterface
	Parent      *StepWellNode
	leftChild   *StepWellNode
	rightChild  *StepWellNode
}

func NewStepwell(numCores uint64, now time.Time, bucketType int, capacity int64, refillRate float64) *StepWell {
	if numCores <= 0 {
		return nil
	}

	var nodes []*StepWellNode
	root := &StepWellNode{TokenBucket: tokenbucket.NewTokenBucketByType(bucketType, capacity, refillRate, now)}
	nodes = append(nodes, root)

	for levelCount := uint64(1); levelCount < numCores; {
		levelCount = 0
		var nextLevel []*StepWellNode
		for _, node := range nodes {
			if levelCount >= numCores {
				break
			}
			leftChild := &StepWellNode{TokenBucket: tokenbucket.NewTokenBucketByType(bucketType, capacity, refillRate, now), Parent: node}
			node.leftChild = leftChild
			nextLevel = append(nextLevel, leftChild)
			levelCount++

			if levelCount >= numCores {
				break
			}
			rightChild := &StepWellNode{TokenBucket: tokenbucket.NewTokenBucketByType(bucketType, capacity, refillRate, now), Parent: node}
			node.rightChild = rightChild
			nextLevel = append(nextLevel, rightChild)
			levelCount++
		}
		nodes = nextLevel
	}

	return &StepWell{
		Cores:      nodes,
		root:       root,
		numCores:   numCores,
		Capacity:   capacity,
		refillRate: refillRate,
		bucketType: bucketType,
	}
}

func (stepwell *StepWell) IsAllowed(port uint64, amount int64, now time.Time) bool {
	var curr *StepWellNode = stepwell.Cores[port]

	if !curr.TokenBucket.IsAllowed(amount, now) {
		return false
	}

	for curr.Parent != nil {
		curr = curr.Parent
		if !curr.TokenBucket.IsAllowed(amount, now) {
			return false
		}
	}
	return true
}

var _ StepWellInterface = (*StepWell)(nil)
