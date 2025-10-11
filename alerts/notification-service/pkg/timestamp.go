package pkg

import "time"

func GenerateTimestamp() int64 {
	return time.Now().UTC().Unix()
}
