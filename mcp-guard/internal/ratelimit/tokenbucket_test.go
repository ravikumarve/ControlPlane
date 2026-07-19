package ratelimit

import (
	"testing"
	"time"
)

func TestParseRate(t *testing.T) {
	tests := []struct {
		input    string
		wantTok  int
		wantDur  time.Duration
		wantErr  bool
	}{
		{"100/m", 100, time.Minute, false},
		{"10/s", 10, time.Second, false},
		{"1000/h", 1000, time.Hour, false},
		{"50/min", 50, time.Minute, false},
		{"5/sec", 5, time.Second, false},
		{"200/hour", 200, time.Hour, false},
		{"", 0, 0, true},
		{"abc", 0, 0, true},
		{"100/x", 0, 0, true},
		{"0/s", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tok, dur, err := ParseRate(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for %q", tt.input)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error for %q: %v", tt.input, err)
				return
			}
			if tok != tt.wantTok {
				t.Errorf("got tokens %d, want %d", tok, tt.wantTok)
			}
			if dur != tt.wantDur {
				t.Errorf("got duration %v, want %v", dur, tt.wantDur)
			}
		})
	}
}

func TestBucketAllow(t *testing.T) {
	b := NewBucket(3, time.Minute)

	// First 3 should be allowed
	if !b.Allow() {
		t.Error("expected allow #1")
	}
	if !b.Allow() {
		t.Error("expected allow #2")
	}
	if !b.Allow() {
		t.Error("expected allow #3")
	}
	// 4th should be denied (bucket empty)
	if b.Allow() {
		t.Error("expected deny #4")
	}
}

func TestBucketAllowN(t *testing.T) {
	b := NewBucket(5, time.Minute)

	if !b.AllowN(3) {
		t.Error("expected allow 3")
	}
	if !b.AllowN(2) {
		t.Error("expected allow 2")
	}
	if b.AllowN(1) {
		t.Error("expected deny after consuming all")
	}
}

func TestBucketRefill(t *testing.T) {
	b := NewBucket(2, 50*time.Millisecond)

	if !b.Allow() {
		t.Error("expected allow #1")
	}
	if !b.Allow() {
		t.Error("expected allow #2")
	}
	if b.Allow() {
		t.Error("expected deny, bucket empty")
	}

	// Wait for refill
	time.Sleep(60 * time.Millisecond)

	if !b.Allow() {
		t.Error("expected allow after refill")
	}
}

func TestBucketRemaining(t *testing.T) {
	b := NewBucket(10, time.Minute)
	b.AllowN(3)
	rem := b.Remaining()
	if rem < 6.5 || rem > 10 {
		t.Errorf("expected ~7 remaining, got %f", rem)
	}
}

func TestBucketReset(t *testing.T) {
	b := NewBucket(5, time.Minute)
	b.AllowN(5)
	if b.Allow() {
		t.Error("expected deny, bucket empty")
	}
	b.Reset()
	if !b.Allow() {
		t.Error("expected allow after reset")
	}
}

func TestKeyedLimiter(t *testing.T) {
	kl := NewKeyedLimiter(2, time.Minute)

	if !kl.Allow("alice") {
		t.Error("expected alice #1")
	}
	if !kl.Allow("alice") {
		t.Error("expected alice #2")
	}
	if kl.Allow("alice") {
		t.Error("expected alice deny #3")
	}

	// Bob should have his own pool
	if !kl.Allow("bob") {
		t.Error("expected bob #1")
	}
	if !kl.Allow("bob") {
		t.Error("expected bob #2")
	}
}

func TestParseAndBuild(t *testing.T) {
	kl, err := ParseAndBuild("100/m")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if kl == nil {
		t.Fatal("expected limiter, got nil")
	}
	if kl.capacity != 100 || kl.interval != time.Minute {
		t.Errorf("unexpected config: capacity=%d interval=%v", kl.capacity, kl.interval)
	}

	// Empty string = nil
	kl2, err := ParseAndBuild("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if kl2 != nil {
		t.Error("expected nil for empty string")
	}
}
