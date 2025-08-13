package events

type EventType int

const (
    EventUnknown EventType = iota
    EventStart
    EventStop
    EventPause
    EventResume
)

func (e EventType) String() string {
    switch e {
    case EventStart:
        return "Start"
    case EventStop:
        return "Stop"
    case EventPause:
        return "Pause"
    case EventResume:
        return "Resume"
    default:
        return "Unknown"
    }
}
