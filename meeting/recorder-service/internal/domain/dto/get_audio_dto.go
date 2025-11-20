package dto

import "recorder-service/internal/domain/dto/request"

type GetAudioDTO struct {
	UrlToAudio string `json:"urlToAudio"`
	request.GetAudioRequest
}
