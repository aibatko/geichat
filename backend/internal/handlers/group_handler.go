package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"backend/internal/auth"
	"backend/internal/chatdb"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GroupHandler manages group operations
type GroupHandler struct {
	Store *chatdb.MongoStore
}

// CreateGroupRequest represents input for creating a group
type CreateGroupRequest struct {
	Name string `json:"name"`
}

// CreateGroup creates a new group for the authenticated user
func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	token := getToken(r)
	claims, err := auth.ValidateToken(token)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var req CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	g, err := h.Store.CreateGroup(r.Context(), req.Name, []string{claims.Username})
	if err != nil {
		http.Error(w, "could not create group", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(g)
}

// InviteRequest represents input for inviting a user to a group
type InviteRequest struct {
	Username string `json:"username"`
}

// Invite adds a user to the specified group
func (h *GroupHandler) Invite(w http.ResponseWriter, r *http.Request) {
	token := getToken(r)
	if _, err := auth.ValidateToken(token); err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	idHex := vars["id"]
	gid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		http.Error(w, "invalid group", http.StatusBadRequest)
		return
	}
	var req InviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Username == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := h.Store.AddMember(r.Context(), gid, req.Username); err != nil {
		http.Error(w, "failed to add member", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListGroups returns all groups the authenticated user is part of
func (h *GroupHandler) ListGroups(w http.ResponseWriter, r *http.Request) {
	token := getToken(r)
	claims, err := auth.ValidateToken(token)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	groups, err := h.Store.GetGroupsForUser(r.Context(), claims.Username)
	if err != nil {
		http.Error(w, "failed to fetch groups", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(groups)
}

func getToken(r *http.Request) string {
	authz := r.Header.Get("Authorization")
	if strings.HasPrefix(authz, "Bearer ") {
		return strings.TrimPrefix(authz, "Bearer ")
	}
	return r.URL.Query().Get("token")
}
