// github.com/hduplooy/gomapper
// Author: Hannes du Plooy
// Mapper function calls in go
// Map - will apply a function to a slice of slices and return a slice based on the application of the function
// ForEach - similar to Map but does not return a slice
// MapConc - similar to Map but does it with Go Routines. This can be used for MapReduce type of actions
// ForEachConc - similar to ForEach but does it with Go Routines
// Fold - applies a function to the first two values in a slice and then iteratively take the answer and the next value in the slice and apply the function on it
package gomapper

import (
	"errors"
	"reflect"
)

type MapFunc func(elms []interface{}) (interface{}, error)
type MapConcFunc func(elms []interface{}, pos int) (interface{}, error)
type ForEachFunc func(elms []interface{}) error
type ForEachConcFunc func(elms []interface{}, pos int) error
type FilterFunc func(elm interface{}) bool

type IFArr []interface{}

// Map apply a function f successively to the slice of slices vals and return a slice of type dstif
// dstif - the type of slice that need to be returned
// f - the function to apply
// vals - the slice of slices on which f is applied
//
// Map("",func(vals[] interface{}) (interface{},error) {
//     return vals[0].(int)*vals[1].(int),nil
// }, []int{1,2,3,4},[]int{5,4,3,2})
// Will return []int{5,10,9,8}
// which is []int{1*5,2*4,3*3,4*2}
func Map(dstif interface{}, f MapFunc, vals ...interface{}) (interface{}, error) {
	sz := -1
	for _, val := range vals {
		if reflect.TypeOf(val).Kind() != reflect.Slice {
			return nil, errors.New("not all provided parameters for map are slices")
		}
		if sz == -1 {
			sz = reflect.ValueOf(val).Len()
		} else if sz != reflect.ValueOf(val).Len() {
			return nil, errors.New("all the slices provided to map are not of the same length")
		}
	}
	slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(dstif)), 0, sz)
	var err error
	for i := 0; i < sz; i++ {
		parms := make([]interface{}, len(vals))
		for j, val := range vals {
			parms[j] = reflect.ValueOf(val).Index(i).Interface()
			if err != nil {
				return nil, err
			}
		}
		ans, err2 := f(parms)
		err = err2
		slice = reflect.Append(slice, reflect.ValueOf(ans))
	}
	return slice.Interface(), err
}

// doMapFunc will call the provided func and notify the caller on donech when it was successful else the error will be
// returned on the errch channel
// f - the function to apply
// slice - is the slice where the results are saved to
// pos - is the position where the result must be saved to
// donech - a true is send on this channel if no error was experienced
// errch - if an error was experienced it will be send back on this channel
// vals - are the parameters to the provided function
func doMapFunc(f MapConcFunc, slice reflect.Value, pos int, donech chan bool, errch chan error, vals []interface{}) {
	ans, err := f(vals, pos)
	slice.Index(pos).Set(reflect.ValueOf(ans))
	if err != nil {
		errch <- err
	} else {
		donech <- true
	}
}

// MapConc - similar to Map but all elements are done concurrently
// The provided function f is applied to the successive elements of the slices of the the slice vals, the results are placed in a
// slice and this is returned
// dstif - the type of slice we want to return
// f - the function that need to be applied
// vals - the slice of slices on which the function is applied
func MapConc(dstif interface{}, f MapConcFunc, vals ...interface{}) (interface{}, error) {
	// Everything is similar to Map up to where the function has to be called
	sz := -1
	for _, val := range vals {
		if reflect.TypeOf(val).Kind() != reflect.Slice {
			return nil, errors.New("not all provided parameters for map are slices")
		}
		if sz == -1 {
			sz = reflect.ValueOf(val).Len()
		} else if sz != reflect.ValueOf(val).Len() {
			return nil, errors.New("all the slices provided to map are not of the same length")
		}
	}
	slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(dstif)), sz, sz)
	cnt := 0
	// Create the two channels we will use to communicate if the function call was successful
	donech := make(chan bool)
	errch := make(chan error)
	var err error
	for i := 0; i < sz; i++ {
		parms := make([]interface{}, len(vals))
		for j, val := range vals {
			parms[j] = reflect.ValueOf(val).Index(i).Interface()
		}
		// Start a new goroutine that will call the function for us
		go doMapFunc(f, slice, i, donech, errch, parms)
	}
	// Get the results from the goroutines
	for cnt < sz {
		select {
		case <-donech:
		case err2 := <-errch:
			err = err2 // If an error was experienced save it
		}
		cnt++
	}
	return slice.Interface(), err
}

