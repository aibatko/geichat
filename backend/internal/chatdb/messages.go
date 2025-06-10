package chatdb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *MongoStore) SaveMessage(ctx context.Context, msg *Message) error {
	msg.Timestamp = time.Now()
	_, err := m.Messages.InsertOne(ctx, msg)
	return err
}

func (m *MongoStore) GetRecentMessages(ctx context.Context, channelID primitive.ObjectID, limit int64) ([]*Message, error) {
	opts := options.Find().SetSort(bson.D{{"timestamp", -1}}).SetLimit(limit)
	cursor, err := m.Messages.Find(ctx, bson.M{"channel_id": channelID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var messages []*Message
	for cursor.Next(ctx) {
		var msg Message
		if err := cursor.Decode(&msg); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}
	// reverse to chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	return messages, nil
}

func (m *MongoStore) GetMessagesBefore(ctx context.Context, channelID primitive.ObjectID, before time.Time, limit int64) ([]*Message, error) {
	filter := bson.M{"channel_id": channelID, "timestamp": bson.M{"$lt": before}}
	opts := options.Find().SetSort(bson.D{{"timestamp", -1}}).SetLimit(limit)
	cursor, err := m.Messages.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var messages []*Message
	for cursor.Next(ctx) {
		var msg Message
		if err := cursor.Decode(&msg); err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	return messages, nil
}
