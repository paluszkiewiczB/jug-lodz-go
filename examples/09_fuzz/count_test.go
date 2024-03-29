package chars_test

import (
	chars "github.com/paluszkiewiczB/jug-go-lodz/examples/09_fuzz"
	"testing"
)

// Run the fuzz tests with command:
// go test -fuzz=.
//
// See more about fuzzing: https://go.dev/doc/security/fuzz/
func Fuzz_Add(f *testing.F) {
	f.Add(0, "a")
	f.Add(1, "foo")
	f.Add(3, "rododendron")
	f.Add(2, "idk")

	f.Fuzz(func(t *testing.T, a int, b string) {
		chars.At(a, b)
	})
}
