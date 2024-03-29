package race_test

import (
	"fmt"
	race "github.com/paluszkiewiczB/jug-go-lodz/examples/06_race"
	"sync"
	"testing"
)

func setup() func() {
	// setup
	race.Counter = 0

	// cleanup
	return func() {

	}
}

func Example_counter() {
	// safe when accessed non-concurrently
	race.Counter = 0
	fmt.Println(race.Counter)

	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			// must be used for concurrent access
			race.Inc()
		}()
	}

	wg.Wait()
	fmt.Println(race.Counter)

	// Output:
	// 0
	// 10
}

// this test should be run with a flag '-race', e.g.
// go test -race ./...
// task go:test:fast -- -race
//
// If you want to know how you can watch: https://www.youtube.com/watch?v=5erqWdlhQLA
// Slides are available here: https://speakerdeck.com/kavya719/go-run-race-under-the-hood
func Test_Resource(t *testing.T) {
	t.Run("race should be detected", func(t *testing.T) {
		// t.Skip() // remove this line

		t.Cleanup(setup())

		wg := sync.WaitGroup{}
		wg.Add(100_000)
		for i := 0; i < 100_000; i++ {
			go func() {
				defer wg.Done()
				race.Counter++ // concurrent non-serialized access to global variable
			}()
		}

		wg.Wait()
		if race.Counter != 100_000 {
			t.Errorf("Counter should be 100_000, is: %d", race.Counter)
		}
	})

	t.Run("race should not be detected", func(t *testing.T) {
		t.Cleanup(setup())

		wg := sync.WaitGroup{}
		wg.Add(100_000)
		for i := 0; i < 100_000; i++ {
			go func() {
				defer wg.Done()
				race.Inc()
			}()
		}

		wg.Wait()
		if race.Counter != 100_000 {
			t.Errorf("Counter should be 100_000")
		}
	})
}
