package trace

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"
)

var goroutineSpace = []byte("goroutine ")
var mu sync.Mutex
var m = make(map[uint64]int)

// curGoroutineID returns goroutine id, copy as /src/net/http/h2_bundle.go
func curGoroutineID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	// Parse the 4707 out of "goroutine 4707 ["
	b = bytes.TrimPrefix(b, goroutineSpace)
	i := bytes.IndexByte(b, ' ')
	if i < 0 {
		panic(fmt.Sprintf("No space found in %q", b))
	}
	b = b[:i]
	n, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse goroutine ID out of %q: %v", b, err))
	}
	return n
}

func printTrace(id uint64, name, arrow, file string, indent, line int) {
	indents := ""
	for i := 0; i < indent; i++ {
		indents += " "
	}
	fmt.Printf("g[%05d]:%s%s%s file:%s line%d\n", id, indents, arrow, name, file, line)
}

func Trace() func() {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		panic("not found caller")
	}
	fn := runtime.FuncForPC(pc)
	name := fn.Name()

	gId := curGoroutineID()

	mu.Lock()
	indents := m[gId]    // 获取当前gid对应的缩进层次
	m[gId] = indents + 3 // 缩进层次+3后存入map
	mu.Unlock()
	printTrace(gId, name, "->", file, indents+1, line)
	return func() {
		mu.Lock()
		indents := m[gId]    // 获取当前gid对应的缩进层次
		m[gId] = indents - 3 // 缩进层次-3后存入map
		mu.Unlock()
		printTrace(gId, name, "<-", file, indents, line)
	}
}
