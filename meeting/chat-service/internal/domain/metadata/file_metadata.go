package metadata

import (
	"chat-service/internal/domain/event"
	"fmt"
	"log"
	"mime/multipart"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

var SUPPORTED_FILE_TYPES = []string{
	"image/",
	"text/",
	"application/pdf",
	"application/msword",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"application/vnd.ms-powerpoint",
	"application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"application/vnd.oasis.opendocument.text",
	"application/vnd.oasis.opendocument.spreadsheet",
	"application/vnd.oasis.opendocument.presentation",
}

type FileMetadata struct {
	Extension   string
	File        multipart.File
	FileName    string
	ContentType string

	SentAt   int64
	Content  string
	ChatID   string
	SenderID string
}

func (f *FileMetadata) GenerateUniqueFileName() string {
	timestamp := time.Now().UTC().UnixNano()
	name := fmt.Sprintf("%d%s", timestamp, f.Extension)
	return name
}

func (f *FileMetadata) IsValidContentType() bool {
	if f.ContentType == "" {
		log.Printf("No content file")
		return false
	}

	log.Printf("Content type: %s", f.ContentType)

	for _, st := range SUPPORTED_FILE_TYPES {
		if strings.HasPrefix(f.ContentType, st) {
			return true
		}
	}

	return false
}

func (f *FileMetadata) NewFileMessage(urlToLink, filename string) *event.SendMessageEvent {
	return &event.SendMessageEvent{
		MessageID: bson.NewObjectID().Hex(),
		Content:   f.Content,
		SenderID:  f.SenderID,
		SentAt:    f.SentAt,
		FileLink:  urlToLink,
		FileName:  filename,
		ChatID:    f.ChatID,
	}
}
