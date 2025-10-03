package dto

type DeleteNotificationsDTO struct {
	// IDs notifications to remove
	// required: true
	IDs []string `json:"ids"`
}
