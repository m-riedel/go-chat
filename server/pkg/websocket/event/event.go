package event

const (
	JoinRoomEvent = iota
	LeaveRoomEvent
	MessageEvent
)

type Event struct {
	Type int         `json:"type"`
	Data interface{} `json:"data"`
}
