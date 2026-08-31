package runner

import (
	"fmt"
	"github.com/go-co-op/gocron/v2"
)

type Job struct {
	id           string
	cronSchedule string
	handler      func()
}

type Runner struct {
	jobs      []Job
	scheduler gocron.Scheduler
	isRunning bool
}

func New(jobs []Job) (*Runner, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("unable to start cron: %w", err)
	}
	return &Runner{
		jobs:      jobs,
		scheduler: s,
		isRunning: false,
	}, nil
}

func NewJob(id string, cronSchedule string, handler func()) Job {
	return Job{
		id:           id,
		cronSchedule: cronSchedule,
		handler:      handler,
	}
}

func (r *Runner) Run() error {
	if r.isRunning {
		return fmt.Errorf("failed to Run: scheduler already running")
	}
	for _, j := range r.jobs {
		_, err := r.scheduler.NewJob(
			gocron.CronJob(j.cronSchedule, false),
			gocron.NewTask(j.handler),
			gocron.WithName(j.id),
		)
		if err != nil {
			return fmt.Errorf("unable to start a cron job %q: %w", j.id, err)
		}
	}
	r.scheduler.Start()
	r.isRunning = true
	return nil
}

func (r *Runner) Stop() error {
	if r.scheduler == nil {
		return nil
	}
	err := r.scheduler.Shutdown()
	r.isRunning = false
	return err
}
