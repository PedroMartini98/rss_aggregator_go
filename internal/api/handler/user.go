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

type userHandler struct {
	dbQueries *database.Queries
}

func NewUserHandler(dbQueries *database.Queries) *userHandler {
	return &userHandler{
		dbQueries: dbQueries,
	}
}

func (h *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

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

	response.WithJson(w, http.StatusCreated, user)

}

func (h *userHandler) GetUser(w http.ResponseWriter, r *http.Request, user database.User) {

	response.WithJson(w, http.StatusOK, user)

}
