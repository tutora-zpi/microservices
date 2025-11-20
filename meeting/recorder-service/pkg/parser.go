package pkg

import (
	"strconv"
	"time"
)

func StrToTime(defValue time.Duration, s string) time.Duration {
	presignT := defValue

	if s != "" {
		presignTimeAsInt, err := strconv.Atoi(s)
		if err == nil {
			presignT = time.Minute * time.Duration(presignTimeAsInt)
		}
	}

	return presignT
}
