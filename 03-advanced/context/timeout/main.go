// Topic: timeout | Concept: context | Stage: 03-advanced
// Run:   go run ./03-advanced/context/timeout
//
// Explore:
// * context.Context carries deadlines and cancellation across calls
// * context.WithTimeout / WithDeadline: cancelled automatically after a time
// * Always `defer cancel()` to release resources, even if the timeout fires
// * Long-running work checks ctx.Done() in a select
// * ctx.Err() tells you why: context.DeadlineExceeded or context.Canceled
// * Pass ctx as the first parameter; never store it in a struct
// * net/http and database/sql accept a ctx for exactly this
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// slowQuery pretends to be a DB call that takes `cost` to finish.
func slowQuery(ctx context.Context, cost time.Duration) (string, error) {
	select {
	case <-time.After(cost):
		return fmt.Sprintf("rows after %v", cost), nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func main() {
	fmt.Println("Topic: context with timeout")

	for _, cost := range []time.Duration{20 * time.Millisecond, 120 * time.Millisecond} {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		res, err := slowQuery(ctx, cost)
		cancel()

		switch {
		case errors.Is(err, context.DeadlineExceeded):
			fmt.Printf("query costing %v: timed out (%v)\n", cost, err)
		case err != nil:
			fmt.Println("other error:", err)
		default:
			fmt.Printf("query costing %v: %s\n", cost, res)
		}
	}

	// Deadline inspection
	deadline := time.Now().Add(40 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	if d, ok := ctx.Deadline(); ok {
		fmt.Printf("deadline set, about %dms left\n", time.Until(d).Round(10*time.Millisecond).Milliseconds())
	}
	<-ctx.Done()
	fmt.Println("after deadline, ctx.Err() =", ctx.Err())
}
