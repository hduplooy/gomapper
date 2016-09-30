package main

import (
	"fmt"

	mp "github.com/hduplooy/gomapper"
)

func main() {
	names := []string{"John", "Peter", "Susan"}
	ages := []int{12, 15, 13}
	heights := []float64{1.23, 1.5, 1.14}

	persons, _ := mp.Map("", func(vals []interface{}) (interface{}, error) {
		return fmt.Sprintf("%s|%d|%f", vals[0].(string), vals[1].(int), vals[2].(float64)), nil
	}, names, ages, heights)
	fmt.Printf("%v\n", persons)

	mp.ForEach(func(vals []interface{}) error {
		fmt.Printf("Name=%s Age=%d Height=%f\n", vals[0].(string), vals[1].(int), vals[2].(float64))
		return nil
	}, names, ages, heights)

	ans := mp.Fold(func(val1, val2 interface{}) interface{} {
		return val1.(int) + val2.(int)
	}, 1, 2, 3, 4, 5)
	fmt.Println(ans.(int))
}
