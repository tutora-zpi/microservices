package meetinginvitation

import (
	"fmt"
	"notification-serivce/internal/domain/dto"
	"notification-serivce/internal/domain/models"
	"notification-serivce/pkg"
	"reflect"
	"time"
)

type PlannedMeetingEvent struct {
	Members    []dto.UserDTO `json:"members"`
	ClassID    string        `json:"classId"`
	Title      string        `json:"title"`
	FinishDate time.Time     `json:"finishDate"`
	StartDate  time.Time     `json:"startDate"`
}

func (p *PlannedMeetingEvent) Notifications() []models.Notification {
	notifications := make([]models.Notification, len(p.Members))

	for i, receiver := range p.Members {
		notifications[i] = *p.PlannedMeetingNotification(receiver)
	}

	return notifications
}

func (p *PlannedMeetingEvent) PlannedMeetingNotification(user dto.UserDTO) *models.Notification {
	base := models.InitInvitationNotification()
	base.Receiver = models.NewUser(user.ID, user.FirstName, user.LastName)

	diff := time.Until(p.StartDate)

	base.Title = p.generateTitle(diff, user.FirstName)

	base.Body = p.generateBody(diff)

	base.RedirectionLink = p.buildLink()

	return base
}

func (p *PlannedMeetingEvent) generateTitle(diff time.Duration, firstName string) string {
	var title string
	minutesUntil := int(diff.Minutes())
	hoursUntil := int(diff.Hours())

	switch {
	case minutesUntil <= 60:
		title = fmt.Sprintf("%s, your meeting starts soon!", firstName)
	case hoursUntil < 24:
		title = fmt.Sprintf("%s, your meeting is coming up today!", firstName)
	default:
		title = fmt.Sprintf("%s, your meeting has been scheduled!", firstName)
	}

	return title
}

func (p *PlannedMeetingEvent) generateBody(diff time.Duration) string {
	hoursUntil := int(diff.Hours())
	minutesUntil := int(diff.Minutes())
	daysUntil := int(diff.Hours() / 24)
	duration := p.FinishDate.Sub(p.StartDate).Round(time.Minute)

	var timeUntil string
	switch {
	case minutesUntil < 60:
		timeUntil = fmt.Sprintf("in %d minute%s", minutesUntil, pkg.WithS(minutesUntil))
	case hoursUntil < 24:
		timeUntil = fmt.Sprintf("in %d hour%s", hoursUntil, pkg.WithS(hoursUntil))
	default:
		timeUntil = fmt.Sprintf("in %d day%s", daysUntil, pkg.WithS(daysUntil))
	}

	return fmt.Sprintf(
		"It will take place %s, starting at %02d:%02d, and will last about %s. Click down below to join.",
		timeUntil,
		p.StartDate.Hour(),
		p.StartDate.Minute(),
		pkg.FormatDuration(duration),
	)
}

func (p *PlannedMeetingEvent) Name() string {
	return reflect.TypeOf(p).Elem().Name()
}

func (p *PlannedMeetingEvent) buildLink() string {
	return fmt.Sprintf("/room/%s", p.ClassID)
}
