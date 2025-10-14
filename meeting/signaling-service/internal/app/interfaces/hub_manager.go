package interfaces

type HubManager interface {
	AddMeetingMember(meetingID string, client Client)
	AddGlobalMember(client Client)
	RemoveGlobalMember(client Client)
	RemoveMeetingMemeber(meetingID string, client Client)

	// roomID maybe meetingID depends from ctx
	Emit(roomID string, messageType int, payload []byte, pred func(id string) bool)

	EmitGlobal(messageType int, payload []byte)
}
