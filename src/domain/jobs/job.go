package jobs

import (
	"time"
)

type ID string

type Status string

const (
	StatusPending  Status = "pending"
	StatusRunning  Status = "running"
	StatusStopped  Status = "stopped"
	StatusFailed   Status = "failed"
	StatusFinished Status = "finished"
)

type Job struct {
	ID        ID
	Type      string
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
	Params    map[string]string
}
