package server

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Throttle interface {
	IsThrottled(ip string) bool
}

// rateLimiter .
type rateLimiter struct {
	requests *rate.Limiter
	ip       map[string]*rate.Limiter
	mu       *sync.RWMutex
}

// NewThrottle .
func NewThrottle() Throttle {
	return &rateLimiter{
		requests: rate.NewLimiter(rate.Every(1*time.Second), 6000),
		ip:       make(map[string]*rate.Limiter),
		mu:       &sync.RWMutex{},
	}
}

func (i *rateLimiter) IsThrottled(ip string) bool {
	defer i.mu.Unlock()
	i.mu.Lock()

	ipLimiter, exists := i.ip[ip]
	if !exists {
		i.ip[ip] = rate.NewLimiter(rate.Every(1*time.Second), 100)
		return false
	}
	if !ipLimiter.Allow() {
		return true
	}

	return !i.requests.Allow()
}
