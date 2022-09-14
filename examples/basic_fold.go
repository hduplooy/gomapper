package main

import (
	"fmt"

	mp "github.com/hduplooy/gomapper"
)

func main() {
	ans := mp.Fold(func(val1, val2 interface{}) interface{} {
		return val1.(int) + val2.(int)
	}, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	fmt.Println(ans.(int))
}
