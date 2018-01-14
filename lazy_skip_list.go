package lazyskiplist

import (
	"math/rand"
	"sync"
	"time"

	"github.com/petermattis/goid"
)

func init() {
	rand.Seed(time.Now().Unix())
}

var (
	// MaxHeight 128
	MaxHeight = 128
)

// left sentinal
type lSentinal struct{}

// right setinal
type rSentinal struct{}

// Node is list node
type Node struct {
	topLayer    int
	fullyLinked bool
	removed     bool
	lock        sync.Mutex
	nexts       []*Node
	Value       interface{}
}

// LazySkipList is the list structure for the algorithm
type LazySkipList struct {
	head  *Node
	less  func(v1, v2 interface{}) bool
	equal func(v1, v2 interface{}) bool
}

// New receives a less function to help values sorted
func New(less func(v1, v2 interface{}) bool, equal ...func(v1, v2 interface{}) bool) *LazySkipList {
	h := &Node{
		topLayer:    MaxHeight,
		fullyLinked: true,
		nexts:       make([]*Node, MaxHeight),
		Value:       lSentinal{},
	}
	t := &Node{
		topLayer:    MaxHeight,
		fullyLinked: true,
		Value:       rSentinal{},
	}
	for i := range h.nexts {
		h.nexts[i] = t
	}
	l := &LazySkipList{
		head: h,
	}
	l.less = func(v1, v2 interface{}) bool {
		if _, ok := v1.(lSentinal); ok {
			return true
		}
		if _, ok := v1.(rSentinal); ok {
			return false
		}
		if _, ok := v2.(rSentinal); ok {
			return true
		}
		return less(v1, v2)
	}

	if len(equal) != 0 {
		l.equal = func(v1, v2 interface{}) bool {
			if _, ok := v1.(lSentinal); ok {
				return false
			}
			if _, ok := v1.(rSentinal); ok {
				return false
			}
			if _, ok := v2.(rSentinal); ok {
				return false
			}
			return equal[0](v1, v2)
		}
	}
	return l
}

func (l *LazySkipList) findNode(v interface{}, preds []*Node, succs []*Node) (found int) {
	pred := l.head
	found = -1
	for layer := MaxHeight - 1; layer >= 0; layer-- {
		curr := pred.nexts[layer]
		debugf("[%d/findNode] scan value %#v at layer %d, topLayer is %d", goid.Get(), curr.Value, layer, curr.topLayer)
		for l.less(curr.Value, v) {
			pred = curr
			curr = pred.nexts[layer]
		}
		// TODO: customize equal
		if found == -1 && l.equal(v, curr.Value) {
			debugf("[%d/findNode] find value %#v at layer %d", goid.Get(), v, layer)
			found = layer
		}
		preds[layer] = pred
		succs[layer] = curr
	}
	return
}

func randomLevel(maxHeight int) int {
	return rand.Intn(maxHeight)
}

// Add adds element in list
func (l *LazySkipList) Add(v interface{}) {
	topLayer := randomLevel(MaxHeight)
	debugf("[%d/Add] adding value %#v, topLayer %d", goid.Get(), v, topLayer)
	preds := make([]*Node, MaxHeight)
	succs := make([]*Node, MaxHeight)
	for {
		found := l.findNode(v, preds, succs)
		if found != -1 {
			nodeFound := succs[found]
			if !nodeFound.removed {
				debugf("[%d/Add] found value %#v, waiting for node fully linked", goid.Get(), v)
				// Wait fullylinked marked
				for !nodeFound.fullyLinked {
				}
				debugf("[%d/Add] found value %#v, node fully linked, return", goid.Get(), v)
				return
			}
			debugf("[%d/Add] found value %#v at layer %d, but removed, try again", goid.Get(), v, found)
			if debug {
				// time.Sleep(time.Second)
			}
			continue
		}
		highestLocked := -1
		var pred, succ, prevPred *Node
		valid := true
		for layer := 0; valid && (layer <= topLayer); layer++ {
			pred = preds[layer]
			succ = succs[layer]
			if pred != prevPred {
				debugf("[%d/Add] lock node with value %#v at layer %d", goid.Get(), pred.Value, pred.topLayer)
				pred.lock.Lock()
				highestLocked = layer
				prevPred = pred
			}
			valid = !pred.removed && !succ.removed && pred.nexts[layer] == succ
		}
		if !valid {
			unlock(preds, highestLocked)
			continue
		}
		newNode := &Node{Value: v, topLayer: topLayer, nexts: make([]*Node, topLayer+1)}
		for layer := 0; layer <= topLayer; layer++ {
			newNode.nexts[layer] = succs[layer]
			preds[layer].nexts[layer] = newNode
		}
		newNode.fullyLinked = true
		debugf("[%d/Add] value %#v, topLayer %d, is added", goid.Get(), v, topLayer)
		unlock(preds, highestLocked)
		return
	}
}

