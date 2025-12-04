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
	errCh := make(chan error, 3)
	var wg sync.WaitGroup

	wg.Go(func() {
		_, err := m.collectionChat.UpdateOne(
			ctx,
			bson.M{"_id": dto.ChatID},
			bson.M{"$pull": bson.M{"messagesIds": dto.MessageID}},
		)
		if err != nil {
			log.Printf("Not found chat with id: %s", dto.ChatID)
			errCh <- fmt.Errorf("failed to update chat")
		}
	})

	wg.Go(func() {
		filter := bson.M{"_id": dto.MessageID}
		result, err := m.collectionMessage.DeleteOne(ctx, filter)
		if err != nil || result.DeletedCount < 1 {
			log.Printf("Failed to remove message %v, deleted count is %d", err, result.DeletedCount)
			errCh <- fmt.Errorf("not found message with id %s", dto.MessageID)
		}
	})

	wg.Go(func() {
		filter := bson.M{"messageId": dto.MessageID}
		_, err := m.collectionReactions.DeleteMany(ctx, filter)
		if err != nil {
			log.Println("Failed to remove reactions")
			errCh <- fmt.Errorf("failed to delete reactions")
		}
	})

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for e := range errCh {
		if e != nil {
			return e
		}
	}

	return nil
}

// FindMore implements repository.MessageRepository.
func (m *messageRepoImpl) FindMore(ctx context.Context, req requests.GetMoreMessages) ([]*dto.MessageDTO, error) {
	var messages []models.Message

	filter := bson.M{"chatId": req.ID}
	if req.LastMessageId != nil {
		msgID, err := bson.ObjectIDFromHex(*req.LastMessageId)
		if err != nil {
			return nil, fmt.Errorf("invalid object id")
		}
		filter["_id"] = bson.M{"$lt": msgID}
	}

	opts := options.Find().SetSort(bson.M{"sentAt": -1}).SetLimit(int64(req.Limit))
	cursor, err := m.collectionMessage.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("not found more messages")
	}

	if err := cursor.All(ctx, &messages); err != nil {
		return nil, fmt.Errorf("failed to decode messages")
	}

	messageCh := make(chan *dto.MessageDTO, len(messages))
	var wg sync.WaitGroup

	for _, msg := range messages {
		wg.Go(func() {
			var reply *models.Message
			if msg.ReplyToID != nil {
				var r models.Message
				if err := m.collectionMessage.FindOne(ctx, bson.M{"_id": msg.ReplyToID}).Decode(&r); err == nil {
					reply = &r
				}
			}

			var reactions []models.Reaction
			cursor, err := m.collectionReactions.Find(ctx, bson.M{"messageId": msg.ID})
			if err == nil {
				_ = cursor.All(ctx, &reactions)
			}

			messageCh <- dto.NewMessageDTO(msg, reply, reactions)
		})
	}

	wg.Wait()
	close(messageCh)

	var messageDTOs []*dto.MessageDTO
	for m := range messageCh {
		messageDTOs = append(messageDTOs, m)
	}

	if len(messageDTOs) == 0 {
		return nil, fmt.Errorf("No messages found")
	}

	sort.Slice(messageDTOs, func(i, j int) bool {
		return messageDTOs[i].SentAt < messageDTOs[j].SentAt
	})

	return messageDTOs, nil
}

// React implements repository.MessageRepository.
func (m *messageRepoImpl) React(ctx context.Context, event event.ReactOnMessageEvent) error {
	reaction, err := models.NewReaction(event.UserID, event.MessageID, event.Emoji, event.SentAt)
	if err != nil {
		return err
	}

	var existing models.Reaction
	err = m.collectionReactions.FindOne(
		ctx,
		bson.M{"messageId": reaction.MessageID, "userId": reaction.UserID},
	).Decode(&existing)

	if err == nil {
		_, err := m.collectionReactions.UpdateByID(
			ctx,
			existing.ID,
			bson.M{"$set": bson.M{
				"emoji":  reaction.Emoji,
				"sentAt": reaction.SentAt,
			}},
		)
		return err
	}

	if err != mongo.ErrNoDocuments {
		return err
	}

	res, err := m.collectionReactions.InsertOne(ctx, reaction)
	if err != nil {
		return fmt.Errorf("failed to insert reaction: %v", err)
	}

	_, _ = m.collectionMessage.UpdateByID(
		ctx,
		reaction.MessageID,
		bson.M{"$push": bson.M{"reactionsIds": res.InsertedID}},
	)

	return nil
}

// Reply implements repository.MessageRepository.
func (m *messageRepoImpl) Reply(ctx context.Context, event event.ReplyOnMessageEvent) error {
	newReply := models.NewMessage(event.ChatID, event.SenderID, event.Content, event.MessageID, event.SentAt, event.FileLink, event.FileName)
	if newReply == nil {
		return fmt.Errorf("invalid reply message format")
	}

	replyObjectID, err := bson.ObjectIDFromHex(event.ReplyToMessageID)
	if err != nil {
		log.Print("Failed to cast string from event to oid")
		return fmt.Errorf("cast failed, invalid hex from ObjectID")
	}

	newReply.ReplyToID = &replyObjectID

	if _, err := m.collectionMessage.InsertOne(ctx, newReply); err != nil {
		return fmt.Errorf("failed to insert reply message: %w", err)
	}

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
func (m *messageRepoImpl) Save(ctx context.Context, event event.SendMessageEvent) (*dto.MessageDTO, error) {
	log.Printf("Saving message: %s", event.Content)
	newMessage := models.NewMessage(event.ChatID, event.SenderID, event.Content, event.MessageID, event.SentAt, event.FileLink, event.FileName)

	_, err := m.collectionMessage.InsertOne(ctx, newMessage)
	if err != nil {
		log.Printf("Failed to insert message: %v", err)
		return nil, fmt.Errorf("failed to insert message")
	}

	_, err = m.collectionChat.UpdateOne(
		ctx,
		bson.M{"_id": event.ChatID},
		bson.M{"$push": bson.M{"messagesIds": newMessage.ID}},
	)
	if err != nil {
		log.Printf("Failed to update chat with message ID: %v", err)
		return nil, fmt.Errorf("failed to update chat with message")
	}

	return dto.NewMessageDTO(*newMessage, nil, []models.Reaction{}), nil
}

func NewMessageRepository(client *mongo.Client, mongoConfig mongoConn.MongoConfig) repository.MessageRepository {
	database := mongoConfig.Database

	chat := client.Database(database).Collection(mongoConfig.ChatCollection)
	message := client.Database(database).Collection(mongoConfig.MessagesCollection)
	reactions := client.Database(database).Collection(mongoConfig.ReactionCollection)

	return &messageRepoImpl{collectionMessage: message, collectionChat: chat, collectionReactions: reactions}
}
