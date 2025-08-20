package entities

type EventStatus BaseStatus

const (
	EventQueued     EventStatus = "queued"
	EventProcessing EventStatus = "processing"
	EventSuccess    EventStatus = "success"
	EventError      EventStatus = "error"
	EventFailed     EventStatus = "failed"
	EventCancelled  EventStatus = "cancelled"
)