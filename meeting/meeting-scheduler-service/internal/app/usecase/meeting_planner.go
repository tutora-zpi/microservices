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

	mutex sync.RWMutex
	cron  cron.Cron

	meetingManager interfaces.ManageMeeting

	plannerConfig PlannerConfig
}

// RerunNotStartedMeetings implements MeetingPlanner.
func (p *planner) RerunNotStartedMeetings(ctx context.Context) error {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Planner stopped:", ctx.Err())
			return ctx.Err()
		case <-ticker.C:
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
		}
	}
}

// Listen implements MeetingPlanner.
func (p *planner) Listen(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Planner stopped:", ctx.Err())
			return
		case <-ticker.C:
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
		}
	}
}

// LoadScheduledMeetings implements MeetingPlanner.
func (p *planner) LoadScheduledMeetings(ctx context.Context) error {
	plannedMeetings, err := p.meetingManager.LoadMorePlannedMeetings(ctx, p.plannerConfig.Interval())
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

	delete(p.scheduledMeetings, meeting.StartDate.UTC().Unix())
}

func (p *planner) addToNotStarted(meeting dto.PlanMeetingDTO) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.unsentMeetings = append(p.unsentMeetings, meeting)
}

func (p *planner) addToScheduledMeetings(meeting dto.PlanMeetingDTO) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	startTimestamp := meeting.StartDate.UTC().Unix()

	if meetings, ok := p.scheduledMeetings[startTimestamp]; ok {
		meetings = append(meetings, meeting)
		p.scheduledMeetings[startTimestamp] = meetings
	} else {
		p.scheduledMeetings[startTimestamp] = []dto.PlanMeetingDTO{meeting}
	}
}

func (p *planner) getMeetings(now int64) []dto.PlanMeetingDTO {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if meetings, ok := p.scheduledMeetings[now]; !ok {
		return []dto.PlanMeetingDTO{}
	} else {
		return meetings
	}
}

func (p *planner) StartCron(ctx context.Context) error {
	spec := fmt.Sprintf("@every %dm", p.plannerConfig.FetchIntervalMinutes)

	_, err := p.cron.AddFunc(spec, func() {
		log.Println("Cron: loading scheduled meetings...")
		if err := p.LoadScheduledMeetings(ctx); err != nil {
			log.Printf("Cron: failed to load scheduled meetings: %v", err)
		}
	})
	if err != nil {
		return fmt.Errorf("failed to add cron job: %w", err)
	}

	p.cron.Start()
	log.Println("Planner cron started with interval:", spec)

	go func() {
		<-ctx.Done()
		log.Println("Stopping cron...")
		p.cron.Stop()
	}()

	return nil
}

func NewPlanner(rootCtx context.Context, meetingManager interfaces.ManageMeeting, config PlannerConfig) interfaces.MeetingPlanner {
	p := &planner{
		scheduledMeetings: make(map[int64][]dto.PlanMeetingDTO),
		unsentMeetings:    make([]dto.PlanMeetingDTO, 0),
		mutex:             sync.RWMutex{},
		meetingManager:    meetingManager,
		cron:              *cron.New(),
		plannerConfig:     config,
	}

	p.LoadScheduledMeetings(rootCtx)
	p.StartCron(rootCtx)

	return p
}
