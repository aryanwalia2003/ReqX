package runner

import (
	"context"
	"time"
)

// runStages walks through each Stage, adjusting the desired worker count.
// Workers are long-lived VUs; scale changes only toggle idle vs active.
func (s *Scheduler) runStages(ctx context.Context) {
	for _, stage := range s.cfg.Stages {
		s.setDesiredWorkers(stage.TargetWorkers)

		stageTimer := time.NewTimer(stage.Duration)
		select {
		case <-ctx.Done():
			stageTimer.Stop()
			return
		case <-stageTimer.C:
		}

		if ctx.Err() != nil {
			return
		}
	}
}