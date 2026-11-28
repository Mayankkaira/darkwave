package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/Mayankkaira/darkwave/internal/probe"
	"github.com/Mayankkaira/darkwave/internal/stats"
)

type Scheduler struct {
	client   *probe.Client
	store    *stats.StatsStore
	URLs     []string
	workers  int
	interval time.Duration
}

func NewScheduler(client *probe.Client, store *stats.StatsStore, urls []string, w int, interval time.Duration) *Scheduler {
	if w <= 0 {
		w = 1
	}
	return &Scheduler{client: client, store: store, URLs: urls, workers: w, interval: interval}

}

func (s *Scheduler) worker(ctx context.Context, jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for url := range jobs {
		result := s.client.Check(ctx, url)
		s.store.Record(result)
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	jobs := make(chan string, len(s.URLs))
	var wg sync.WaitGroup
  		for w := 0; w < s.workers; w++ {
		wg.Add(1)
		go s.worker(ctx, jobs, &wg)
    }
	}
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
}
