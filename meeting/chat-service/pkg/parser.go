package pkg

import (
	"net/textproto"
	"strings"
)

func GetFileNameFromHeader(header textproto.MIMEHeader) (filename string) {
	info := header.Get("Content-Disposition")
	if info == "" {
		return ""
	}

	tokens := strings.Split(info, "filename=")
	last := tokens[len(tokens)-1]

	filename = strings.Trim(last, "\"")
	return filename
}
