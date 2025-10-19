package repository

import (
	"chat-service/internal/domain/dto"
	"chat-service/internal/domain/dto/requests"
	"chat-service/internal/domain/models"
	"chat-service/internal/domain/repository"
	mongoConn "chat-service/internal/infrastructure/mongo"
	"context"
	"fmt"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type chatRepositoryImpl struct {
	collectionMessage   *mongo.Collection
	collectionChat      *mongo.Collection
	collectionReactions *mongo.Collection
}

// Save implements repository.ChatRepository.
func (c *chatRepositoryImpl) Save(ctx context.Context, memberIDs []string, chatID string) (*dto.ChatDTO, error) {
	var found models.Chat
	result := c.collectionChat.FindOne(ctx, bson.M{"_id": chatID})
	if result.Err() == nil {
		if err := result.Decode(&found); err == nil {
			dto := dto.NewChatDTO(found, []dto.MessageDTO{})
			return &dto, nil
		}
	}

	chatModel := models.NewChat(memberIDs, chatID)
	_, err := c.collectionChat.InsertOne(ctx, *chatModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat: %v", err)
	}

	chatDTO := dto.NewChatDTO(*chatModel, []dto.MessageDTO{})

	return &chatDTO, err
}

// Delete implements repository.ChatRepository.
func (c *chatRepositoryImpl) Delete(ctx context.Context, request requests.DeleteChat) error {
	var wg sync.WaitGroup
	mutex := sync.Mutex{}
	var messageIDs []string
	var deletionError error

	wg.Go(func() {
		filter := bson.M{"_id": request.ChatID}
		result, err := c.collectionChat.DeleteOne(ctx, filter)
		if err != nil || result.DeletedCount < 1 {
			log.Printf("Failed to delete chat: %v", err)
			mutex.Lock()
			deletionError = fmt.Errorf("chat not found: %v", err)
			mutex.Unlock()

			return
		}
	})

	wg.Go(func() {
		filter := bson.M{"chatId": request.ChatID}
		cursor, err := c.collectionMessage.Find(ctx, filter)
		if err != nil {
			return
		}

		var messages []models.Message
		_ = cursor.All(ctx, &messages)

		ids := make([]string, len(messages))
		for i, m := range messages {
			ids[i] = m.ID
		}
		mutex.Lock()
		messageIDs = ids
		mutex.Unlock()
	})

	wg.Wait()

	if len(messageIDs) > 0 {
		wg.Go(func() {
			c.collectionReactions.DeleteMany(ctx, bson.M{"messageId": bson.M{"$in": messageIDs}})
		})

		wg.Go(func() {
			c.collectionMessage.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": messageIDs}})
		})

		wg.Wait()
	}

	return deletionError
}

// Find implements repository.ChatRepository.
func (c *chatRepositoryImpl) Find(ctx context.Context, request requests.GetChat) (*dto.ChatDTO, error) {
	var wg sync.WaitGroup
	filter := bson.M{"_id": request.ID}

	var chat *models.Chat
	var messages []models.Message
	messageDTOs := make([]dto.MessageDTO, 0, request.Limit)
	var mu sync.Mutex

	wg.Go(func() {
		singleResult := c.collectionChat.FindOne(ctx, filter)
		if err := singleResult.Err(); err != nil {
			log.Printf("Failed to find chat: %v", err)
			return
		}
		if err := singleResult.Decode(&chat); err != nil {
			log.Printf("Failed to decode chat: %v", err)
			return
		}
	})

	wg.Go(func() {
		opts := options.Find().
			SetSort(bson.D{{Key: "sentAt", Value: 1}}).
			SetLimit(int64(request.Limit))

		cursor, err := c.collectionMessage.Find(ctx, bson.M{"chatId": request.ID}, opts)
		if err != nil {
			log.Printf("Failed to find messages: %v", err)
			return
		}
		if err := cursor.All(ctx, &messages); err != nil {
			log.Printf("Failed to decode messages: %v", err)
			return
		}
	})

	wg.Wait()

	if chat == nil {
		return nil, fmt.Errorf("chat not found")
	}

	for _, message := range messages {
		msg := message
		wg.Go(func() {
			var reply *models.Message
			if msg.ReplyToID != nil {
				var r models.Message
				if err := c.collectionMessage.FindOne(ctx, bson.M{"_id": msg.ReplyToID}).Decode(&r); err == nil {
					reply = &r
				}
			}

			var reactions []models.Reaction
			cursor, err := c.collectionReactions.Find(ctx, bson.M{"messageId": msg.ID})
			if err == nil {
				_ = cursor.All(ctx, &reactions)
			}

			messageDTO := dto.NewMessageDTO(msg, reply, reactions)

			mu.Lock()
			messageDTOs = append(messageDTOs, *messageDTO)
			mu.Unlock()
		})
	}

	wg.Wait()

	chatDTO := dto.NewChatDTO(*chat, messageDTOs)
	return &chatDTO, nil
}

func NewChatRepository(client *mongo.Client, mongoConfig mongoConn.MongoConfig) repository.ChatRepository {
	database := mongoConfig.Database

	chat := client.Database(database).Collection(mongoConfig.ChatCollection)
	message := client.Database(database).Collection(mongoConfig.MessagesCollection)
	reactions := client.Database(database).Collection(mongoConfig.ReactionCollection)
	return &chatRepositoryImpl{collectionMessage: message, collectionChat: chat, collectionReactions: reactions}
}
