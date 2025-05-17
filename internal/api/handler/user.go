package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/PedroMartini98/rss_aggregator_go/internal/database"
	"github.com/PedroMartini98/rss_aggregator_go/internal/response"
	"github.com/google/uuid"
)

type UserHandler struct {
	dbQueries *database.Queries
}

func NewUserHandler(dbQueries *database.Queries) *UserHandler {
	return &UserHandler{
		dbQueries: dbQueries,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	type requestBody struct {
		Name string `json:"name"`
	}

	var req requestBody

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.WithError(w, 500, "Error unmarshaling request body in CreateUser")
		return
	}

	user, err := h.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      req.Name,
	})
	if err != nil {
		response.WithError(w, 500, fmt.Sprintf("Couldn't create user: %v", err))
	}

	userWithTags := response.AddJsonTagToUserStuct(user)

	response.WithJson(w, http.StatusCreated, userWithTags)

}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ApiKey := r.Header.Get("ApiKey")
	if ApiKey == "" {
		response.WithError(w, http.StatusForbidden, fmt.Sprintf("failed to get ApiKey"))
		return
	}

	user, err := h.dbQueries.GetUserByApiKey(r.Context(), ApiKey)
	if err != nil {
		response.WithError(w, http.StatusForbidden, fmt.Sprintf("failed to get user: %v", err))
		return
	}

	userWithTags := response.AddJsonTagToUserStuct(user)

	response.WithJson(w, http.StatusOK, userWithTags)

}
