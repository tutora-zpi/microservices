package cache

import "fmt"

type Key func(suffix string) string

var (
	MeetingKey Key = func(suffix string) string { return fmt.Sprintf("meeting:%s", suffix) }
)
