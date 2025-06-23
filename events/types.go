package events

const (
	UnknownEvent = iota
	MessageEvent
)

type EventType int
type Event struct {
	Type EventType
	Text string
	Meta interface{}
}

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(e Event) error
}
