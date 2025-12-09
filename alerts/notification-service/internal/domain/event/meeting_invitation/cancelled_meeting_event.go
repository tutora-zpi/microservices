package meetinginvitation

import (
	"fmt"
	"notification-serivce/internal/domain/dto"
	"notification-serivce/internal/domain/models"
	"reflect"
	"time"
)

type CancelledMeetingEvent struct {
	Title       string        `json:"title"`
	Receivers   []dto.UserDTO `json:"members"`
	StartedDate time.Time     `json:"startedDate"`
}

func (c *CancelledMeetingEvent) Name() string {
	return reflect.TypeOf(*c).Name()
}

func (c *CancelledMeetingEvent) Notifications() []models.Notification {
	var ns []models.Notification = make([]models.Notification, len(c.Receivers))

	for i, user := range c.Receivers {
		base := models.BaseNotification()
		base.Title = "Meeting has been cancelled"
		base.Body = fmt.Sprintf("Deleted meeting %s scheduled on %s", c.Title, c.StartedDate.Format(time.RFC1123))

		base.Receiver = models.NewUser(user.ID, user.FirstName, user.LastName)

		ns[i] = *base
	}

	return ns
}
