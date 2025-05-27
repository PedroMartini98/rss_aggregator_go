package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/PedroMartini98/rss_aggregator_go/internal/database"
	"github.com/PedroMartini98/rss_aggregator_go/internal/response"
	"github.com/go-chi/chi"
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

func (h *userHandler) Follow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedUnvalidated := chi.URLParam(r, "feedID")

	feedValidated, err := uuid.Parse(feedUnvalidated)
	if err != nil {
		response.WithError(w, http.StatusBadRequest, fmt.Sprintf("please submit a valid feed id: %v", err))
		return
	}

	follow, err := h.dbQueries.CreateFollow(r.Context(), database.CreateFollowParams{
		UserID:    user.ID,
		CreatedAt: time.Now(),
		FeedID:    feedValidated,
	})

	if err != nil {
		response.WithError(w, http.StatusBadRequest, fmt.Sprintf("unable to create follow: %v", err))
		return
	}
	response.WithJson(w, http.StatusCreated, follow)
}

func (h *userHandler) Unfollow(w http.ResponseWriter, r *http.Request, user database.User) {

	feedUnvalidated := chi.URLParam(r, "feedID")

	feedValidated, err := uuid.Parse(feedUnvalidated)
	if err != nil {
		response.WithError(w, http.StatusBadRequest, fmt.Sprintf("please submit a valid feedID: %v", err))
		return
	}

	_, err = h.dbQueries.DeleteFollow(r.Context(), database.DeleteFollowParams{
		UserID: user.ID,
		FeedID: feedValidated,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			response.WithError(w, http.StatusBadRequest, fmt.Sprintf("there is no feed follow from this user in this id"))
			return
		}
		response.WithError(w, http.StatusInternalServerError, fmt.Sprintf("failed to delete follow in the database:%v", err))
		return
	}
	response.WithJson(w, http.StatusOK, "Sucessfully deleted follow")
}

func (h *userHandler) GetFollows(w http.ResponseWriter, r *http.Request, user database.User) {

	feeds, err := h.dbQueries.GetUserFollows(r.Context(), user.ID)
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, fmt.Sprintf("failed to get user follows: %v", err))
		return
	}

	response.WithJson(w, http.StatusAccepted, feeds)

}

func (h *userHandler) GetPosts(w http.ResponseWriter, r *http.Request, user database.User) {

	limitString := chi.URLParam(r, "limit")

	limitNumber := 10
	var err error

	if limitString != "" {
		limitNumber, err = strconv.Atoi(limitString)
		if err != nil || limitNumber <= 0 {
			response.WithError(w, http.StatusBadRequest, fmt.Sprintf("please enter a valid limit number:%v", err))
			return
		}

		posts, err := h.dbQueries.GetPostsForUser(r.Context(), database.GetPostsForUserParams{UserID: user.ID, Limit: int32(limitNumber)})
		if err != nil {
			response.WithError(w, http.StatusInternalServerError, fmt.Sprintf("failed to get posts from database: %v", err))
			return
		}
		response.WithJson(w, http.StatusOK, posts)
	}
}
