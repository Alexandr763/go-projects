package ratelimit

import (
	"sync"
	"time"
)

type Limiter struct {
	mu     sync.Mutex
	hits   map[string][]time.Time
	limit  int
	window time.Duration
	now    func() time.Time
}

func New(limit int, window time.Duration) *Limiter {
	return &Limiter{
		hits:   map[string][]time.Time{},
		limit:  limit,
		window: window,
		now:    time.Now,
	}
}

func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	prev := l.hits[key]
	recent := make([]time.Time, 0, len(prev)+1)
	for _, at := range prev {
		if now.Sub(at) < l.window {
			recent = append(recent, at)
		}
	}
	if len(recent) >= l.limit {
		l.hits[key] = recent
		return false
	}
	l.hits[key] = append(recent, now)
	return true
}
