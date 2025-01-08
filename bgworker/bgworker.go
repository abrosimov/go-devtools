//go:build ignore

package bgworker

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"sync"
	"time"
)

var (
	ErrJobAlreadyRegistered = errors.New("job already registered")
)

type Job interface {
	Name() string
	Run(context.Context) error
	RepeatEvery() time.Duration
	Ctx() context.Context
}

type JobsRunner struct {
	jobs    map[time.Duration][]Job // TODO: not a string, but group by duration
	jobsMtx sync.Mutex
	exitCh  chan struct{}
}

func NewJobsRunner() *JobsRunner {
	return &JobsRunner{
		jobs: make(map[time.Duration][]Job, 10),
	}
}

func (j *JobsRunner) AddJob(job Job) error {
	j.jobsMtx.Lock()
	defer j.jobsMtx.Unlock()

	_, ok := j.jobs[job.Name()]
	if ok {
		return fmt.Errorf("%q %w", job.Name(), ErrJobAlreadyRegistered)
	}

	j.jobs[job.Name()] = job
	return nil
}

func (j *JobsRunner) RemoveJob(name string) {
	j.jobsMtx.Lock()
	defer j.jobsMtx.Unlock()

	delete(j.jobs, name)
}

func (j *JobsRunner) Run(ctx context.Context) error {
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := j.runRegisteredJobs(ctx)
				_ = err
			case <-j.exitCh:
				return
			}
		}
	}()

	return nil
}

func (j *JobsRunner) runRegisteredJobs(ctx context.Context) error {
	j.jobsMtx.Lock()
	defer j.jobsMtx.Unlock()

	eg, ctx := errgroup.WithContext(ctx)
	for _, job := range j.jobs {
		eg.Go(func() error {
			return job.Run(ctx)
		})
	}
	return eg.Wait()
}

func (j *JobsRunner) Stop(ctx context.Context) error {
	j.exitCh <- struct{}{}
	return nil
}
