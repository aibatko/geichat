package hub

import (
	"context"
	"time"
)

type StoredMessage struct {
	ID        string
	ChannelID string
	Username  string
	Content   string
	Timestamp time.Time
}

type MessageStore interface {
	Save(ctx context.Context, msg StoredMessage) error
}
