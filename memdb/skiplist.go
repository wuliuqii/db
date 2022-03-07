package memdb

import (
	"bytes"
	"errors"
	"math/rand"
)

const (
	tMaxHeight = 12
	Branching  = 4
)

var ErrNotFound = errors.New("ErrNotFound")

type comparator interface {
	// compare returns -1, 0, or +1 depending on whether a is 'less than',
	// 'equal to' or 'greater than' b. The two arguments can only be 'equal'
	// if their contents are exactly equal. Furthermore, the empty slice
	// must be 'less than' any non-empty slice.
	compare(a, b []byte) int
}

type bytesComparator struct{}

func (bytesComparator) compare(a, b []byte) int {
	return bytes.Compare(a, b)
}

type element struct {
	levels []*element

	key   []byte
	value []byte
}

func newElement(key, value []byte, h int) *element {
	return &element{
		levels: make([]*element, h),
		key:    key,
		value:  value,
	}
}

type skipList struct {
	cmp comparator
	rnd *rand.Rand

	header    *element // dummy header
	maxHeight int
	n         int // nodes of skipList, excluded header
}

func (sl *skipList) randomHeight() int {
	h := 1
	for h < tMaxHeight && sl.rnd.Int()%Branching == 0 {
		h++
	}
	return h
}

// find node greater or equal to given key, if contained, second return args will be true
// prev recode every level previous element node
func (sl *skipList) findGE(key []byte, prev []*element) (*element, bool) {
	x := sl.header
	h := sl.maxHeight - 1
	for {
		next := x.levels[h]
		cmp := 1
		if next != nil {
			cmp = sl.cmp.compare(next.key, key)
		}
		if cmp < 0 {
			// keep searching in this list
			x = next
		} else {
			if prev != nil {
				prev[h] = x
			} else if cmp == 0 {
				// equal
				return next, true
			}
			if h == 0 {
				// bottom
				return next, cmp == 0
			}
			h--
		}
	}
}

func (sl *skipList) put(key, value []byte) error {
	prev := make([]*element, tMaxHeight)
	if node, exist := sl.findGE(key, prev); exist {
		// current list contained the key, overwrite value
		node.value = value
		return nil
	}

	h := sl.randomHeight()
	if h > sl.maxHeight {
		for i := sl.maxHeight; i < h; i++ {
			prev[i] = sl.header
		}
		sl.maxHeight = h
	}

	ele := newElement(key, value, h)
	for i := 0; i < h; i++ {
		ele.levels[i] = prev[i].levels[i]
		prev[i].levels[i] = ele
	}
	sl.n++
	return nil
}

func (sl *skipList) get(key []byte) (value []byte, err error) {
	if node, exist := sl.findGE(key, nil); exist {
		value = node.value
	} else {
		err = ErrNotFound
	}
	return
}

func newSkipList(cmp comparator) *skipList {
	return &skipList{
		cmp:       cmp,
		rnd:       rand.New(rand.NewSource(0xdeadbeaf)),
		header:    &element{levels: make([]*element, tMaxHeight)},
		maxHeight: 1,
		n:         0,
	}
}
