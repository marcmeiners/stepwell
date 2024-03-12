package core

import (
    "fmt"
    "time"
    "stepwell/tokenbucket"
)

//Problems to think about:

//1. What about the case when a tokenBucket on the path returns false when calling IsAllowed()? 
//   Should the tokens issued before be revoked again?
//2. How to "glue" the buckets together in an efficient way such that the isAllowed() function of stepwell is fast?
//   Is a linked tree structure the right way?

//"Ports" are the top layer token buckets of the stepwell construct
type StepWellInterface interface {
    Build(numPorts uint64, now time.Time)
    //Get a token for all the buckets on the path to the single bucket which is the root of the tree and the bottom layer of the StepWell tree structure
    IsAllowed(port uint64, amount uint64, now time.Time) bool
    //Only refill the buckets along the way in the graph starting from the corresponding port
    refillTokens(port uint64, now time.Time)
}

type StepWell struct {
    //Ports are the top layer token buckets or the leaves in the StepWell tree structure
    Ports []StepWellNode
    numPorts uint64
    //is the same for all token buckets
	capacity uint64
}

//idea: use a tree structure similar to a linked list 
//later on we use the tree flipped aroung, i.e. we start at the leaves and iteratively use "prev"
//also: we make it possible to use different tokenbucket types in the stepwell tree structure
type StepWellNode struct {
    tokenbucket TokenBucket
    prev *StepWellNode
    nextLeft *StepWellNode
    nextRight *StepWellNode
}