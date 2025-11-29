package usecase

import (
	"context"
	"fmt"
	"log"
	"meeting-scheduler-service/internal/app/interfaces"
	"meeting-scheduler-service/internal/domain/dto"
	"sync"
	"time"
)

type endTimestamp = int64

type meetingTerminatorImpl struct {
	ongoingMeetings map[endTimestamp][]dto.EndMeetingDTO
	mu              sync.Mutex
}

// AppendNewMeeting implements MeetingTerminator.
func (m *meetingTerminatorImpl) AppendNewMeeting(endMeetingDTO dto.EndMeetingDTO, expectedEndTimestamp int64) error {
	if expectedEndTimestamp < 1 {
		return fmt.Errorf("invalid end timestamp")
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	meetings, _ := m.ongoingMeetings[expectedEndTimestamp]
	m.ongoingMeetings[expectedEndTimestamp] = append(meetings, endMeetingDTO)

	t := time.Unix(expectedEndTimestamp, 0)
	ttt := time.Until(t)

	log.Printf("Successfully added meeting %s to termination in %d", endMeetingDTO.MeetingID, ttt)

	return nil
}

// Run implements MeetingTerminator.
func (m *meetingTerminatorImpl) Run(ctx context.Context, stopHandler func(ctx context.Context, endMeetingDTO dto.EndMeetingDTO) error) {
	log.Printf("Temrinator has been started, awaiting for meetings to terminate")
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Context done")
			return

		case t := <-ticker.C:
			now := t.UTC().Unix()

			m.mu.Lock()
			var dueKeys []int64
			for ts := range m.ongoingMeetings {
				if ts <= now {
					dueKeys = append(dueKeys, ts)
				}
			}

			var meetingsToTerminate []dto.EndMeetingDTO
			for _, ts := range dueKeys {
				meetingsToTerminate = append(meetingsToTerminate, m.ongoingMeetings[ts]...)
				delete(m.ongoingMeetings, ts)
			}
			m.mu.Unlock()

			for _, meeting := range meetingsToTerminate {
				go func(meeting dto.EndMeetingDTO) {
					if err := stopHandler(context.Background(), meeting); err != nil {
						log.Printf("Terminating meeting: %v", err)
					}
				}(meeting)
			}
		}
	}
}

func NewMeetingTerminator() interfaces.MeetingTerminator {
	return &meetingTerminatorImpl{
		ongoingMeetings: make(map[int64][]dto.EndMeetingDTO),
	}
}
