package request

import "path"

type GetAudioRequest struct {
	MeetingID string `json:"meetingId"`
	AudioName string `json:"audioName,omitempty"`
}

func (g *GetAudioRequest) Key() string {
	return path.Join(g.MeetingID, g.AudioName)
}
