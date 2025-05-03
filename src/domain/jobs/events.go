package jobs

type EventType string
type EventData string

// Define constants for the specific event types
const (
	EventStarted   EventType = "job.started"
	EventCompleted EventType = "job.completed"
	EventFailed    EventType = "job.failed"
	EventStopped   EventType = "job.stopped"
)

type Event interface {
	GetType() EventType
	GetInfo() EventData
}

type JobEvent struct {
	JobID ID
	Type  EventType
	Data  EventData // Can be a message, error, or other info related to the event
}

func (e JobEvent) GetInfo() EventData {
	return e.Data
}

func (e JobEvent) GetType() EventType {
	return e.Type
}
