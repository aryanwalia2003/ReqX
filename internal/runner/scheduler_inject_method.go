package runner

import (
	"context"
)

// conduct is the controller goroutine: it manages worker lifecycle and stage ramps.
// When done it waits for workers and closes s.results.
func (s *Scheduler) conduct(ctx context.Context) {
	defer func() {
		// Ensure all workers stop, then close the result stream.
		if s.cancel != nil {
			s.cancel()
		}
		s.wg.Wait()
		close(s.results)
	}()

	// Spawn a fixed pool of long-lived VUs once; stages just idle/un-idle them.
	max := s.cfg.MaxWorkers
	if len(s.cfg.Stages) > 0 {
		max = maxTargetWorkers(s.cfg.Stages)
	}
	if max < 1 {
		max = 1
	}
	for i := 1; i <= max; i++ {
		s.spawnWorker(ctx, i)
	}

	if len(s.cfg.Stages) > 0 {
		s.runStages(ctx)
		return
	}

	// Duration mode: keep desired workers fixed until ctx expires.
	s.setDesiredWorkers(s.cfg.MaxWorkers)
	<-ctx.Done()
}

func maxTargetWorkers(stages []Stage) int {
	m := 0
	for _, st := range stages {
		if st.TargetWorkers > m {
			m = st.TargetWorkers
		}
	}
	return m
}

func (s *Scheduler) setDesiredWorkers(n int) {
	if n < 0 {
		n = 0
	}
	s.desiredWorkers.Store(int64(n))

	// Wake any idled workers.
	s.wakeMu.Lock()
	close(s.wakeCh)
	s.wakeCh = make(chan struct{})
	s.wakeMu.Unlock()
}