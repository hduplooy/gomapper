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
		return fmt.Sprintf("<tr><td>%s</td><td>%d</td><td>%f</td></tr>", vals[0].(string), vals[1].(int), vals[2].(float64)), nil
	}, names, ages, heights)
	fmt.Printf("%v\n", persons)
}
