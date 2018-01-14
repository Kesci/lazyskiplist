package lazyskiplist

import (
	"testing"
)

func Test(t *testing.T) {
	l := New(func(v1, v2 interface{}) bool { return v1.(int) < v2.(int) }, func(v1, v2 interface{}) bool { return v1.(int) == v2.(int) })
	l.Add(1)
	l.Add(2)
	l.Add(3)
	if !l.Contains(1) {
		t.Fatalf("not contains %d", 1)
	}
	if !l.Contains(2) {
		t.Fatalf("not contains %d", 2)
	}
	if !l.Contains(3) {
		t.Fatalf("not contains %d", 3)
	}
	l.Remove(2)
	if l.Contains(2) {
		t.Fatalf("contains %d", 2)
	}
	/*
		iter := l.Iterator()
		ints := []int{}
		for {
			if v, ok := iter.Next(); ok {
				ints = append(ints, v.(int))
			} else {
				break
			}
		}
		if ints[0] != 1 || ints[1] != 3 {
			t.Fail()
		}
	*/
}
