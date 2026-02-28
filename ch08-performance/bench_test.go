package main

import "testing"

func BenchmarkSumSlice(b *testing.B) {
	data := make([]int, 10000)
	for i := range data {
		data[i] = i
	}
	for b.Loop() {
		_ = sumSlice(data)
	}
}
