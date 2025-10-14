package enum

type UserStatus string

const (
	Online    UserStatus = "online"
	Busy      UserStatus = "busy"
	OnMeeting UserStatus = "on-meeting"
)
