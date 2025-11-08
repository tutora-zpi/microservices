// @title Chat Serivce API
// @version 2.0
// @description Service responsible for chat feature
// @host localhost:8002
package handlers

import (
	"chat-service/internal/app/interfaces"
	"chat-service/internal/domain/dto/requests"
	"chat-service/internal/domain/metadata"
	"chat-service/internal/infrastructure/server"
	"chat-service/pkg"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Handlable interface {
	IsAuth(next http.Handler) http.Handler
	FindChat(w http.ResponseWriter, r *http.Request)
	DeleteChat(w http.ResponseWriter, r *http.Request)
	FetchMoreMessages(w http.ResponseWriter, r *http.Request)
	CreateGeneralChat(w http.ResponseWriter, r *http.Request)
	UpdateMembersInChat(w http.ResponseWriter, r *http.Request)

	UploadFile(w http.ResponseWriter, r *http.Request)
}

type handlers struct {
	chatService    interfaces.ChatService
	messageService interfaces.MessageService
	fileService    interfaces.FileService
}

// UpdateMembersInChat godoc
// @Summary      Update chat members
// @Description  Updates the list of members in a specific chat. Requires chat ID and at least one member ID.
// @Tags         Chats
// @Accept       json
// @Produce      json
// @Param        request body requests.UpdateChatMembers true "Chat members update payload"
// @Success      200 {object} server.Response
// @Failure      400 {object} server.Response
// @Failure      401 {object} server.Response
// @Router       /api/v1/chats/update-members [put]
func (h *handlers) UpdateMembersInChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dto, ok := ctx.Value(dtoKey).(*requests.UpdateChatMembers)
	if !ok {
		server.NewResponse(w, pkg.Ptr("Invalid bodies structure"), http.StatusBadRequest, nil)
		return
	}

	if err := h.chatService.UpdateChatMember(ctx, *dto); err != nil {
		server.NewResponse(w, pkg.Ptr(fmt.Sprintf("Something went wrong: %s", err.Error())), http.StatusBadRequest, nil)
		return
	}

	server.NewResponse(w, nil, http.StatusOK, nil)
}

// UploadFile godoc
// @Summary      Upload a file and create a file message
// @Description  Uploads a file to the server, saves its metadata, and creates a chat message referencing the uploaded file.
// @Tags         Messages
// @Accept       multipart/form-data
// @Produce      json
// @Param        id path string true "Chat ID"
// @Param        file formData file true "File to upload"
// @Param        content formData string false "Optional message content"
// @Param        senderId formData string true "UUID of the sender"
// @Param        sentAt formData int64 true "Unix timestamp when the message was sent"
// @Success      201 {object} dto.MessageDTO "File uploaded successfully"
// @Failure      400 {object} server.Response "Invalid parameters or file metadata"
// @Failure      401 {object} server.Response "Unauthorized access"
// @Failure      415 {object} server.Response "Unsupported media type"
// @Failure      422 {object} server.Response "File metadata has been lost"
// @Failure      500 {object} server.Response "Internal server error while saving the file"
// @Router       /api/v1/chats/{id}/upload-file [post]
func (h *handlers) UploadFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fileMetadata, ok := ctx.Value(FileMetadata).(metadata.FileMetadata)
	if !ok {
		server.NewResponse(w, pkg.Ptr("File metadata has been lost"), http.StatusUnprocessableEntity, nil)
		return
	}

	urlToFile, err := h.fileService.Save(ctx, &fileMetadata)
	log.Print(urlToFile)
	if err != nil {
		server.NewResponse(w, pkg.Ptr("Failed to save data"), http.StatusInternalServerError, nil)
		return
	}

	message := *fileMetadata.NewFileMessage(urlToFile)
	log.Print(message.FileLink)

	result, err := h.messageService.SaveFileMessage(ctx, message)
	if err != nil {
		server.NewResponse(w, pkg.Ptr(err.Error()), http.StatusBadRequest, nil)
		return
	}

	server.NewResponse(w, nil, http.StatusCreated, *result)
}

