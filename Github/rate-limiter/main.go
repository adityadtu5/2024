package main

import (
	"context"
	"fmt"
	"rate-limiter/ratelimit"
	"time"
)

func main() {
	bucket := ratelimit.NewTokenBucket(2, 5)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bucket.Start(ctx)
	for i := 0; i < 100; i++ {
		start := time.Now()
		bucket.Wait()
		end := time.Now()
		fmt.Printf("Processing request[%v] at [%v] processed at [%v]\n", i+1, start, end)
	}
}
