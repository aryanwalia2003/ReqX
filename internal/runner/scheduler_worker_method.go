package runner

import (
	"context"
	"time"

	"reqx/internal/http_executor"
	"reqx/internal/scripting"
)

// spawnWorker launches one self-driven virtual user (VU).
// The worker sets up isolated state once (env clone, persona, cookie jar),
// then loops executing the plan until ctx is cancelled.
func (s *Scheduler) spawnWorker(ctx context.Context, id int) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		// Per-worker isolated state (persist across iterations).
		rtCtx := NewRuntimeContext()
		if s.cfg.BaseEnv != nil {
			rtCtx.SetEnvironment(s.cfg.BaseEnv.Clone())
		}
		if len(s.cfg.Personas) > 0 {
			p := s.cfg.Personas[(id-1)%len(s.cfg.Personas)]
			applyPersona(rtCtx, p)
		}

		exec := http_executor.NewDefaultExecutor()
		if s.cfg.NoCookies {
			exec.DisableCookies()
		}

		engine := NewCollectionRunner(exec, nil, nil, scripting.NewGojaRunner())
		engine.SetVerbosity(s.cfg.Verbosity)
		if s.cfg.ClearCookies {
			engine.SetClearCookiesPerRequest(true)
		}

		for {
			// Stage-controlled idling (no job queue contention).
			for {
				if ctx.Err() != nil {
					return
				}
				if int64(id) <= s.desiredWorkers.Load() {
					break
				}
				s.activeWorkers.Store(s.desiredWorkers.Load())
				select {
				case <-ctx.Done():
					return
				case <-s.wake():
				case <-time.After(100 * time.Millisecond):
				}
			}

			s.activeWorkers.Store(s.desiredWorkers.Load())

			iter := int(s.completedIterations.Add(1)) // unique-ish monotonically increasing index
			metrics, err := engine.Run(s.cfg.Plan, rtCtx)
			if err != nil {
				s.failedIterations.Add(1)
			}
			for i := range metrics {
				metrics[i].WorkerID = id
			}
			s.results <- WorkerResult{IterationIndex: iter, Metrics: metrics, Err: err}

			// Optional global-ish RPS control without a central job queue.
			if s.cfg.RPS > 0 {
				desired := s.desiredWorkers.Load()
				if desired < 1 {
					desired = 1
				}
				// Per-worker interval that approximates cfg.RPS overall:
				// perWorkerRPS ≈ cfg.RPS / desired => interval ≈ desired/cfg.RPS seconds.
				interval := time.Duration(float64(time.Second) * float64(desired) / s.cfg.RPS)
				if interval > 0 {
					select {
					case <-ctx.Done():
						return
					case <-time.After(interval):
					}
				}
			}
		}
	}()
}

func (s *Scheduler) wake() <-chan struct{} {
	s.wakeMu.Lock()
	ch := s.wakeCh
	s.wakeMu.Unlock()
	return ch
}
