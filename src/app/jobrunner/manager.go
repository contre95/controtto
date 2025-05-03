package jobrunner

import (
	"context"
	"controtto/src/domain/jobs"
	"fmt"
	"sync"
	"time"
)

type JobRunnerFunc func(ctx context.Context, params map[string]string) error

type Manager struct {
	mu       sync.Mutex
	running  map[jobs.ID]runningJob
	registry map[string]JobRunnerFunc
	events   chan jobs.Event
}

type runningJob struct {
	cancel context.CancelFunc
	info   jobs.Job
}

func NewManager() *Manager {
	return &Manager{
		running:  make(map[jobs.ID]runningJob),
		registry: make(map[string]JobRunnerFunc),
		events:   make(chan jobs.Event, 100),
	}
}

// Register a jobs type with its runner function
func (m *Manager) Register(jobType string, fn JobRunnerFunc) {
	m.registry[jobType] = fn
}

func (m *Manager) StartJob(ctx context.Context, j jobs.Job) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.registry[j.Type]; !ok {
		return fmt.Errorf("jobs type '%s' not registered", j.Type)
	}

	if _, exists := m.running[j.ID]; exists {
		return fmt.Errorf("jobs %s is already running", j.ID)
	}

	runCtx, cancel := context.WithCancel(ctx)
	runner := m.registry[j.Type]

	// Save jobs state
	j.Status = jobs.StatusRunning
	j.UpdatedAt = time.Now()
	m.running[j.ID] = runningJob{cancel: cancel, info: j}

	go func(j jobs.Job) {
		m.events <- jobs.Event{JobID: j.ID, Type: jobs.EventStarted}
		err := runner(runCtx, j.Params)

		m.mu.Lock()
		defer m.mu.Unlock()
		delete(m.running, j.ID)

		if err != nil {
			m.events <- jobs.Event{JobID: j.ID, Type: jobs.EventFailed, Data: err.Error()}
		} else {
			m.events <- jobs.Event{JobID: j.ID, Type: jobs.EventCompleted}
		}
	}(j)

	return nil
}

func (m *Manager) StopJob(id jobs.ID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	rj, ok := m.running[id]
	if !ok {
		return fmt.Errorf("jobs %s not found", id)
	}
	rj.cancel()
	delete(m.running, id)
	return nil
}

func (m *Manager) GetJob(id jobs.ID) (jobs.Job, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	rj, ok := m.running[id]
	return rj.info, ok
}

func (m *Manager) ListJobs() []jobs.Job {
	m.mu.Lock()
	defer m.mu.Unlock()

	var jobs []jobs.Job
	for _, rj := range m.running {
		jobs = append(jobs, rj.info)
	}
	return jobs
}

func (m *Manager) Events() <-chan jobs.Event {
	return m.events
}
