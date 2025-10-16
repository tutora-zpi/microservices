package enum

import "fmt"

type Key func(suffix string) string

var (
	SnapshotKey Key = func(suffix string) string { return fmt.Sprintf("snapshot:%s", suffix) }
	EventKey    Key = func(suffix string) string { return fmt.Sprintf("room:%s:event", suffix) }
)
