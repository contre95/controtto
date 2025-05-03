package jobs

import "context"

type Manager interface {
	StartJob(ctx context.Context, j Job) error
	StopJob(id ID) error
	GetJob(id ID) (Job, bool)
	ListJobs() []Job
}