func Filter(f FilterFunc, vals interface{}) interface{} {
	dstif := reflect.ValueOf(vals).Index(0).Interface()

	sz := reflect.ValueOf(vals).Len()

	slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(dstif)), 0, 0)

	for i := 0; i < sz; i++ {
		val := reflect.ValueOf(vals).Index(i).Interface()
		if f(val) {
			slice = reflect.Append(slice, reflect.ValueOf(val))
		}
	}
	return slice.Interface()
}

func (arr IFArr) Filter(f FilterFunc) interface{} {
	return Filter(f, arr)
}

func Count(f FilterFunc, vals interface{}) int {
	sz := reflect.ValueOf(vals).Len()

	cnt := 0

	for i := 0; i < sz; i++ {
		val := reflect.ValueOf(vals).Index(i).Interface()
		if f(val) {
			cnt++
		}
	}
	return cnt
}

// ForEach will apply f successively on the slice of slices vals and only returns an error if any
// f - function to apply
// vals - the slice of slices to apply
//
// ForEach(func(v []interface{}) error {
//     fmt.Fprintf(w,"<tr><td>%s</td><td>%d</td></tr>\n",v[0].(string),v[1].(int))
// }, []string{"John","Peter"},[]int{12,44})
// This will send to stream w the following:
// <tr><td>John</td><td>12</td></tr>
// <tr><td>Peter</td><td>44</td></tr>
func ForEach(f ForEachFunc, vals ...interface{}) error {
	sz := -1
	for _, val := range vals {
		if reflect.TypeOf(val).Kind() != reflect.Slice {
			return errors.New("not all provided parameters for map are slices")
		}
		if sz == -1 {
			sz = reflect.ValueOf(val).Len()
		} else if sz != reflect.ValueOf(val).Len() {
			return errors.New("all the slices provided to map are not of the same length")
		}
	}
	for i := 0; i < sz; i++ {
		parms := make([]interface{}, len(vals))
		for j, val := range vals {
			parms[j] = reflect.ValueOf(val).Index(i).Interface()
		}
		err := f(parms)
		if err != nil {
			return err
		}
	}
	return nil
}

// doForEachFunc will call the provided func and notify the caller on donech when it was successful else the error will be
// returned on the errch channel
// f - the function to apply
// pos - the position in the provided slice this is performed on (only added when this needed to be used for decision making in the provided function)
// donech - a true is send on this channel if no error was experienced
// errch - if an error was experienced it will be send back on this channel
// vals - are the parameters to the provided function
func doForEachFunc(f ForEachConcFunc, pos int, donech chan bool, errch chan error, vals []interface{}) {
	err := f(vals, pos)
	if err != nil {
		errch <- err
	} else {
		donech <- true
	}
}

// ForEachConc - similar to ForEach but all elements are done concurrently
// The provided function f is applied to the successive elements of the slices of the the slice vals
// f - the function that need to be applied
// vals - the slice of slices on which the function is applied
func ForEachConc(f ForEachConcFunc, vals ...interface{}) error {
	sz := -1
	for _, val := range vals {
		if reflect.TypeOf(val).Kind() != reflect.Slice {
			return errors.New("not all provided parameters for map are slices")
		}
		if sz == -1 {
			sz = reflect.ValueOf(val).Len()
		} else if sz != reflect.ValueOf(val).Len() {
			return errors.New("all the slices provided to map are not of the same length")
		}
	}
	cnt := 0
	donech := make(chan bool)
	errch := make(chan error)
	for i := 0; i < sz; i++ {
		parms := make([]interface{}, len(vals))
		for j, val := range vals {
			parms[j] = reflect.ValueOf(val).Index(i).Interface()
		}
		// Do goroutine to call the provided function
		go doForEachFunc(f, i, donech, errch, parms)
	}
	var err error
	for cnt < sz {
		select {
		case <-donech:
		case err2 := <-errch:
			err = err2
		}
		cnt++
	}

	return err
}

func Fold(f func(val1, val2 interface{}) interface{}, vals ...interface{}) interface{} {
	if len(vals) < 2 {
		return nil
	}
	ans := vals[0]
	for _, val := range vals[1:] {
		ans = f(ans, val)
	}
	return ans
}

func ToInterface(vals interface{}) []interface{} {
	valsr := reflect.ValueOf(vals)
	if valsr.Kind() != reflect.Slice {
		return nil
	}
	ans := make([]interface{}, valsr.Len())
	for i := 0; i < valsr.Len(); i++ {
		ans = append(ans, valsr.Index(i).Interface())
	}
	return ans
}
