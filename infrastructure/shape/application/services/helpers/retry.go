package helpers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Space-Core-Dev/Mistaken.Utils/result"
)

func Retry[T any](ctx context.Context, e func(string) (T, error), connStr string, attempts int, delay, offset int) result.Result[T] {
	if attempts <= 0 {
		attempts = 1
	}

	var err error
	var conn T

	currentDelay := time.Duration(delay) * time.Second
	currentAttempt := 0

	for currentAttempt < attempts {
		conn, err = e(connStr)
		if err == nil {
			return result.OK(conn)
		}

		log.Printf("attempt %d failed. Delay is %v. error %s", currentAttempt+1, currentDelay, err)

		if currentAttempt == attempts-1 {
			break
		}

		select {
		case <-time.After(currentDelay):
		case <-ctx.Done():
			return result.Fail[T](fmt.Errorf("context error %w. Call Error %v", ctx.Err(), err))
		}

		currentDelay += time.Duration(offset) * time.Second
		currentAttempt++
	}

	return result.Fail[T](fmt.Errorf("the connection is failed. %w", err))
}
