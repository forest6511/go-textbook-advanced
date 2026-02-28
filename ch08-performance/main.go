package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
)

// sumSlice is a simple function to benchmark
func sumSlice(data []int) int {
	total := 0
	for _, v := range data {
		total += v
	}
	return total
}

// allocateStrings demonstrates memory allocation patterns
func allocateStrings(n int) []string {
	// Pre-allocate with known capacity to avoid reallocation
	result := make([]string, 0, n)
	for i := range n {
		result = append(result, fmt.Sprintf("item-%d", i))
	}
	return result
}

func main() {
	// Green Tea GC is default in Go 1.26
	// Disable with: GOEXPERIMENT=nogreenteagc go build
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	fmt.Printf("GC runs: %d\n", ms.NumGC)

	// pprof CPU profile example
	f, err := os.Create(os.DevNull) // Use /dev/null for demo
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}

	// Do some work
	data := make([]int, 10000)
	for i := range data {
		data[i] = i
	}
	result := sumSlice(data)
	fmt.Println("sum:", result)

	pprof.StopCPUProfile()

	// Show allocated strings
	strs := allocateStrings(5)
	fmt.Println("strings:", len(strs))
}
