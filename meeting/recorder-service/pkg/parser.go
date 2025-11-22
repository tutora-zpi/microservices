package pkg

import (
	"strconv"
	"strings"
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

func GetFileName(s string) string {
	tokens := strings.Split(s, "/")
	last := tokens[len(tokens)-1]
	return last
}
