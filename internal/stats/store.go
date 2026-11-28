package stats

import (
	"sync"
	"time"

	"github.com/Mayankkaira/darkwave/internal/probe"
)

type URLStats struct {
	URL          string
	TotalChecks  int
	Successes    int
	Failures     int
	TotalLatency time.Duration
	LastStatus   int
	LastError    string
	IsHealthy    bool
}

func (u URLStats) AvgLatency() time.Duration {
	if u.TotalChecks == 0 {
		return 0
	}
	return u.TotalLatency / time.Duration(u.TotalChecks)
}

func (u URLStats) UpTimePercent() float64 {
	if u.TotalChecks == 0 {
		return 0.0
	}
	return (float64(u.Successes) / float64(u.TotalChecks)) * 100
}

type StatsStore struct {
	mu    sync.RWMutex
	items map[string]*URLStats
}

func NewStore() *StatsStore {
	items := make(map[string]*URLStats)
	return &StatsStore{items: items}
}

func (s *StatsStore) Record(result probe.Result) {
	s.mu.Lock()
	defer s.mu.Unlock()
	stats, ok := s.items[result.URL]
	if !ok {
		stats = &URLStats{URL: result.URL}
		s.items[result.URL] = stats
	}
	stats.TotalChecks++
	stats.TotalLatency += result.Latency
	stats.LastStatus = result.StatusCode
	if result.Err != nil {
		stats.Failures++
		stats.IsHealthy = false
		stats.LastError = result.Err.Error()
	} else {
		stats.Successes++
		stats.IsHealthy = true
		stats.LastError = ""
	}
}

func (s *StatsStore) Snapshot() map[string]URLStats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	snap := make(map[string]URLStats, len(s.items))
	for key, value := range s.items {
		snap[key] = *value
	}
	return snap
}
