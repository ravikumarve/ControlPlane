package proxy

import "sync/atomic"

// Stats tracks proxy activity counters for the admin API.
type Stats struct {
	TotalCalls    atomic.Int64
	Allowed       atomic.Int64
	Blocked       atomic.Int64
	HITLPending   atomic.Int64
	RateLimited   atomic.Int64
	InjectionBlock atomic.Int64
}

// Snapshot returns a point-in-time copy of all stats.
func (s *Stats) Snapshot() StatsSnapshot {
	return StatsSnapshot{
		TotalCalls:     s.TotalCalls.Load(),
		Allowed:        s.Allowed.Load(),
		Blocked:        s.Blocked.Load(),
		HITLPending:    s.HITLPending.Load(),
		RateLimited:    s.RateLimited.Load(),
		InjectionBlock: s.InjectionBlock.Load(),
	}
}

// StatsSnapshot is a point-in-time copy of proxy statistics.
type StatsSnapshot struct {
	TotalCalls     int64 `json:"total_calls"`
	Allowed        int64 `json:"allowed"`
	Blocked        int64 `json:"blocked"`
	HITLPending    int64 `json:"hitl_pending"`
	RateLimited    int64 `json:"rate_limited"`
	InjectionBlock int64 `json:"injection_blocked"`
}
