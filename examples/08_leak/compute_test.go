package compute_test

import (
	"context"
	compute "github.com/paluszkiewiczB/jug-go-lodz/examples/08_leak"
	"go.uber.org/goleak"
	"testing"
)

func TestMain(m *testing.M) {
	defer goleak.VerifyTestMain(m)
}

func Test_SumConcurrently(t *testing.T) {
	// Enable to leak a goroutine
	compute.HasLeak = true
	ctx, cf := context.WithCancel(context.Background())
	defer cf()

	res, err := compute.SumConcurrently(ctx, 1_000_000)
	if err != nil {
		t.Errorf("unexpected error: %q", err)
	}

	expect := sumUpTo(1_000_000)
	if expect != res {
		t.Errorf("expected %d, got: %d", expect, res)
	}
}

// https://cseweb.ucsd.edu/groups/tatami/handdemos/sum/
func sumUpTo(n int) int {
	n -= 1 // it is a sum of [0;n) thus n is excluded
	return n * (n + 1) / 2
}
