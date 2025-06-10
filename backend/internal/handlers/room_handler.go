package handlers

import (
	"backend/internal/auth"
	"backend/internal/chatdb"
	"backend/internal/hub"
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		allowedOrigins := []string{"http://localhost:3000", "http://localhost:8080"}
		origin := r.Header.Get("Origin")
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				return true
			}
		}
		return false
	},
}

type RoomHandler struct {
	hub   *hub.Hub
	store *chatdb.MongoStore
}

func NewRoomHandler(maxPlayers int, store *chatdb.MongoStore) *RoomHandler {
	return &RoomHandler{
		hub:   hub.NewHub(maxPlayers),
		store: store,
	}
}

func (h *RoomHandler) Start() {
	go h.hub.Run()
}

func (h *RoomHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	channelParam := r.URL.Query().Get("channel")

	if token == "" {
		http.Error(w, "Token required", http.StatusBadRequest)
		return
	}

	claims, err := auth.ValidateToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	if channelParam == "" {
		http.Error(w, "Channel required", http.StatusBadRequest)
		return
	}

	chID, err := primitive.ObjectIDFromHex(channelParam)
	if err != nil {
		http.Error(w, "Invalid channel", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection", http.StatusInternalServerError)
		return
	}

	client := &hub.Client{
		Hub:      h.hub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		Username: claims.Username,
		RoomID:   chID.Hex(),
		Store:    &MongoMessageStore{Store: h.store},
	}

	// send recent messages
	msgs, err := h.store.GetRecentMessages(r.Context(), chID, 50)
	if err == nil {
		for _, m := range msgs {
			data, _ := json.Marshal(hub.Message{
				Type:      "message",
				Room:      chID.Hex(),
				Username:  m.Username,
				Content:   m.Content,
				Timestamp: m.Timestamp.Unix(),
			})
			client.Send <- data
		}
	}

	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
