package usecase

import (
	"context"
	"fmt"
	"log"
	"meeting-scheduler-service/internal/app/interfaces"
	"meeting-scheduler-service/internal/domain/dto"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

type PlannerConfig struct {
	FetchIntervalMinutes int
}

func (p *PlannerConfig) Interval() time.Duration {
	return time.Minute * time.Duration(p.FetchIntervalMinutes)
}

type planner struct {
	scheduledMeetings map[int64][]dto.PlanMeetingDTO
	unsentMeetings    []dto.PlanMeetingDTO

	mutex sync.Mutex
	cron  cron.Cron

	meetingManager interfaces.ManageMeeting

	plannerConfig PlannerConfig
}

// RerunNotStartedMeetings implements MeetingPlanner.
func (p *planner) RerunNotStartedMeetings(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Planner stopped:", ctx.Err())
			return ctx.Err()

		default:
			p.mutex.Lock()
			for i := 0; i < len(p.unsentMeetings); i++ {
				meeting := p.unsentMeetings[i]
				_, err := p.meetingManager.Start(ctx, meeting.StartMeetingDTO)
				if err != nil {
					p.unsentMeetings = append(
						append(p.unsentMeetings[:i], p.unsentMeetings[i+1:]...),
						meeting,
					)
					i--
				}
			}
			p.mutex.Unlock()
			time.Sleep(10 * time.Second)
		}
	}
}

// Listen implements MeetingPlanner.
func (p *planner) Listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Planner stopped:", ctx.Err())
			return
		default:
			now := time.Now().UTC().Unix()

			meetings := p.getMeetings(now)

			for _, meeting := range meetings {
				_, err := p.meetingManager.Start(ctx, meeting.StartMeetingDTO)
				if err != nil {
					p.addToNotStarted(meeting)
				} else {
					p.removeFromScheduled(meeting)
				}
			}

			time.Sleep(time.Second * 1)
		}
	}
}

// LoadScheduledMeetings implements MeetingPlanner.
func (p *planner) LoadScheduledMeetings(ctx context.Context) error {
	plannedMeetings, err := p.meetingManager.GetPlannedMeetings(ctx, p.plannerConfig.Interval())
	if err != nil {
		return err
	}

	for _, meeting := range plannedMeetings {
		p.addToScheduledMeetings(meeting)
	}

	return nil
}

func (p *planner) removeFromScheduled(meeting dto.PlanMeetingDTO) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	delete(p.scheduledMeetings, meeting.StartDate.Unix())
}

func (p *planner) addToNotStarted(meeting dto.PlanMeetingDTO) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.unsentMeetings = append(p.unsentMeetings, meeting)
}

func (p *planner) addToScheduledMeetings(meeting dto.PlanMeetingDTO) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if meetings, ok := p.scheduledMeetings[meeting.StartDate.Unix()]; ok {
		meetings = append(meetings, meeting)
		p.scheduledMeetings[meeting.StartDate.Unix()] = meetings
	} else {
		p.scheduledMeetings[meeting.StartDate.Unix()] = []dto.PlanMeetingDTO{meeting}
	}
}

func (p *planner) getMeetings(now int64) []dto.PlanMeetingDTO {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if meetings, ok := p.scheduledMeetings[now]; !ok {
		return []dto.PlanMeetingDTO{}
	} else {
		return meetings
	}
}

func (p *planner) StartCron(ctx context.Context) {
	spec := fmt.Sprintf("@every %dm", p.plannerConfig.FetchIntervalMinutes)

	_, err := p.cron.AddFunc(spec, func() {
		log.Println("Cron: loading scheduled meetings...")
		if err := p.LoadScheduledMeetings(ctx); err != nil {
			log.Printf("Cron: failed to load scheduled meetings: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("Failed to start cron job: %v", err)
	}

	p.cron.Start()
	log.Println("Planner cron started with interval:", spec)
}

func NewPlanner(ctx context.Context, meetingManager interfaces.ManageMeeting, config PlannerConfig) interfaces.MeetingPlanner {
	p := &planner{
		scheduledMeetings: make(map[int64][]dto.PlanMeetingDTO),
		unsentMeetings:    make([]dto.PlanMeetingDTO, 0),
		mutex:             sync.Mutex{},
		meetingManager:    meetingManager,
		cron:              *cron.New(),
		plannerConfig:     config,
	}

	p.LoadScheduledMeetings(ctx)
	p.StartCron(ctx)

	return p
}
