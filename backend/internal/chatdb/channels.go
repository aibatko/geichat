package chatdb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m *MongoStore) CreateChannel(ctx context.Context, name string) (*Channel, error) {
	ch := &Channel{Name: name}
	res, err := m.Channels.InsertOne(ctx, ch)
	if err != nil {
		return nil, err
	}
	ch.ID = res.InsertedID.(primitive.ObjectID)
	return ch, nil
}

func (m *MongoStore) GetChannelByID(ctx context.Context, id primitive.ObjectID) (*Channel, error) {
	var ch Channel
	err := m.Channels.FindOne(ctx, bson.M{"_id": id}).Decode(&ch)
	if err != nil {
		return nil, err
	}
	return &ch, nil
}
