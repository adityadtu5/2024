package ratelimit

import (
	"context"
	"fmt"
	"time"
)

type tokens chan struct{}

type TokenBucket struct {
	count  int64
	tokens tokens
	ticker *time.Ticker
}

func NewTokenBucket(count int64, rate int64) *TokenBucket {
	tokens := make(chan struct{}, count)
	for i := 0; i < int(count); i++ {
		tokens <- struct{}{}
	}
	fmt.Printf("current token count %v\n", len(tokens))
	return &TokenBucket{
		count:  count,
		tokens: tokens,
		ticker: time.NewTicker(time.Duration(int64(1/float64(rate)*1000) * int64(time.Millisecond))),
	}
}

func (tb *TokenBucket) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-tb.ticker.C:
				select {
				case tb.tokens <- struct{}{}:
					fmt.Printf("added token, current token count %v\n", len(tb.tokens))
				default:
				}
			case <-ctx.Done():
				fmt.Println("context cancelled")
				return
			}
		}
	}()
}

func (tb *TokenBucket) Wait() {
	<-tb.tokens
	fmt.Printf("consumed token, current token count %v\n", len(tb.tokens))
}
