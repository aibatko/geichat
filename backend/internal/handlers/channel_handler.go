package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"backend/internal/chatdb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChannelHandler struct {
	Store *chatdb.MongoStore
}

type createChannelRequest struct {
	Name string `json:"name"`
}

type channelResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *ChannelHandler) CreateChannel(w http.ResponseWriter, r *http.Request) {
	var req createChannelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	ch, err := h.Store.CreateChannel(r.Context(), req.Name)
	if err != nil {
		http.Error(w, "could not create channel", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(channelResponse{ID: ch.ID.Hex(), Name: ch.Name})
}

func (h *ChannelHandler) ListMessages(w http.ResponseWriter, r *http.Request) {
	channelID := r.URL.Query().Get("channel")
	before := r.URL.Query().Get("before")
	if channelID == "" {
		http.Error(w, "channel required", http.StatusBadRequest)
		return
	}
	chID, err := primitive.ObjectIDFromHex(channelID)
	if err != nil {
		http.Error(w, "invalid channel", http.StatusBadRequest)
		return
	}
	var beforeTime time.Time
	if before != "" {
		ts, err := strconv.ParseInt(before, 10, 64)
		if err == nil {
			beforeTime = time.Unix(ts, 0)
		}
	}
	var msgs []*chatdb.Message
	if beforeTime.IsZero() {
		msgs, err = h.Store.GetRecentMessages(context.Background(), chID, 50)
	} else {
		msgs, err = h.Store.GetMessagesBefore(context.Background(), chID, beforeTime, 50)
	}
	if err != nil {
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(msgs)
}
