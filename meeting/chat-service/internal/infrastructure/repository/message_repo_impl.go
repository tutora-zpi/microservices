package repository

import (
	"chat-service/internal/domain/dto"
	"chat-service/internal/domain/dto/requests"
	"chat-service/internal/domain/event"
	"chat-service/internal/domain/models"
	"chat-service/internal/domain/repository"
	mongoConn "chat-service/internal/infrastructure/mongo"
	"context"
	"fmt"
	"log"
	"sort"
	"sync"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type messageRepoImpl struct {
	collectionMessage   *mongo.Collection
	collectionChat      *mongo.Collection
	collectionReactions *mongo.Collection
}

// Delete implements repository.MessageRepository.
func (m *messageRepoImpl) Delete(ctx context.Context, dto requests.DeleteMessage) error {
	var wg sync.WaitGroup
	var err error
	var mutex sync.Mutex

	wg.Go(func() {
		_, err := m.collectionChat.UpdateOne(
			ctx,
			bson.M{"_id": dto.ChatID},
			bson.M{"$pull": bson.M{"messagesIds": dto.MessageID}},
		)
		if err != nil {
			log.Printf("Not found chat with id: %s", dto.ChatID)
		}
	})

	wg.Go(func() {
		filter := bson.M{"_id": dto.MessageID}
		result, err := m.collectionMessage.DeleteOne(ctx, filter)
		if err != nil || result.DeletedCount < 1 {
			log.Printf("Failed to remove message %v, deleted count is %d", err, result.DeletedCount)
			mutex.Lock()
			err = fmt.Errorf("Not found message with %s", dto.MessageID)
			mutex.Unlock()
		}
	})

	wg.Go(func() {
		filter := bson.M{"messageId": dto.MessageID}
		result, err := m.collectionReactions.DeleteMany(ctx, filter)
		if err != nil || result.DeletedCount < 1 {
			log.Println("Failed to remove reactions maybe not found")
		}
	})

	wg.Wait()

	return err
}

// FindMore implements repository.MessageRepository.
func (m *messageRepoImpl) FindMore(ctx context.Context, req requests.GetMoreMessages) ([]*dto.MessageDTO, error) {
	var messages []models.Message
	var messageDTOs []*dto.MessageDTO
	var mu sync.Mutex
	var wg sync.WaitGroup

	filter := bson.M{"chatId": req.ID}
	if req.LastMessageId != nil {
		filter["_id"] = bson.M{"$lt": req.LastMessageId}
	}

	opts := options.Find().SetSort(bson.M{"sentAt": -1}).SetLimit(int64(req.Limit))
	cursor, err := m.collectionMessage.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("not found more messages")
	}

	if err := cursor.All(ctx, &messages); err != nil {
		return nil, fmt.Errorf("failed to decode messages")
	}

	for _, message := range messages {
		wg.Go(func() {
			var reply *models.Message = nil
			if message.ReplyToID != nil {
				var r models.Message
				//find parent message
				if err := m.collectionMessage.FindOne(ctx, bson.M{"_id": message.ReplyToID}).Decode(&r); err == nil {
					reply = &r
				}
			}

			var reactions []models.Reaction
			cursor, err := m.collectionReactions.Find(ctx, bson.M{"messageId": message.ID})
			if err == nil {
				_ = cursor.All(ctx, &reactions)
			}

			messageDTO := dto.NewMessageDTO(message, reply, reactions)

			mu.Lock()
			messageDTOs = append(messageDTOs, messageDTO)
			mu.Unlock()
		})
	}
	wg.Wait()

	sort.Slice(messageDTOs, func(i, j int) bool {
		return messageDTOs[i].SentAt.Before(messageDTOs[j].SentAt)
	})

	return messageDTOs, nil
}

// React implements repository.MessageRepository.
func (m *messageRepoImpl) React(ctx context.Context, event event.ReactMessageOnEvent) error {
	reaction, err := models.NewReaction(event.UserID, event.MessageID, event.Emoji, event.SentAt)
	if err != nil {
		return err
	}

	if _, err := m.collectionReactions.InsertOne(ctx, reaction); err != nil {
		return fmt.Errorf("failed to insert reaction: %v", err)
	}

	_, _ = m.collectionMessage.UpdateByID(
		ctx,
		event.MessageID,
		bson.M{"$push": bson.M{"reactionsIds": reaction.ID}},
	)

	return nil
}

// Reply implements repository.MessageRepository.
func (m *messageRepoImpl) Reply(ctx context.Context, event event.ReplyOnMessageEvent) error {
	newReply, err := models.NewMessage(event.ChatID, event.SenderID, event.Content, event.SentAt)
	if err != nil {
		return err
	}

	newReply.ReplyToID = &event.ReplyToMessageID

	if _, err := m.collectionMessage.InsertOne(ctx, newReply); err != nil {
		return fmt.Errorf("failed to insert reply message: %w", err)
	}

	// if _, err := m.collectionMessage.UpdateByID(
	// 	ctx,
	// 	event.ReplyToMessageID,
	// 	bson.M{"$set": bson.M{"replyToId": newReply.ID}},
	// ); err != nil {
	// 	return fmt.Errorf("failed to update parent message: %w", err)
	// }

	if _, err := m.collectionChat.UpdateByID(
		ctx,
		event.ChatID,
		bson.M{"$push": bson.M{"messagesIds": newReply.ID}},
	); err != nil {
		return fmt.Errorf("failed to update chat with reply message ID: %w", err)
	}

	return nil
}

// Save implements repository.MessageRepository.
func (m *messageRepoImpl) Save(ctx context.Context, event event.SendMessageEvent) error {
	log.Printf("Saving message: %s", event.Content)
	newMessage, err := models.NewMessage(event.ChatID, event.SenderID, event.Content, event.SentAt)
	if err != nil {
		return err
	}

	_, err = m.collectionMessage.InsertOne(ctx, newMessage)
	if err != nil {
		return fmt.Errorf("failed to insert message: %v", err)
	}

	_, err = m.collectionChat.UpdateOne(
		ctx,
		bson.M{"_id": event.ChatID},
		bson.M{"$push": bson.M{"messagesIds": newMessage.ID}},
	)
	if err != nil {
		return fmt.Errorf("failed to update chat with message ID: %v", err)
	}

	return nil
}

func NewMessageRepository(client *mongo.Client, mongoConfig mongoConn.MongoConfig) repository.MessageRepository {
	database := mongoConfig.Database

	chat := client.Database(database).Collection(mongoConfig.ChatCollection)
	message := client.Database(database).Collection(mongoConfig.MessagesCollection)
	reactions := client.Database(database).Collection(mongoConfig.ReactionCollection)

	return &messageRepoImpl{collectionMessage: message, collectionChat: chat, collectionReactions: reactions}
}