func okToDelete(candidate *Node, l int) bool {
	return candidate.fullyLinked && candidate.topLayer == l && !candidate.removed
}

// Remove removes a element in list
func (l *LazySkipList) Remove(v interface{}) {
	debugf("[%d/Remove] removing value %#v", goid.Get(), v)
	var nodeToDelete *Node
	isRemoved := false
	topLayer := -1
	preds := make([]*Node, MaxHeight)
	succs := make([]*Node, MaxHeight)
	for {
		found := l.findNode(v, preds, succs)
		if isRemoved || found != -1 && okToDelete(succs[found], found) {
			if !isRemoved {
				debugf("[%d/Remove] removing value %#v, isRemoved=%t", goid.Get(), v, isRemoved)
				nodeToDelete = succs[found]
				topLayer = nodeToDelete.topLayer
				debugf("[%d/Remove] lock node with value %#v, at layer %d, isRemoved=%t", goid.Get(), nodeToDelete.Value, nodeToDelete.topLayer, isRemoved)
				nodeToDelete.lock.Lock()
				if nodeToDelete.removed {
					debugf("[%d/Remove] value %#v, at layer %d, has been removed, unlock it", goid.Get(), v, nodeToDelete.topLayer)
					nodeToDelete.lock.Unlock()
					return
				}
				debugf("[%d/Remove] logically remove value %#v, at layer %d, isRemoved=%t", goid.Get(), nodeToDelete.Value, nodeToDelete.topLayer, isRemoved)
				nodeToDelete.removed = true // logically removed
				isRemoved = true
			}
			highestLocked := -1
			var pred, succ, prevPred *Node
			valid := true
			for layer := 0; valid && (layer <= topLayer); layer++ {
				pred = preds[layer]
				succ = succs[layer]
				if pred != prevPred {
					debugf("[%d/Remove] lock node with value %#v at layer %d", goid.Get(), pred.Value, pred.topLayer)
					pred.lock.Lock() // [2342]
					highestLocked = layer
					prevPred = pred
				}
				valid = !pred.removed && pred.nexts[layer] == succ
			}
			if !valid {
				debugf("[%d/Remove] removing value %#v invalid, at layer %d, retry", goid.Get(), v, nodeToDelete.topLayer)
				unlock(preds, highestLocked)
				continue
			}
			for layer := topLayer; layer >= 0; layer-- {
				preds[layer].nexts[layer] = nodeToDelete.nexts[layer]
			}
			debugf("[%d/Remove] value %#v is physically removed", goid.Get(), v)
			unlock(preds, highestLocked)
			debugf("[%d/Remove] unlock node with %#v, at layer %d", goid.Get(), v, nodeToDelete.topLayer)
			nodeToDelete.lock.Unlock()
			return
		} else {
			debugf("[%d/Remove] value %#v has been removed", goid.Get(), v)
			return
		}
	}
}

func unlock(preds []*Node, highestLocked int) {
	var prevPred, pred *Node
	for layer := 0; layer <= highestLocked; layer++ {
		pred = preds[layer]
		if pred != prevPred {
			debugf("[%d/unlock] unlock node with value %#v at layer %d", goid.Get(), pred.Value, pred.topLayer)
			preds[layer].lock.Unlock()
		}
		prevPred = pred
	}
}

// Contains is a test function for v
func (l *LazySkipList) Contains(v interface{}) bool {
	preds := make([]*Node, MaxHeight)
	succs := make([]*Node, MaxHeight)
	found := l.findNode(v, preds, succs)
	return found != -1 && succs[found].fullyLinked && !succs[found].removed
}

// Iterator is used to iterate the list
type Iterator struct {
	curr *Node
}

// Next returns value until cont is false
func (i *Iterator) Next() (value interface{}, cont bool) {
	if _, ok := i.curr.Value.(rSentinal); ok {
		return nil, false
	}
	value = i.curr.Value
	i.curr = i.curr.nexts[0]
	return value, true
}

// Iterator returns a Iterator
func (l *LazySkipList) Iterator() Iterator {
	return Iterator{
		curr: l.head.nexts[0],
	}
}
