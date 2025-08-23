package entities

type EventStatus BaseStatus

const (
	EventProcessing EventStatus = "booster.processing"
	EventSuccess EventStatus = "booster.success"
	EventError EventStatus = "booster.error"
	EventFailed EventStatus = "booster.failed"
	EventQueued EventStatus = "booster.queued"
	EventBatchQueued EventStatus = "booster.batch_queued"
	EventCancelled EventStatus = "booster.cancelled"
)
