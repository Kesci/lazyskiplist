package lazyskiplist_test

import (
	"fmt"

	"github.com/Kesci/lazyskiplist"
)

func Example() {
	l := lazyskiplist.New(func(v1, v2 interface{}) bool { return v1.(int) < v2.(int) })
	l.Add(1)
	fmt.Println(l.Contains(1))
	l.Remove(1)
	fmt.Println(l.Contains(1))

	// Output: true
	// false
}
