package compute

import (
	"context"
	"fmt"
)

var (
	HasLeak = false
)

// SumConcurrently sums all the integers starting in range [0; upto).
// It means upto is out of range and will not be added to the sum.
func SumConcurrently(ctx context.Context, upto int) (int, error) {
	select {
	case <-ctx.Done():
		return 0, fmt.Errorf("summing invoked with cancelled context, %s", ctx.Err())
	default:
	}

	res := make(chan int)

	workers := 0
	for i := 0; i < upto; {
		workers++

		from := i
		to := min(upto, from+1000)
		go func() {
			res <- sum(ctx, from, to)
		}()
		i = to
	}

	total := 0
	for i := 0; i < workers; i++ {
		select {
		case <-ctx.Done():
			return total, fmt.Errorf("summing cancelled before all workers finished, result is not complete, %w", ctx.Err())
		case r := <-res:
			total += r
		}
	}

	if HasLeak {
		go func() {
			fmt.Printf("hehe leak: %d\n", <-res)
		}()
	}

	return total, nil
}

func sum(ctx context.Context, from, to int) int {
	total := 0
	for i := from; i < to; i++ {
		// checking it in every operation is an overkill and most likely a performance killer
		select {
		case <-ctx.Done():
			return total
		default:
		}
		total += i
	}

	return total
}
