package interfaces

type HubManager interface {
	AddRoomMember(roomID string, client Client) (roomUsers []string)
	AddGlobalMember(client Client)

	RemoveGlobalMember(client Client)
	RemoveRoomMember(roomID string, client Client) (roomUsers []string)

	Emit(roomID string, payload []byte, pred func(id string) bool)
	EmitGlobal(payload []byte)
	EmitToClient(clientID string, payloads [][]byte)
	EmitToClientInRoom(roomID, clientID string, payloads [][]byte)

	Close()
}
