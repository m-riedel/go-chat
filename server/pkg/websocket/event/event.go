package event

const (
	JoinRoomEvent = iota
	LeaveRoomEvent
	MessageEvent
	SetUsernameEvent
)

type Event struct {
	Type int          `json:"type"`
	Data EventMessage `json:"data"`
}

type EventMessage struct {
	Client    string `json:"client"`
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}
