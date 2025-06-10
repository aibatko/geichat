package chatdb

import (
	"context"
	"time"

	"backend/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	databaseName = "chat"
	channelsCol  = "channels"
	messagesCol  = "messages"
)

type MongoStore struct {
	Client   *mongo.Client
	Channels *mongo.Collection
	Messages *mongo.Collection
}

func Open(ctx context.Context) (*MongoStore, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err = client.Connect(ctx); err != nil {
		return nil, err
	}
	db := client.Database(databaseName)
	store := &MongoStore{
		Client:   client,
		Channels: db.Collection(channelsCol),
		Messages: db.Collection(messagesCol),
	}
	// create index on messages
	_, _ = store.Messages.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    map[string]int{"channel_id": 1, "timestamp": -1},
		Options: options.Index().SetBackground(true),
	})
	return store, nil
}

func (m *MongoStore) Close(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}
