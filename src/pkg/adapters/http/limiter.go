package http

import (
	"net/http"
	"sync"
	"time"
)

const (
	maxRequestsPerSecond = 100
	timeWindow           = time.Second
)

type RateLimiter struct {
	mu          sync.Mutex
	requests    chan struct{}
	resetTicker *time.Ticker
}

func NewRateLimiter(maxRequests int, timeDuration time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests:    make(chan struct{}, maxRequestsPerSecond),
		resetTicker: time.NewTicker(timeWindow),
	}

	go func() {
		for range rl.resetTicker.C {
			rl.mu.Lock()
			for len(rl.requests) > 0 {
				<-rl.requests
			}
			//close(rl.requests)
			//rl.requests = make(chan struct{}, maxRequestsPerSecond)
			rl.mu.Unlock()
		}
	}()
	return rl
}

func (rl *RateLimiter) Allow() bool {
	select {
	case rl.requests <- struct{}{}:
		return true
	default:
		return false
	}
}

func (rl *RateLimiter) RateLimiting(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
