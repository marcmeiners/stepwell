package core

import (
    "fmt"
    "time"
    "stepwell/tokenbucket"
)

//Problems to think about:

//1. What about the case when a tokenBucket on the path returns false when calling IsAllowed()? 
//   Should the tokens issued before be revoked again?

//"Ports" are the top layer token buckets of the stepwell construct
type StepWellInterface interface {
    //Types - 1: trivial / 2: atomic / 3: lock
    Build(numPorts uint64, now time.Time, type int64, capacity uint64)
    //Get a token for all the buckets on the path to the single bucket which is the root of the tree and the bottom layer of the StepWell tree structure
    IsAllowed(port uint64, amount uint64, now time.Time) bool
}

type StepWell struct {
    //Ports are the top layer token buckets or the leaves in the StepWell tree structure
    ports []*StepWellNode
    root *StepWellNode
    numPorts uint64
	capacity uint64
    refillRate float64
    bucketType int
}

//idea: use a tree structure similar to a linked list 
//later on we use the tree flipped aroung, i.e. we start at the leaves and iteratively use "prev"
//also: we make it possible to use different tokenbucket types in the stepwell tree structure
type StepWellNode struct {
    tokenBucket TokenBucket
	parent      *StepWellNode
	leftChild   *StepWellNode
	rightChild  *StepWellNode
}

func (stepwell *StepWell) Build(numPorts uint64, now time.Time, type int64, capacity uint64, refillRate float64) bool {
    if numPorts <= 0 {
		return false
	}

    var currentLevel []*StepWellNode
    root := &StepWellNode{tokenBucket: NewTokenBucketByType(type, capacity, refillRate, now)}
	currentLevel = append(currentLevel, root)

    var currentLeaves uint64 = 1
    var addedLeaves uint64 = 0

    for currentLeaves < numPorts {
        var nextLevel []*StepWellNode
        // First pass: Add a left child to each node in the current level
        for(index, node := range currentLevel){
            leftChild := &StepWellNode{tokenBucket: NewTokenBucketByType(type, capacity, refillRate, now), parent: node}
            node.leftChild = leftChild
            nextLevel = append(nextLevel, leftChild)
        }
        addedLeaves = currentLeaves;

        // Second pass: Add a right child to nodes in the current level until reaching numPorts
        for(index, node := range currentLevel){
            if(addedLeaves < numPorts){
                rightChild := &StepWellNode{TokenBucket: NewTokenBucketByType(type, capacity, refillRate, now), parent: node}
                node.rightChild = rightChild
                //After a loop iteration the nextLevel array first contains all "left" children and and the all "right" children. 
                //In the next iteration, if it is the last one, there might be nodes with two children and nodes with only one child.
                //Because of the order mentioned, the nodes with 2 children will be distributed over the whole tree "width". Thus we have better load balancing inside the tree.
                nextLevel = append(nextLevel, rightChild)
                addedLeaves++
            }
            else{
                break;
            }
        }
        currentLeaves = addedLeaves
        currentLevel = nextLevel
    }

    stepwell.ports = currentLevel
    stepwell.root = root
    stepwell.numPorts = numPorts
	stepwell.capacity = capacity
    stepwell.refillRate = refillRate
    stepwell.bucketType = type
    
    return true
}

func (stepwell *StepWell) IsAllowed(port uint64, amount uint64, now time.Time) bool {
    StepwellNode curr = stepwell.ports[port]

    if(!curr.tokenBucket.IsAllowed(amount, now)){
        return false
    }

    for(curr.parent != nil){
        curr = curr.parent
        if(!curr.tokenBucket.IsAllowed(amount, now)){
            return false
        }
    }
    return true
}

var _ StepWellInterface = (*StepWell)(nil)