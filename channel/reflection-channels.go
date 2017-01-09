package main

import (
	"reflect"
	"strings"
)

func makeChannel(t reflect.Type, chanDir reflect.ChanDir, buffer int) reflect.Value {
	ctype := reflect.ChanOf(chanDir, t)
	return reflect.MakeChan(ctype, buffer)
}

type T interface{}

type Query struct {
	input reflect.Value
}

func (q *Query) Apply(f T) *Query {
	value := reflect.ValueOf(f)
	if value.Kind() != reflect.Func {
		panic("Apply() parameter must be a function")
	}

	rtype := value.Type().Out(0)
	output := makeChannel(rtype, reflect.BothDir, 0)
	go func() {
		var elem reflect.Value
		for ok := true; ok; {
			if elem, ok = q.input.Recv(); ok {
				result := value.Call([]reflect.Value{elem})
				output.Send(result[0])
			}
		}
		output.Close()
	}()
	return &Query{output}
}

func From(array T) *Query {
	value := reflect.ValueOf(array)
	if value.Kind() != reflect.Slice {
		panic("From() parameter must be a slice")
	}

	etype := value.Type().Elem()
	output := makeChannel(etype, reflect.BothDir, 0)
	go func() {
		for i := 0; i != value.Len(); i++ {
			output.Send(value.Index(i))
		}
		output.Close()
	}()
	return &Query{output}
}

func (q *Query) Items() <-chan T {
	output := make(chan T)
	go func() {
		for ok := true; ok; {
			var elem reflect.Value
			if elem, ok = q.input.Recv(); ok {
				output <- elem.Interface()
			}
		}
		close(output)
	}()
	return output
}

func (q *Query) StringItems() <-chan string {
	output := make(chan string)
	go func() {
		for elem := range q.Items() {
			output <- elem.(string)
		}
		close(output)
	}()
	return output
}

func main() {
	t := []int{0, 1, 2}
	m := map[int]string{0: "zero", 1: "one", 2: "two"}
	for elem := range From(t).Apply(func(x int) string { return m[x] }).Apply(strings.ToUpper).StringItems() {
		println(elem)
	}
}
