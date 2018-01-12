package lazyskiplist

import (
	"math/rand"
	"sync"

	"github.com/petermattis/goid"
)

// Fuzz fuzz test
func Fuzz(data []byte) int {
	//TurnOnDebug()
	debugf("input data")
	for _, d := range data {
		debugf("%d ,", int(d))
	}
	debugf("")

	var l = New(func(v1, v2 interface{}) bool {
		i1 := v1.(int)
		i2 := v2.(int)
		return i1 < i2
	})

	var wait sync.WaitGroup
	for _, d := range data {
		if rand.Int()%2 == 0 {
			debugf("[%d] add %d", goid.Get(), int(d))
			wait.Add(1)
			go func(i int) {
				l.Add(i)
				wait.Done()
			}(int(d))
		} else {
			debugf("[%d] delete %d", goid.Get(), int(d))
			wait.Add(1)
			go func(i int) {
				l.Remove(i)
				wait.Done()
			}(int(d))
		}

		if rand.Int()%7 == 0 {
			wait.Add(1)
			go func() {
				defer wait.Done()
				iter := l.Iterator()
				v, ok := iter.Next()
				if !ok {
					return
				}
				debugf("[%d] [", goid.Get())
				defer debugf("[%d] ]", goid.Get())

				debugf("[%d] %d", goid.Get(), v.(int))
				prev := v
				for {
					v, ok := iter.Next()
					if !ok {
						return
					}
					debugf("[%d] %d", goid.Get(), v.(int))
					if prev.(int) >= v.(int) {

						panic("not sorted")
					}
					prev = v
				}
			}()
		}
	}
	wait.Wait()
	return 1
}
