package sse

import "time"

func (conn *SSEConnection) GetStats() ConnectionStats {
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
	ClientID          string        `json:"client_id"`
	NotificationsSent int           `json:"notifications_sent"`
	HeartbeatsSent    int           `json:"heartbeats_sent"`
	Duration          time.Duration `json:"duration"`
	Connected         bool          `json:"connected"`
	LastActivity      time.Time     `json:"last_activity"`
}
