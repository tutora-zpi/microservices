package handlers

import (
	"chat-service/internal/domain/dto/requests"
	"chat-service/internal/domain/metadata"
	"chat-service/internal/infrastructure/security"
	"chat-service/internal/infrastructure/server"
	"chat-service/pkg"
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

const (
	Auth         string = "Authorization"
	BearerPrefix string = "Bearer "
	ID           string = "userID"
	Token        string = "token"

	File         string = "file"
	FileMetadata string = "fileMetadata"

	ContentType string = "Content-Type"
)

func findToken(r *http.Request) (string, error) {
	var token string

	cookie, err := r.Cookie(Token)
	if err == nil {
		token = cookie.Value
		return token, nil
	}
	log.Println("Not found token in cookie going to find in query")

	token = r.URL.Query().Get(Token)
	if token != "" {
		return token, nil
	}

	log.Println("Not found token in query going to find in header")

	auth := r.Header.Get(Auth)
	if !strings.HasPrefix(auth, BearerPrefix) {
		return "", fmt.Errorf("no bearer prefix in header")
	}

	token = strings.TrimPrefix(auth, BearerPrefix)
	token = strings.TrimSpace(token)

	if token == "" {
		return "", fmt.Errorf("token is empty")
	}

	return token, nil
}

func (h *handlers) IsAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := findToken(r)

		if err != nil {
			log.Printf("Not found token: %s\n", err)
			server.NewResponse(w, pkg.Ptr("Unauthorized access"), http.StatusUnauthorized, nil)
			return
		}

		userID, err := security.DecodeJWT(tokenStr)

		if err != nil {
			log.Printf("Unauthorized access: %s\n", err)
			server.NewResponse(w, pkg.Ptr("Unauthorized access"), http.StatusUnauthorized, nil)
			return
		}

		ctx := r.Context()

		ctx = context.WithValue(ctx, ID, userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func ValidateFileFormData(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(http.DefaultMaxHeaderBytes)
		if err != nil {
			server.NewResponse(w, pkg.Ptr("File is too big, max size is 1MB"), http.StatusRequestEntityTooLarge, nil)
			return
		}

		file, fileHeader, err := r.FormFile(File)
		if err != nil {
			log.Printf("An error occurred during getting file: %v", err)
			server.NewResponse(w, pkg.Ptr("An error occurred during getting file"), http.StatusBadRequest, nil)
			return
		}

		defer file.Close()

		content := r.FormValue("content")
		senderID := r.FormValue("senderId")
		chatID := r.FormValue("chatId")
		sentAt := r.FormValue("sentAt")

		req, err := requests.NewSaveFileMessage(content, senderID, chatID, sentAt)
		if err != nil {
			log.Printf("Invalid body: %v", err)
			server.NewResponse(w, pkg.Ptr(fmt.Sprintf("Invalid params %s", err.Error())), http.StatusBadRequest, nil)
			return
		}

		fileMetadata := metadata.FileMetadata{
			File:        file,
			Extension:   filepath.Ext(fileHeader.Filename),
			ContentType: fileHeader.Header.Get("Content-Type"),
			SentAt:      int64(req.SentAt),
			Content:     req.Content,
			ChatID:      req.ChatID,
			SenderID:    req.SenderID,
		}

		if fileMetadata.IsValidContentType() {
			server.NewResponse(w, pkg.Ptr(fmt.Sprintf("You have %s and supported file types are: %s", fileMetadata.ContentType, strings.Join(metadata.SUPPORTED_FILE_TYPES, ", "))), http.StatusUnsupportedMediaType, nil)
			return
		}

		ctx := context.WithValue(r.Context(), FileMetadata, fileMetadata)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
