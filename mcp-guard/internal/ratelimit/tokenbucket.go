package ratelimit

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Bucket implements a token bucket rate limiter.
// Each bucket allows up to 'capacity' tokens per 'interval'.
type Bucket struct {
	mu        sync.Mutex
	capacity  int
	tokens    float64
	interval  time.Duration
	lastRefill time.Time
}

// ParseRate parses a rate string like "100/m" or "10/s" or "1000/h" into tokens + interval.
func ParseRate(s string) (tokens int, interval time.Duration, err error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, 0, fmt.Errorf("empty rate string")
	}

	// Find the split between number and unit
	split := -1
	for i, r := range s {
		if r < '0' || r > '9' {
			split = i
			break
		}
	}
	if split <= 0 || split >= len(s) {
		return 0, 0, fmt.Errorf("invalid rate format: %q (expected e.g. 100/m)", s)
	}

	numStr := s[:split]
	unit := s[split:]

	tokens, err = strconv.Atoi(numStr)
	if err != nil || tokens <= 0 {
		return 0, 0, fmt.Errorf("invalid rate count: %q", numStr)
	}

	switch unit {
	case "/s", "/sec", "/second":
		interval = time.Second
	case "/m", "/min", "/minute":
		interval = time.Minute
	case "/h", "/hr", "/hour":
		interval = time.Hour
	default:
		return 0, 0, fmt.Errorf("unsupported rate unit: %q (use /s, /m, or /h)", unit)
	}

	return tokens, interval, nil
}

// NewBucket creates a token bucket with the given capacity and refill interval.
func NewBucket(capacity int, interval time.Duration) *Bucket {
	now := time.Now()
	return &Bucket{
		capacity:   capacity,
		tokens:     float64(capacity),
		interval:   interval,
		lastRefill: now,
	}
}

// Allow checks if a token can be consumed. Returns true if allowed.
func (b *Bucket) Allow() bool {
	return b.AllowN(1)
}

// AllowN checks if n tokens can be consumed. Returns true if allowed.
func (b *Bucket) AllowN(n int) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(b.lastRefill)

	// Refill based on elapsed time
	refillPerInterval := float64(b.capacity)
	refillFraction := elapsed.Seconds() / b.interval.Seconds()
	b.tokens += refillFraction * refillPerInterval

	// Clamp
	if b.tokens > float64(b.capacity) {
		b.tokens = float64(b.capacity)
	}
	b.lastRefill = now

	// Try to consume
	if b.tokens >= float64(n) {
		b.tokens -= float64(n)
		return true
	}
	return false
}

// Remaining returns the approximate number of tokens available.
func (b *Bucket) Remaining() float64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	now := time.Now()
	elapsed := now.Sub(b.lastRefill)
	refillFraction := elapsed.Seconds() / b.interval.Seconds()
	tokens := b.tokens + refillFraction*float64(b.capacity)
	if tokens > float64(b.capacity) {
		tokens = float64(b.capacity)
	}
	return tokens
}

// Reset resets the bucket to full capacity.
func (b *Bucket) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.tokens = float64(b.capacity)
	b.lastRefill = time.Now()
}

// KeyedLimiter manages a set of token buckets keyed by identity or tool.
type KeyedLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*Bucket
	capacity int
	interval time.Duration
}

// NewKeyedLimiter creates a new keyed rate limiter.
func NewKeyedLimiter(capacity int, interval time.Duration) *KeyedLimiter {
	return &KeyedLimiter{
		buckets:  make(map[string]*Bucket),
		capacity: capacity,
		interval: interval,
	}
}

// Allow checks if the given key is allowed to proceed.
func (kl *KeyedLimiter) Allow(key string) bool {
	kl.mu.Lock()
	b, ok := kl.buckets[key]
	if !ok {
		b = NewBucket(kl.capacity, kl.interval)
		kl.buckets[key] = b
	}
	kl.mu.Unlock()
	return b.Allow()
}

// ParseAndBuild parses a rate string and creates a KeyedLimiter for per-identity limiting.
func ParseAndBuild(rateStr string) (*KeyedLimiter, error) {
	if rateStr == "" {
		return nil, nil
	}
	tokens, interval, err := ParseRate(rateStr)
	if err != nil {
		return nil, err
	}
	return NewKeyedLimiter(tokens, interval), nil
}