// CreateGeneralChat godoc
// @Summary Create a new general chat
// @Description Creates a new general chat with specified members
// @Tags Chats
// @Accept json
// @Produce json
// @Param body body requests.CreateGeneralChat true "CreateGeneralChat DTO"
// @Success 200 {object} dto.ChatDTO
// @Failure 400 {object} server.Response "Invalid request"
// @Router /api/v1/chats/general [post]
func (h *handlers) CreateGeneralChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dto, ok := ctx.Value(dtoKey).(*requests.CreateGeneralChat)
	if !ok {
		server.NewResponse(w, pkg.Ptr("Invalid bodies structure"), http.StatusBadRequest, nil)
		return
	}

	newChat, err := h.chatService.CreateChat(ctx, *dto)
	if err != nil {
		server.NewResponse(w, pkg.Ptr(fmt.Sprintf("Failed to create chat: %v", err)), http.StatusBadRequest, nil)
		return
	}

	server.NewResponse(w, nil, http.StatusOK, *newChat)
}

// DeleteChat godoc
// @Summary Delete a chat
// @Description Deletes chat by ID
// @Tags Chats
// @Param id path string true "Chat ID"
// @Success 204 "No Content"
// @Failure 400 {object} server.Response "No id in url"
// @Failure 404 {object} server.Response "Failed to remove chat"
// @Router /api/v1/chats/{id} [delete]
func (h *handlers) DeleteChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, ok := mux.Vars(r)["id"]

	if !ok {
		server.NewResponse(w, pkg.Ptr("No id in url"), http.StatusBadRequest, nil)
		return
	}

	if err := h.chatService.DeleteChat(ctx, requests.DeleteChat{ChatID: id}); err != nil {
		server.NewResponse(w, pkg.Ptr(fmt.Sprintf("Failed to remove chat: %v", err)), http.StatusNotFound, nil)
		return
	}

	server.NewResponse(w, nil, http.StatusNoContent, nil)
}

// FetchMoreMessages godoc
// @Summary Get more messages
// @Description Fetches more messages from a chat with optional pagination
// @Tags Messages
// @Param id path string true "Chat ID"
// @Param limit query int false "Number of messages to fetch"
// @Param lastMessageId query string false "Last message ID for pagination"
// @Success 200 {array} dto.MessageDTO
// @Failure 400 {object} server.Response "Invalid request"
// @Failure 404 {object} server.Response "No more messages"
// @Router /api/v1/chats/{id}/messages [get]
func (h *handlers) FetchMoreMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, ok := mux.Vars(r)["id"]
	if !ok {
		server.NewResponse(w, pkg.Ptr("No id in url"), http.StatusBadRequest, nil)
		return
	}

	lastMessageID := r.URL.Query().Get("lastMessageId")
	limit := r.URL.Query().Get("limit")

	dto := requests.NewGetMoreMessages(id, limit, lastMessageID)

	messages, err := h.messageService.GetMoreMessages(ctx, *dto)
	if err != nil {
		server.NewResponse(w, pkg.Ptr(fmt.Sprintf("Failed to fetch more messages: %v", err)), http.StatusNotFound, nil)
		return
	}

	server.NewResponse(w, nil, http.StatusOK, messages)
}

// FindChat godoc
// @Summary Get chat by ID
// @Description Retrieves chat details by ID
// @Tags Chats
// @Param id path string true "Chat ID"
// @Success 200 {object} dto.ChatDTO
// @Failure 400 {object} server.Response "Invalid ID"
// @Failure 404 {object} server.Response "Chat not found"
// @Router /api/v1/chats/{id} [get]
func (h *handlers) FindChat(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, ok := mux.Vars(r)["id"]
	if !ok {
		server.NewResponse(w, pkg.Ptr("No id in url"), http.StatusBadRequest, nil)
		return
	}

	dto := requests.GetChat{
		ID: id,
	}

	chatDto, err := h.chatService.FindChat(ctx, dto)
	if err != nil {
		server.NewResponse(w, pkg.Ptr(fmt.Sprintf("Chat not found: %v", err)), http.StatusNotFound, nil)
		return
	}

	server.NewResponse(w, nil, http.StatusOK, *chatDto)
}

func NewHandlers(chatService interfaces.ChatService,
	messageService interfaces.MessageService, fileService interfaces.FileService) Handlable {
	return &handlers{chatService: chatService, messageService: messageService, fileService: fileService}
}
