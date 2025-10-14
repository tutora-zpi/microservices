package interfaces

type HubManager interface {
	AddMeetingMember(roomID string, client Client)
	AddGlobalMember(client Client)
	RemoveGlobalMember(client Client)

	// roomID maybe meetingID depends from ctx
	Emit(roomID string, messageType int, payload []byte)

	EmitGlobal(messageType int, payload []byte)
}
