package main

import (
	"fmt"

	mp "github.com/hduplooy/gomapper"
)

func main() {
	names := []string{"John", "Peter", "Susan"}
	ages := []int{12, 15, 13}
	heights := []float64{1.23, 1.5, 1.14}

	mp.ForEach(func(vals []interface{}) error {
		fmt.Printf("<tr><td>%s</td><td>%d</td><td>%f</td></tr>\n", vals[0].(string), vals[1].(int), vals[2].(float64))
		return nil
	}, names, ages, heights)
}
