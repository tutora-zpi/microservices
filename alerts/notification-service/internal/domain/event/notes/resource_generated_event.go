package notes

import (
	"fmt"
	"log"
	"notification-serivce/internal/domain/models"
	"reflect"
)

type ResourcesGeneratedEvent struct {
	ClassID   string   `json:"class_id"`
	MeetingID string   `json:"meeting_id"`
	Status    string   `json:"status"`
	MemberIDs []string `json:"member_ids"`
}

func (r *ResourcesGeneratedEvent) Notifications() ([]models.Notification, error) {
	var ns []models.Notification = make([]models.Notification, len(r.MemberIDs))
	if r.Status == "FAILURE" {
		log.Printf("ResourcesGeneratedEvent has status: %s", r.Status)
		return nil, fmt.Errorf("something went wrong during generating resources")
	}

	for i, id := range r.MemberIDs {
		base := models.BaseNotification()
		base.Title = "New meeting resources available"
		base.Body = "New notes and the recording from the meeting are now available in your class."
		base.Receiver = models.NewUser(id, "", "")

		base.RedirectionLink = r.buildLink()

		ns[i] = *base
	}

	return ns, nil
}

func (r *ResourcesGeneratedEvent) Name() string {
	return reflect.TypeOf(r).Elem().Name()
}

func (r *ResourcesGeneratedEvent) buildLink() string {
	return fmt.Sprintf("/room/%s", r.ClassID)
}
