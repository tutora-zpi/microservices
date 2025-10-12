package pkg

import (
	"fmt"
	"time"
)

func WithS(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

func FormatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	if h > 0 && m > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}
	if h > 0 {
		return fmt.Sprintf("%dh", h)
	}
	return fmt.Sprintf("%dm", m)
}
