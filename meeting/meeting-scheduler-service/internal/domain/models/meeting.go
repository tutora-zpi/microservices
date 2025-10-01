package models

import "time"

type Meeting struct {
	ClassID   string    `json:"classId"`
	Timestamp time.Time `json:"timestamp"`
}
