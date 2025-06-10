package chatdb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m *MongoStore) CreateGroup(ctx context.Context, name string, members []string) (*Group, error) {
	g := &Group{Name: name, Members: members}
	res, err := m.Groups.InsertOne(ctx, g)
	if err != nil {
		return nil, err
	}
	g.ID = res.InsertedID.(primitive.ObjectID)
	return g, nil
}

func (m *MongoStore) AddMember(ctx context.Context, id primitive.ObjectID, username string) error {
	_, err := m.Groups.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$addToSet": bson.M{"members": username}})
	return err
}

func (m *MongoStore) GetGroupsForUser(ctx context.Context, username string) ([]*Group, error) {
	cursor, err := m.Groups.Find(ctx, bson.M{"members": username})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var groups []*Group
	for cursor.Next(ctx) {
		var g Group
		if err := cursor.Decode(&g); err != nil {
			return nil, err
		}
		groups = append(groups, &g)
	}
	return groups, nil
}
