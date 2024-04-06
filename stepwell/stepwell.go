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
	IsAllowed(port uint64, amount uint64, now time.Time) bool
}

type StepWell struct {
	//Cores are the top layer token buckets or the leaves in the StepWell tree structure
	cores      []*StepWellNode
	root       *StepWellNode
	numCores   uint64
	capacity   uint64
	refillRate float64
	bucketType int
}

// idea: use a tree structure similar to a linked list
// later on we use the tree flipped aroung, i.e. we start at the leaves and iteratively use "prev"
// also: we make it possible to use different tokenbucket types in the stepwell tree structure
type StepWellNode struct {
	tokenBucket tokenbucket.TokenBucketInterface
	parent      *StepWellNode
	leftChild   *StepWellNode
	rightChild  *StepWellNode
}

func NewStepwell(numCores uint64, now time.Time, bucketType int, capacity uint64, refillRate float64) *StepWell {
	if numCores <= 0 {
		return nil
	}

	var currentLevel []*StepWellNode
	root := &StepWellNode{tokenBucket: tokenbucket.NewTokenBucketByType(bucketType, capacity, refillRate, now)}
	currentLevel = append(currentLevel, root)

	var currentLeaves uint64 = 1
	var addedLeaves uint64 = 0

	for currentLeaves < numCores {
		var nextLevel []*StepWellNode
		// First pass: Add a left child to each node in the current level
		for _, node := range currentLevel {
			leftChild := &StepWellNode{tokenBucket: tokenbucket.NewTokenBucketByType(bucketType, capacity, refillRate, now), parent: node}
			node.leftChild = leftChild
			nextLevel = append(nextLevel, leftChild)
		}
		addedLeaves = currentLeaves

		// Second pass: Add a right child to nodes in the current level until reaching numCores
		for _, node := range currentLevel {
			if addedLeaves < numCores {
				rightChild := &StepWellNode{tokenBucket: tokenbucket.NewTokenBucketByType(bucketType, capacity, refillRate, now), parent: node}
				node.rightChild = rightChild
				//After a loop iteration the nextLevel array first contains all "left" children and and the all "right" children.
				//In the next iteration, if it is the last one, there might be nodes with two children and nodes with only one child.
				//Because of the order mentioned, the nodes with 2 children will be distributed over the whole tree "width". Thus we have better load balancing inside the tree.
				nextLevel = append(nextLevel, rightChild)
				addedLeaves++
			} else {
				break
			}
		}
		currentLeaves = addedLeaves
		currentLevel = nextLevel
	}

	return &StepWell{
		cores:      currentLevel,
		root:       root,
		numCores:   numCores,
		capacity:   capacity,
		refillRate: refillRate,
		bucketType: bucketType,
	}
}

func (stepwell *StepWell) IsAllowed(port uint64, amount uint64, now time.Time) bool {
	var curr *StepWellNode = stepwell.cores[port]

	if !curr.tokenBucket.IsAllowed(amount, now) {
		return false
	}

	for curr.parent != nil {
		curr = curr.parent
		if !curr.tokenBucket.IsAllowed(amount, now) {
			return false
		}
	}
	return true
}

var _ StepWellInterface = (*StepWell)(nil)
