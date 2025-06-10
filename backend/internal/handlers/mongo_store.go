package handlers

import (
	"backend/internal/chatdb"
	"backend/internal/hub"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoMessageStore struct {
	Store *chatdb.MongoStore
}

func (m *MongoMessageStore) Save(ctx context.Context, msg hub.StoredMessage) error {
	channelID, err := primitive.ObjectIDFromHex(msg.ChannelID)
	if err != nil {
		return err
	}
	return m.Store.SaveMessage(ctx, &chatdb.Message{
		ChannelID: channelID,
		Username:  msg.Username,
		Content:   msg.Content,
		Timestamp: msg.Timestamp,
	})
}
