package cache

import "fmt"

type CacheKey func(suffix string) string

var (
	BotKey CacheKey = func(suffix string) string { return fmt.Sprintf("room:%s:bot", suffix) }
)
