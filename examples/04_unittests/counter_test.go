package counter_test

import (
	counter "github.com/paluszkiewiczB/jug-go-lodz/examples/04_unittests"
	"github.com/paluszkiewiczB/jug-go-lodz/examples/04_unittests/foo/footest"
	"testing"
	"time"
)

func Test_Counter(t *testing.T) {
	c := counter.NewUnsigned()

	c.Inc()

	if v := c.Get(); 1 != v {
		t.Fatalf("counter should be 1, is: %d", v)
	}

	t.Run("counter should be incremented in a loop", func(t *testing.T) {
		tt := map[string]struct {
			iterations int
			expected   uint
		}{
			"ten": {
				iterations: 10,
				expected:   10,
			},
			"hundred": {
				iterations: 100,
				expected:   100,
			},
		}

		for name, tc := range tt {
			t.Run("increment "+name+" times", func(t *testing.T) {
				c := counter.NewUnsigned()
				for i := 0; i < tc.iterations; i++ {
					c.Inc()
				}

				if v := c.Get(); tc.expected != v {
					t.Fatalf("counter should be %d, is: %d", tc.expected, v)
				}
			})
		}
	})

	t.Run("this test is really slow", func(t *testing.T) {
		if testing.Short() {
			t.Skipf("it is way too slow")
		}

		time.Sleep(5 * time.Second)
		t.Log("just woke up")
	})
}

func Test_FromFoo(t *testing.T) {
	footest.ShouldReturnOne(t, func() int {
		return 1
	})
}
