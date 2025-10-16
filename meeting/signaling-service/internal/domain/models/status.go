package models

import (
	"encoding/json"
	"fmt"
	"signaling-service/internal/domain/enum"
)

type Status struct {
	UserID     string          `json:"userId"`
	UserStatus enum.UserStatus `json:"userStatus"`
}

func DecodeStatus(body []byte) (*Status, error) {
	var status Status
	if err := json.Unmarshal(body, &status); err != nil {
		return nil, err
	}

	return &status, nil
}

func EncodeStatus(userID string, status enum.UserStatus) ([]byte, error) {

	result, err := json.Marshal(&Status{
		UserID:     userID,
		UserStatus: status,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to encode data")
	}

	return result, nil
}
