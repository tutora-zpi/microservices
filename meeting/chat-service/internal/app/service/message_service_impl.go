package service

import (
	"chat-service/internal/app/interfaces"
	"chat-service/internal/domain/broker"
	"chat-service/internal/domain/dto"
	"chat-service/internal/domain/dto/requests"
	"chat-service/internal/domain/event"
	"chat-service/internal/domain/repository"
	"context"
	"fmt"
	"log"
	"sync"
)

type messageServiceImpl struct {
	repo      repository.MessageRepository
	broker    interfaces.Broker
	fileQueue string
}

// SaveFileMessage implements interfaces.MessageService.
func (m *messageServiceImpl) SaveFileMessage(ctx context.Context, evt event.SendMessageEvent) (*dto.MessageDTO, error) {
	var wg sync.WaitGroup
	fileEvent := event.NewSendFileMessageEvent(&evt)
	var result chan *dto.MessageDTO = make(chan *dto.MessageDTO, 1)
	var errCh chan error = make(chan error, 2)

	wg.Go(func() {
		dto, err := m.repo.Save(ctx, evt)
		if err != nil {
			errCh <- err
			return
		}

		result <- dto
	})

	wg.Go(func() {
		dest := broker.NewQueueDestination(fileEvent, m.fileQueue)
		err := m.broker.Publish(ctx, fileEvent, dest)
		if err != nil {
			log.Printf("Error during publishing message: %v", err)
			errCh <- fmt.Errorf("failed to publish message")
		}
	})

	wg.Wait()
	close(result)
	close(errCh)

	for err := range errCh {
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}

	select {
	case dto := <-result:
		return dto, nil
	default:
		return nil, fmt.Errorf("message not saved due to earlier error")
	}
}

// GetMoreMessages implements interfaces.MessageService.
func (m *messageServiceImpl) GetMoreMessages(ctx context.Context, req requests.GetMoreMessages) ([]*dto.MessageDTO, error) {
	log.Println("Getting more messages...")
	return m.repo.FindMore(ctx, req)
}

func NewMessageService(repo repository.MessageRepository, broker interfaces.Broker, fileQueue string) interfaces.MessageService {
	return &messageServiceImpl{repo: repo, broker: broker, fileQueue: fileQueue}
}
