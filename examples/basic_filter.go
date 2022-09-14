package main

import (
	"fmt"

	mp "github.com/hduplooy/gomapper"
)

func main() {
	ans := mp.Filter(func(val1 interface{}) bool {
		return val1.(int)%2 == 0
	}, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	fmt.Printf("%v\n", ans)
}
