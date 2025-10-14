package general

type HeartbeatEvent struct {
	Timestamp int `json:"timestamp"`
}

func (u *HeartbeatEvent) IsValid() error {
	return nil
}

func (u *HeartbeatEvent) Name() string {
	return "heartbeat"
}
