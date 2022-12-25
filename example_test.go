package trace_test

import (
	"github.com/JoyZF/go-trace"
)

func a() {
	defer trace.Trace()
	b()
}

func b() {
	defer trace.Trace()()
	c()
}

func c() {
	defer trace.Trace()()
	d()
}

func d() {
	defer trace.Trace()()
}

func ExampleTrace() {
	a()
	// Output:
	//g[00001]: ->github.com/JoyZF/go-trace_test.b file:E:/workspace/go-trace/example_test.go line13
	//g[00001]:    ->github.com/JoyZF/go-trace_test.c file:E:/workspace/go-trace/example_test.go line18
	//g[00001]:       ->github.com/JoyZF/go-trace_test.d file:E:/workspace/go-trace/example_test.go line23
	//g[00001]:         <-github.com/JoyZF/go-trace_test.d file:E:/workspace/go-trace/example_test.go line23
	//g[00001]:      <-github.com/JoyZF/go-trace_test.c file:E:/workspace/go-trace/example_test.go line18
	//g[00001]:   <-github.com/JoyZF/go-trace_test.b file:E:/workspace/go-trace/example_test.go line13
	//g[00001]: ->github.com/JoyZF/go-trace_test.a file:E:/workspace/go-trace/example_test.go line10
}
