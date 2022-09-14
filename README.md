# hduplooy/gomapper

## Basic map and for-each functionality in golang (also concurrent versions)

**These are changed to a better working set of functions**

### Map

These functions apply a provided functions to successive elements of provided slices. The concurrent versions can be used for basic map-reduce actions.

For example if we want to combine the successive
elements of [John,Peter,Susan], [12,15,13] and [1.23, 1.5, 1.14] to generate a html table we can do it as follows:

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
 
 This will produce:
 
    <tr><td>John</td><td>12</td><td>1.23</td></tr>
    <tr><td>Peter</td><td>15</td><td>1.5</td></tr>
    <tr><td>Susan</td><td>13</td><td>1.14</td></tr>

### ForEach

If we didn't want to return the results but rather just send it out directly with fmt.Printf or fmt.Fprintf we can do it as follows:

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

This will just print the same results out, but only return an error if any.

### Fold

The Fold function will apply a function successively on values in a slice and return the result when applied to all values. If for
example we want to add the integers provided in a slice (and let's assume it was a much more involved operation than simple sums).

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

This will apply the function first to 1 and 2, and then to the result and 3 and then to the result and 4 etc. Effectively it does this:  

((((((((1+2) + 3) + 4) + 5) + 6) + 7) + 8) + 9)

### Filter

To filter an array of values one can use the Filter function with arguments a function that will return true or false on an individual value and an array of values. The returned array is of the same type as the input array.

    package main        
    
    import (
        "fmt"
    
        mp "github.com/hduplooy/gomapper"
    )
    
    func main() {
        ans := mp.Filter(func(val1 interface{}) bool {
            return val1.(int) % 2 == 0
        }, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
        fmt.Printf("%v\n",ans)
    }

Here we filter an array of integers 1 to 10 based on if the element is even (divisible by 2).

### Count

Similar to Filter you can also Count the number of elements for which the provided function returns true.

    package main
    
    import (
        "fmt"
    
        mp "github.com/hduplooy/gomapper"
    )
    
    func main() {
        ans := mp.Count(func(val1 interface{}) bool {
            return val1.(int) % 2 == 0
        }, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
        fmt.Printf("%v\n",ans)
    }

This will count the number of even numbers in the provided integer array.

### MapConc

To do Map and ForEach calling of the provided function concurrently use MapConc and ForEachConc. For example, let's say once again
summation was really hard (which it obviously isn't) and we wanted several computers to help in doing it, we can accomplish it as follows:

Let's make use of json-rpc to be the golang server that does the summations for us, then we can do that as follows:

    // Simple json-rpc server to summate a slice of integers and return it
    package main

    import (
	    "log"
	    "net"
	    "net/rpc"
	    "net/rpc/jsonrpc"
    )

    type Agg struct{}

    // Actual function that does the summation of integers
    func (t *Agg) Sum(args *[]int, reply *int) error {
	    var ans int
	    for _, val := range *args {
		    ans += val
	    }
	    *reply = ans
	    return nil
    }

    func main() {
	    cal := new(Agg)
	    server := rpc.NewServer()
	    server.Register(cal)
	    server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	    listener, e := net.Listen("tcp", ":9999")
	    if e != nil {
		    log.Fatal("listen error:", e)
	    }
	    for {
		    if conn, err := listener.Accept(); err != nil {
			    log.Fatal("accept error: " + err.Error())
		    } else {
			    log.Printf("new connection established\n")
			    go server.ServeCodec(jsonrpc.NewServerCodec(conn))
		    }
	    }
    }

We can compile this and place it on all the servers were want it to run and start them on each machine. Now for the client part we can do the following:

    package main

    import (
	    "fmt"
	    "log"
	    "net"
	    "net/rpc/jsonrpc"

	    mp "github.com/hduplooy/gomapper"
    )

    // Function to do the actual calling of the json-rpc server to summate a slice of integers for us
    func rpcAdd(add string, args []int) int {
	    client, err := net.Dial("tcp", add)
	    if err != nil {
		    log.Fatal("dialing:", err)
		    return -1
	    }
	    var reply int
	    c := jsonrpc.NewClient(client)
	    err = c.Call("Agg.Sum", args, &reply)
	    if err != nil {
		    log.Fatal("arith error:", err)
		    return -1
	    }
	    return reply
    }

    func main() {
	    // IP Addresses of machines that we want to use to process our function
      // Obviously the correct ones have to be used here
	    addrs := []string{"10.0.0.10:9999", "10.0.0.11:9999", "10.0.0.13:9999"}

      // Call MapConc to distribute the work for us
	    sums, err := mp.MapConc(3, func(vals []interface{}, pos int) (interface{}, error) {
        // Convert the interface{} to ints for the call to the rpc function
		    ints := make([]int, len(vals))
		    for i := 0; i < len(vals); i++ {
			    ints[i] = vals[i].(int)
		    }
		    // let rpcAdd do the json-rpc call to the different servers to do the summations for us (use pos to select the server)
		    return rpcAdd(addrs[pos], ints), nil
	    }, []int{1, 2, 3}, []int{5, 6, 7}, []int{9, 10, 11})
	    fmt.Printf("%v\n", sums)
	    if err != nil {
		    log.Println(err)
		    return
	    }

	    // Just add the results together
	    ans := mp.Fold(func(val1, val2 interface{}) interface{} {
		    return val1.(int) + val2.(int)
	    }, mp.ToInterface(sums.([]int))...)
	    fmt.Printf("%d\n", ans)
    }


