package sse

import "time"

func (conn *NotificationStreamConnection) GetStats() ConnectionStats {
	return ConnectionStats{
		ClientID:          conn.ClientID,
		NotificationsSent: conn.NotificationsSent,
		HeartbeatsSent:    conn.HeartbeatsSent,
		Duration:          time.Since(conn.ConnectionTime),
		Connected:         conn.HeartbeatTicker != nil && conn.isHealthy,
		LastActivity:      conn.lastActivity,
	}
}

type ConnectionStats struct {
	ClientID          string        `json:"clientId"`
	NotificationsSent int           `json:"notificationsSent"`
	HeartbeatsSent    int           `json:"heartbeatsSent"`
	Duration          time.Duration `json:"duration"`
	Connected         bool          `json:"connected"`
	LastActivity      time.Time     `json:"lastActivity"`
}
