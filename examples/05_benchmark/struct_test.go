package data_test

import (
	data "github.com/paluszkiewiczB/jug-go-lodz/examples/05_benchmark"
	"io"
	"testing"
)

// Run the benchmarks with command:
// go test -bench=. -benchmem
//
// Show compiler optimizations like stack escape analysis and method devirtualization with command:
// go test -gcflags='-m'
//
// You can increase details level setting -m=2 or -m=3
func Benchmark_Empty(b *testing.B) {
	for range b.N {
		w := data.Empty{}
		w.Write(nil)
	}
}

func Benchmark_NotEmpty(b *testing.B) {
	n := "name"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := data.NotEmpty{Name: n}
		w.Write(nil)
	}
}

func Benchmark_EmptyVirtual(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := data.NewEmpty()
		w.Write(nil)
	}
}

func Benchmark_NotEmptyVirtual(b *testing.B) {
	n := "name"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := data.NewNotEmpty(n)
		w.Write(nil)
	}
}

func Benchmark_NotEmptyDevirtualized(b *testing.B) {
	n := "name"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var w io.Writer = &data.NotEmpty{Name: n}
		w.Write(nil)
	}
}

func Benchmark_SmallSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = make([]byte, 65536) // does not allocate, 2^16
	}
}

func Benchmark_BigSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = make([]byte, 65537) // does allocate, 2^16 + 1
	}
}
